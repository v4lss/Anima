// main is the entry point for the Animas API server.
// Wires config, MongoDB, Redis, repositories, router, workers and HTTP.
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/v4lss/animas/internal/infrastructure/config"
	apphttp "github.com/v4lss/animas/internal/infrastructure/http"
	"github.com/v4lss/animas/internal/infrastructure/mongodb"
	redisinfra "github.com/v4lss/animas/internal/infrastructure/redis"
	"github.com/v4lss/animas/internal/workers"
	"github.com/v4lss/animas/pkg/jwt"
	"github.com/v4lss/animas/pkg/logger"
)

func main() {
	// Config
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		logger.Fatal("failed to load config: %v", err)
	}

	// MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Mongo.URI))
	if err != nil {
		logger.Fatal("failed to connect to MongoDB: %v", err)
	}
	if err := mongoClient.Ping(ctx, nil); err != nil {
		logger.Fatal("MongoDB ping failed: %v", err)
	}
	logger.Info("MongoDB connected\n")

	db := mongoClient.Database(cfg.Mongo.DB)

	// Repositories
	userRepo    := mongodb.NewUserRepository(db)
	monitorRepo := mongodb.NewMonitorRepository(db)
	checkRepo   := mongodb.NewCheckRepository(db)

	// Redis
	queue := redisinfra.NewQueue(cfg.Redis.Addr, cfg.Redis.Password)
	logger.Info("Redis queue ready\n")

	// JWT
	expiry, _ := time.ParseDuration(cfg.JWT.Expiry)
	jwtSvc := jwt.New(cfg.JWT.Secret, expiry)

	// HTTP
	router := apphttp.NewRouter(jwtSvc, userRepo, monitorRepo, checkRepo)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.App.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Workers
	workerCtx, workerCancel := context.WithCancel(context.Background())

	scheduler   := workers.NewScheduler(monitorRepo, queue)
	httpChecker := workers.NewHTTPChecker(monitorRepo, checkRepo, queue)
	tcpChecker  := workers.NewTCPChecker(monitorRepo, checkRepo, queue)

	go scheduler.Run(workerCtx)
	go httpChecker.Run(workerCtx)
	go tcpChecker.Run(workerCtx)

	// Graceful shutdown
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

	workerCancel()

	shutCtx, shutCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutCancel()
	srv.Shutdown(shutCtx)
	mongoClient.Disconnect(shutCtx)
}
