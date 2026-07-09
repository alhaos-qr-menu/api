package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/alhaos-qr-menu/api/internal/auth"
)

type ProtectedExample struct {
}

func (h *Handler) ProtectedExample(c *gin.Context) {
	userID, exists := auth.UserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id not found in context"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Это защищённый маршрут",
		"user_id": userID,
		"status":  "success",
	})
}
