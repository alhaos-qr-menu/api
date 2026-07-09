package router

import (
	"github.com/gin-gonic/gin"

	"github.com/alhaos-qr-menu/api/internal/auth"
	"github.com/alhaos-qr-menu/api/internal/http/handlers"
)

// Setup registers all API routes on the given gin engine.
func Setup(r *gin.Engine, h *handlers.Handler, jwtSecret []byte) {
	api := r.Group("/api")

	// Public routes.
	api.GET("/health", h.Health.Health)
	api.POST("/auth/register", h.Auth.Register)
	api.POST("/auth/login", h.Auth.Login)
	api.POST("/auth/refresh", h.Auth.Refresh)
	api.POST("/auth/logout", h.Auth.Logout)

	// Routes requiring a valid access token.
	protected := api.Group("")
	protected.Use(auth.RequireAuth(jwtSecret))
	{
		protected.GET("/me", h.User.Me)
	}
}
