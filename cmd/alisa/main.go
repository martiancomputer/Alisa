package main

import (
	"database/sql"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed dashboard/dist/*
var dashboardAssets embed.FS

type Application struct {
	DB      *sql.DB
	Discord *discordgo.Session
}

func main() {
	// Configure garbage collection threshold for low-memory environments
	os.Setenv("GOGC", "50")

	// 1. Initialize Relational Database
	db, err := sql.Open("sqlite3", "./alisa.db?_journal=WAL&_sync=NORMAL&_busy_timeout=5000")
	if err != nil {
		log.Fatalf("FATAL: SQLite instantiation failed: %v", err)
	}
	defer db.Close()
	
	// Restrict connection pool to prevent RAM spikes and SQLite lock contention
	db.SetMaxOpenConns(1)

	// 2. Initialize Discord Gateway
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("FATAL: DISCORD_TOKEN environment variable is undefined.")
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("FATAL: Discord Gateway initialization failed: %v", err)
	}

	app := &Application{
		DB:      db,
		Discord: dg,
	}

	// Request specific gateway intents required for audit and automation logic
	dg.Identify.Intents = discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildMembers |
		discordgo.IntentsGuildModeration

	if err := dg.Open(); err != nil {
		log.Fatalf("FATAL: Failed to open Gateway WebSocket: %v", err)
	}
	defer dg.Close()

	log.Println("INFO: Discord Gateway connected successfully.")

	// 3. Initialize Internal REST & Dashboard Server
	mux := http.NewServeMux()
	
	mux.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "operational"}`))
	})

	// Mount embedded filesystem
	distFS, err := fs.Sub(dashboardAssets, "dashboard/dist")
	if err != nil {
		log.Printf("WARN: Dashboard assets missing or failed to mount: %v", err)
	} else {
		mux.Handle("/", http.FileServer(http.FS(distFS)))
		log.Println("INFO: Dashboard UI mounted from memory.")
	}

	go func() {
		log.Println("INFO: Local HTTP interface listening on :8080")
		if err := http.ListenAndServe(":8080", mux); err != nil {
			log.Fatalf("FATAL: HTTP server crashed: %v", err)
		}
	}()

	// 4. Await termination signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	log.Println("INFO: Intercepted termination signal. Committing graceful shutdown.")
}