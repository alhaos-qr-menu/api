package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/alhaos-qr-menu/api/internal/config"
	"github.com/alhaos-qr-menu/api/internal/database"
	"github.com/alhaos-qr-menu/api/internal/http/handlers"
	"github.com/alhaos-qr-menu/api/internal/http/router"
	"github.com/alhaos-qr-menu/api/internal/logging"
)

func main() {
	// Bootstrap logger (до загрузки конфига)
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	// Load config
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Init structured logger with rotation
	logger := logging.New(logging.Config{
		Path:       cfg.LogPath,
		MaxSizeMB:  cfg.LogMaxSizeMB,
		MaxBackups: cfg.LogMaxBackups,
		MaxAgeDays: cfg.LogMaxAgeDays,
		Compress:   cfg.LogCompress,
	})
	slog.SetDefault(logger)

	// Background context
	ctx := context.Background()

	// Init database
	pool, err := database.New(ctx, cfg.DatabaseURL)

	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Init handlers
	h := handlers.New(pool)

	// Init Gin router
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	// CORS
	corsConfig := cors.Config{
		AllowOrigins:     []string{cfg.AllowedOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}
	r.Use(cors.New(corsConfig))

	// Setup all routes
	router.Setup(r, h)

	slog.Info("server starting", "port", cfg.Port, "env", os.Getenv("ENV"))

	if err := r.Run(":" + cfg.Port); err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}
