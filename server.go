package main

import (
	"blog_app/app"
	"blog_app/config"
	"blog_app/db"
	"blog_app/utils/logger"

	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	// Set Gin mode based on environment
	switch env {
	case "prod", "production":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	// Load config and initialize logger
	config.LoadConfig()
	logger.InitLogger(env)

	// Initialize DB
	if err := db.Init(); err != nil {
		logger.Error("APP_ERROR: Database initialization failed", zap.Error(err))
		os.Exit(1)
	}

	// Health check
	if err := db.GetInstance().HealthCheck(); err != nil {
		logger.Error("APP_ERROR: Database health check failed", zap.Error(err))
		os.Exit(1)
	}

	// Start Gin app in a goroutine
	go app.StartApp()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Received shutdown signal, closing server gracefully...")

	// Close DB
	if err := db.GetInstance().Close(); err != nil {
		logger.Error("Error during DB shutdown", zap.Error(err))
	}

	logger.Info("Server closed successfully")
}
