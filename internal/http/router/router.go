package router

import (
	"github.com/alhaos-qr-menu/api/internal/auth"
	"github.com/alhaos-qr-menu/api/internal/http/handlers"
	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, h *handlers.Handler, jwtSecret []byte) {
	api := r.Group("/api")

	// Public routes
	api.GET("/health", h.Health.Health)

	// Protected routes
	protected := api.Group("/protected")
	protected.Use(auth.RequireAuth(jwtSecret))
	{
		protected.GET("/me", h.ProtectedExample)
	}
}
