package router

import (
	"github.com/gin-gonic/gin"

	"github.com/alhaos-qr-menu/api/internal/auth"
	"github.com/alhaos-qr-menu/api/internal/http/handlers"
)

func Setup(r *gin.Engine, h *handlers.Handler, jwtSecret []byte) {
	api := r.Group("/api")

	// Public
	api.GET("/health", h.Health.Health)

	// Protected
	protected := api.Group("/protected")
	protected.Use(auth.RequireAuth(jwtSecret))
	{
		protected.GET("/me", h.ProtectedExample)
	}
}
