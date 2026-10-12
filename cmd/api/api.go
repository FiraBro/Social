package main

import (
	"log"
	"net/http"
	"time"

	"github.com/FiraBro/social/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// application holds the dependencies required by the API.
type application struct {
	config config
	store  store.Storage
}

// dbConfig contains database connection settings.
type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

// config contains application configuration.
type config struct {
	addr string
	db   dbConfig
	env  string
}

// mount registers routes and middleware.
func (app *application) mount() *chi.Mux {
	r := chi.NewRouter()

	// Middleware.
	r.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
	)

	// Routes.
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Welcome to Social API"))
	})

	r.Get("/v1/health", app.healthHandler)

	// Posts routes can be registered here when ready.
	// Example:
	// r.Post("/v1/posts", app.createPostHandler)

	return r
}

// run starts the HTTP server.
func (app *application) run(mux *chi.Mux) error {
	srv := &http.Server{
		Addr:              app.config.addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("Server running on %s", app.config.addr)

	return srv.ListenAndServe()
}