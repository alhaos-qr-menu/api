package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/alhaos-qr-menu/api/internal/config"
	"github.com/alhaos-qr-menu/api/internal/db"
	"github.com/alhaos-qr-menu/api/internal/logging"
)

func main() {
	// Bootstrap logger: stdout only, used until config is loaded.
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	// Init config
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Init logger
	logger := logging.New(logging.Config{
		Path:       cfg.LogPath,
		MaxSizeMB:  cfg.LogMaxSizeMB,
		MaxBackups: cfg.LogMaxBackups,
		MaxAgeDays: cfg.LogMaxAgeDays,
		Compress:   cfg.LogCompress,
	})
	slog.SetDefault(logger)

	// Get context
	ctx := context.Background()

	// Init database pool
	pool, err := db.New(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Init router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// CORS config init
	corsConfig := cors.Config{
		AllowOrigins:     []string{cfg.AllowedOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}
	router.Use(cors.New(corsConfig))

	router.GET("/api/health", func(c *gin.Context) {
		if err := pool.Ping(c.Request.Context()); err != nil {
			slog.Error("health check: database unreachable", "error", err)
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "error", "database": "unreachable"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok", "database": "reachable"})
	})

	slog.Info("server starting", "port", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
