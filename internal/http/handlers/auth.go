package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/alhaos-qr-menu/api/internal/auth"
	"github.com/alhaos-qr-menu/api/internal/config"
)

// AuthHandler handles user registration, login, and token lifecycle.
type AuthHandler struct {
	db  *pgxpool.Pool
	cfg *config.Config
}

// NewAuthHandler constructs an AuthHandler.
func NewAuthHandler(db *pgxpool.Pool, cfg *config.Config) *AuthHandler {
	return &AuthHandler{db: db, cfg: cfg}
}

type credentialsRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Register creates a new user account and returns a token pair.
// POST /api/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req credentialsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not process password"})
		return
	}

	var userID int64
	err = h.db.QueryRow(c.Request.Context(),
		`INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id`,
		req.Email, hash,
	).Scan(&userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			c.JSON(http.StatusConflict, gin.H{"error": "user with this email already exists"})
			return
		}
		slog.Error("register: failed to insert user", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	h.issueTokens(c, userID)
}

// Login verifies credentials and returns a token pair.
// POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req credentialsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userID int64
	var hash string
	err := h.db.QueryRow(c.Request.Context(),
		`SELECT id, password_hash FROM users WHERE email = $1`, req.Email,
	).Scan(&userID, &hash)
	if errors.Is(err, pgx.ErrNoRows) || !auth.CheckPassword(hash, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	h.issueTokens(c, userID)
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Refresh rotates a valid refresh token for a new token pair.
// POST /api/auth/refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token is required"})
		return
	}

	hash := auth.HashRefreshToken(req.RefreshToken)

	var userID int64
	var expiresAt time.Time
	var revoked bool
	err := h.db.QueryRow(c.Request.Context(),
		`SELECT user_id, expires_at, revoked FROM refresh_tokens WHERE token_hash = $1`, hash,
	).Scan(&userID, &expiresAt, &revoked)
	if err != nil || revoked || time.Now().After(expiresAt) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token"})
		return
	}

	// Rotate: revoke the used refresh token before issuing a new pair.
	_, _ = h.db.Exec(c.Request.Context(),
		`UPDATE refresh_tokens SET revoked = true WHERE token_hash = $1`, hash)

	h.issueTokens(c, userID)
}

// Logout revokes the given refresh token.
// POST /api/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err == nil {
		hash := auth.HashRefreshToken(req.RefreshToken)
		_, _ = h.db.Exec(c.Request.Context(),
			`UPDATE refresh_tokens SET revoked = true WHERE token_hash = $1`, hash)
	}
	c.Status(http.StatusNoContent)
}

// issueTokens generates and persists a fresh access/refresh token pair
// for userID, and writes it to the response.
func (h *AuthHandler) issueTokens(c *gin.Context, userID int64) {
	access, err := auth.GenerateAccessToken([]byte(h.cfg.JWTSecret), userID, h.cfg.AccessTokenTTL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not issue access token"})
		return
	}

	rawRefresh, refreshHash, err := auth.NewRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not issue refresh token"})
		return
	}

	_, err = h.db.Exec(c.Request.Context(),
		`INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`,
		userID, refreshHash, time.Now().Add(h.cfg.RefreshTokenTTL),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not store refresh token"})
		return
	}

	c.JSON(http.StatusOK, tokenResponse{AccessToken: access, RefreshToken: rawRefresh})
}
