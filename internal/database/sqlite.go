package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB instantiates the SQLite database with memory and concurrency optimizations.
func InitDB(filepath string) *sql.DB {
	// Connect using explicit configuration arguments
	db, err := sql.Open("sqlite3", filepath+"?_journal=WAL&_sync=NORMAL&_busy_timeout=5000&_fk=true")
	if err != nil {
		log.Fatalf("FATAL: Failed to bind SQLite storage interface: %v", err)
	}

	// Enforce single-connection isolation to eliminate thread locking or memory overhead
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	// Apply low-footprint performance configurations
	pragmas := []string{
		"PRAGMA journal_mode = WAL;",
		"PRAGMA synchronous = NORMAL;",
		"PRAGMA foreign_keys = ON;",
		"PRAGMA temp_store = MEMORY;",
		"PRAGMA cache_size = -2000;", // Limit page cache memory consumption (~2MB allocation ceiling)
	}

	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			log.Fatalf("FATAL: SQL configuration directive failed (%s): %v", pragma, err)
		}
	}

	return db
}