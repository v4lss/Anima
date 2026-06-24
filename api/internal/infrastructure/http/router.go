// Package http - router wires chi routes to handler functions.
package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/v4lss/animas/internal/infrastructure/http/middleware"
	"github.com/v4lss/animas/pkg/jwt"
)

// NewRouter creates and returns the configured chi router.
// Handler functions are stubs - they will be filled in progressively.
func NewRouter(jwtSvc *jwt.Service) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Public routes
	r.Post("/api/auth/register", stubHandler)
	r.Post("/api/auth/login", stubHandler)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(jwtSvc))

		r.Post("/api/monitors", stubHandler)
		r.Get("/api/monitors", stubHandler)
		r.Get("/api/monitors/{id}", stubHandler)
		r.Delete("/api/monitors/{id}", stubHandler)
		r.Get("/api/monitors/{id}/history", stubHandler)
	})

	return r
}

// stubHandler is a placeholder returned while handlers are being implemented.
func stubHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(`{"error":"not implemented"}`))
}
