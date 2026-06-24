// Package http - router wires chi routes to real handler functions.
package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/v4lss/animas/internal/domain/monitor"
	"github.com/v4lss/animas/internal/domain/user"
	"github.com/v4lss/animas/internal/infrastructure/http/handler"
	"github.com/v4lss/animas/internal/infrastructure/http/middleware"
	"github.com/v4lss/animas/pkg/jwt"
)

// NewRouter creates and returns the configured chi router with all real handlers.
func NewRouter(
	jwtSvc     *jwt.Service,
	userRepo   user.Repository,
	monitorRepo monitor.Repository,
	checkRepo  monitor.CheckRepository,
) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Handlers
	authHandler    := handler.NewAuthHandler(userRepo, jwtSvc)
	monitorHandler := handler.NewMonitorHandler(monitorRepo, checkRepo)

	// Public routes
	r.Post("/api/auth/register", authHandler.Register)
	r.Post("/api/auth/login",    authHandler.Login)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(jwtSvc))

		r.Post("/api/monitors",              monitorHandler.Create)
		r.Get("/api/monitors",               monitorHandler.List)
		r.Get("/api/monitors/{id}",          monitorHandler.Get)
		r.Delete("/api/monitors/{id}",       monitorHandler.Delete)
		r.Get("/api/monitors/{id}/history",  monitorHandler.History)
	})

	return r
}
