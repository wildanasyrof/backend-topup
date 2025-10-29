package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/config"
	"github.com/wildanasyrof/backend-topup/internal/di"
	"github.com/wildanasyrof/backend-topup/internal/http/router"
	"github.com/wildanasyrof/backend-topup/internal/http/server"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error load config:%v:", err)
	}

	di := di.InitDI(cfg)
	app := fiber.New(
		fiber.Config{
			ErrorHandler: server.ErrorHandler(di.Logger),
		},
	)
	router.SetupRouter(app, di, cfg)
	// Channel to listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		di.Logger.Info(fmt.Sprintf("Starting server on port %s...", cfg.Server.Port))
		if err := app.Listen(":" + cfg.Server.Port); err != nil {
			di.Logger.Error(err, fmt.Sprintf("Server Listen error: %v", err))
		}
	}()

	// Block until a signal is received
	sig := <-quit
	di.Logger.Info(fmt.Sprintf("Received signal: %s. Shutting down server...", sig))

	// --- Graceful Shutdown ---
	// Give existing requests some time to finish (e.g., 30 seconds)
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown of the Fiber app
	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		di.Logger.Error(err, "Server forced to shutdown:")
	}

	// --- Cleanup Resources ---
	// Close database connection (assuming GetDB method exists in DI or similar)
	if sqlDB, err := di.GetDB().DB(); err == nil {
		di.Logger.Info("Closing database connection pool...")
		if err := sqlDB.Close(); err != nil {
			di.Logger.Error(err, "Error closing database connection pool:")
		}
	} else {
		di.Logger.Error(err, "Failed to get DB instance for closing:")
	}

	di.Logger.Info("Server shutdown complete.")
}
