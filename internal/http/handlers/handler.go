package handlers

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/alhaos-qr-menu/api/internal/config"
)

// Handler aggregates all HTTP handler groups and their shared dependencies.
type Handler struct {
	Health *HealthHandler
	Auth   *AuthHandler
	User   *UserHandler
}

// New constructs a Handler with all sub-handlers wired up.
func New(pool *pgxpool.Pool, cfg *config.Config) *Handler {
	return &Handler{
		Health: NewHealthHandler(pool),
		Auth:   NewAuthHandler(pool, cfg),
		User:   NewUserHandler(pool),
	}
}
