package router

import (
	"github.com/gin-gonic/gin"
	"github.com/alhaos-qr-menu/api/internal/handlers"
)

func Setup(router *gin.Engine, h *handlers.Handler) {
	api := router.Group("/api")

	// Public routes
	api.GET("/health", h.Health.Health)

	// Protected routes later
	// api.Use(auth.RequireAuth(...))
	// v1 := api.Group("/v1")
}
