package main

import (
	"log"

	"github.com/FiraBro/social/internal/db"
	"github.com/FiraBro/social/internal/env"
	"github.com/FiraBro/social/internal/store"
)

const version = "0.0.1"

func main() {
	// Load configuration from environment variables or defaults.
	cfg := config{
		addr: env.GetString("ADDR", ":4000"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	// Connect to the database.
	dbConn, err := db.New(
		cfg.db.addr,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
		cfg.db.maxOpenConns,
	)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	defer dbConn.Close()

	log.Println("Database connected successfully")

	// Initialize storage.
	storage := store.NewSQLStorage(dbConn)

	// Create the application.
	app := &application{
		config: cfg,
		store:  storage,
	}

	// Start the HTTP server.
	mux := app.mount()
	log.Fatal(app.run(mux))
}