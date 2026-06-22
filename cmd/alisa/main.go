package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/martiancomputer/Alisa/internal/api"
	"github.com/martiancomputer/Alisa/internal/automation"
	"github.com/martiancomputer/Alisa/internal/bot/handlers"
	"github.com/martiancomputer/Alisa/internal/database"
	"github.com/martiancomputer/Alisa/internal/database/queries"
)

//go:embed dashboard/dist/*
var dashboardAssets embed.FS

func main() {
	// Constrain runtime memory growth for 512MB RAM optimization target
	os.Setenv("GOGC", "50")

	log.Println("SYSTEM: Initializing Alisa Monolithic Engine...")

	// 1. Initialize Sub-process SQLite Storage Unit
	dbConn := database.InitDB("./alisa.db")
	defer dbConn.Close()

	// Instantiate specific data repository layers
	caseRepository := &queries.CaseRepository{DB: dbConn}
	logRepository  := &queries.LogRepository{DB: dbConn}
	statsRepository := &queries.StatsRepository{DB: dbConn}

	// 2. Initialize Low-overhead Event Automation Engine
	autoEngine := &automation.Engine{
		Rules: []models.Rule{}, // Dynamically populated via migrations / DB fetches
	}

	// 3. Initialize Discord Gateway Controller
	botToken := os.Getenv("DISCORD_TOKEN")
	if botToken == "" {
		log.Fatal("CRITICAL: DISCORD_TOKEN environment variable is unassigned.")
	}

	dgSession, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatalf("CRITICAL: Gateway wrapper allocation failed: %v", err)
	}

	// Register deterministic message handler bridging into the rules loop
	msgHandler := &handlers.MessageHandler{Engine: autoEngine}
	dgSession.AddHandler(msgHandler.OnMessageCreate)

	// Configure network constraints to intercept exact operational flows
	dgSession.Identify.Intents = discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildMembers |
		discordgo.IntentsGuildModeration

	if err := dgSession.Open(); err != nil {
		log.Fatalf("CRITICAL: Failed to tie Discord Gateway socket link: %v", err)
	}
	defer dgSession.Close()
	log.Println("SYSTEM: Discord Gateway channel streaming live.")

	// 4. Construct REST API Routing Core
	apiContext := &api.APIHandler{
		CaseRepo:  caseRepository,
		LogRepo:   logRepository,
		StatsRepo: statsRepository,
	}

	httpMux := http.NewServeMux()
	
	// Bind endpoints protected through low-overhead sub-bit validation checking
	httpMux.HandleFunc("/api/v1/cases", api.AuthMiddleware(apiContext.HandleGetRecentCases))
	httpMux.HandleFunc("/api/v1/audit", api.AuthMiddleware(apiContext.HandleGetAuditFeed))

	// Mount embedded static assets into HTTP thread space
	distSubtree, err := fs.Sub(dashboardAssets, "dashboard/dist")
	if err != nil {
		log.Printf("SYSTEM: Dashboard asset subsystem unpopulated or detached: %v", err)
	} else {
		httpMux.Handle("/", http.FileServer(http.FS(distSubtree)))
		log.Println("SYSTEM: Dashboard single-page asset arrays embedded into runtime memory.")
	}

	// Execute listening context inside an isolated light thread
	go func() {
		log.Println("SYSTEM: REST/Dashboard HTTP service listening on :8080")
		if err := http.ListenAndServe(":8080", httpMux); err != nil {
			log.Fatalf("CRITICAL: Unified HTTP network interface crashed: %v", err)
		}
	}()

	// 5. Establish Signal Trap for Orderly Resource Releases
	terminationSignalChannel := make(chan os.Signal, 1)
	signal.Notify(terminationSignalChannel, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-terminationSignalChannel

	log.Println("SYSTEM: Termination flag intercepted. Deallocating memory descriptors and breaking gateway loops.")
}