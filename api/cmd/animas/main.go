// main is the entry point for the Animas API server.
// It wires together config, infrastructure, application services, and HTTP.
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/v4lss/animas/internal/infrastructure/config"
	apphttp "github.com/v4lss/animas/internal/infrastructure/http"
	"github.com/v4lss/animas/pkg/jwt"
	"github.com/v4lss/animas/pkg/logger"
)

func main() {
	// ── Config ──────────────────────────────────────────────────────────────
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		logger.Fatal("failed to load config: %v", err)
	}

	// ── JWT ─────────────────────────────────────────────────────────────────
	expiry, _ := time.ParseDuration(cfg.JWT.Expiry)
	jwtSvc := jwt.New(cfg.JWT.Secret, expiry)

	// ── HTTP ─────────────────────────────────────────────────────────────────
	router := apphttp.NewRouter(jwtSvc)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.App.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// ── Graceful shutdown ────────────────────────────────────────────────────
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Info("Animas listening on :%s\n", cfg.App.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server error: %v", err)
		}
	}()

	<-done
	logger.Info("shutting down...\n")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
