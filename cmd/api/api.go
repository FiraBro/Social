package main

import (
	"log"
	"net/http"
	"time"

	"Todo-List-Go/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Main app struct
type application struct {
	config config
	store  store.Storage
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

// App config
type config struct {
	addr string
	db   dbConfig
	env  string
}

// Setup routes
func (app *application) mount() *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Recoverer, middleware.RequestID, middleware.Logger)

	// Routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Get("/v1/health", app.healthHandler)

	return r
}

// Run HTTP server
func (app *application) run(mux *chi.Mux) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	log.Printf("Server running on %s", app.config.addr)
	return srv.ListenAndServe()
}
