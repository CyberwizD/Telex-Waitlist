package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/CyberwizD/Telex-Waitlist/internal/config"
	"github.com/CyberwizD/Telex-Waitlist/internal/database"
	"github.com/CyberwizD/Telex-Waitlist/internal/handlers"
	"github.com/CyberwizD/Telex-Waitlist/internal/repository"
	"github.com/CyberwizD/Telex-Waitlist/internal/routes"
	"github.com/CyberwizD/Telex-Waitlist/internal/services"
	"github.com/CyberwizD/Telex-Waitlist/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		// Fallback to stdout if logger cannot initialize.
		slog.Error("config error", "err", err)
		os.Exit(1)
	}

	appLog := logger.New(cfg.LogLevel)

	migrateCtx, cancelMigrate := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelMigrate()

	appLog.Info("Starting database migration...")

	if err := database.Migrate(migrateCtx, appLog, cfg); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbPool, err := database.Connect(ctx, cfg)
	if err != nil {
		appLog.Error("db connection error", "err", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	waitlistRepo := repository.NewWaitlistRepository(dbPool)
	emailSvc := services.NewEmailService(cfg)
	waitlistSvc := services.NewWaitlistService(waitlistRepo, emailSvc)
	waitlistHandler := handlers.NewWaitlistHandler(waitlistSvc, cfg.AdminAPIKey)

	router := routes.SetupRouter(cfg, waitlistHandler)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		appLog.Info("server listening", "port", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	waitForShutdown(server, appLog)
}

func waitForShutdown(server *http.Server, log *slog.Logger) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("graceful shutdown failed", "err", err)
	}
	log.Info("server stopped")
}
