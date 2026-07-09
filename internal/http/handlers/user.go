package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/alhaos-qr-menu/api/internal/auth"
	"github.com/alhaos-qr-menu/api/internal/models"
)

// UserHandler handles requests about the currently authenticated user.
type UserHandler struct {
	db *pgxpool.Pool
}

// NewUserHandler constructs a UserHandler.
func NewUserHandler(db *pgxpool.Pool) *UserHandler {
	return &UserHandler{db: db}
}

// Me returns the profile of the currently authenticated user.
// GET /api/me (requires a valid access token, see auth.RequireAuth)
func (h *UserHandler) Me(c *gin.Context) {
	userID, ok := auth.UserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	var u models.User
	err := h.db.QueryRow(c.Request.Context(),
		`SELECT id, email, created_at FROM users WHERE id = $1`, userID,
	).Scan(&u.ID, &u.Email, &u.CreatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, u)
}
