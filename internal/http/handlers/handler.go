package handlers

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	Health *HealthHandler
	// Auth   *AuthHandler
	// Menu   *MenuHandler
}

func New(pool *pgxpool.Pool, jwtSecret []byte) *Handler {
	return &Handler{
		Health: NewHealthHandler(pool),
		// Auth:   NewAuthHandler(pool, jwtSecret),
	}
}
