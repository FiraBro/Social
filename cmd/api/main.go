package main

import (
	"log"

	"Todo-List-Go/internal/db"
	"Todo-List-Go/internal/env"
	"Todo-List-Go/internal/store"
)

const version string = "0.0.1"

func main() {
	// Load config from env or defaults
	cfg := config{
		addr: env.GetString("ADDR", ":4000"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	// Connect to DB
	dbConn, err := db.New(cfg.db.addr, cfg.db.maxIdleConns, cfg.db.maxIdleTime, cfg.db.maxOpenConns)
	if err != nil {
		log.Panic(err)
	}
	defer dbConn.Close() // ✅ Corrected

	log.Println("Database is connected")

	// Initialize store
	storage := store.NewSQLStorage(dbConn) // ✅ Renamed variable

	// Create application
	app := &application{
		config: cfg,
		store:  storage,
	}

	// Start server
	mux := app.mount()
	log.Fatal(app.run(mux))
}
