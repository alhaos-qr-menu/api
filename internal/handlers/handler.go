package handlers

import "github.com/jackc/pgx/v5/pgxpool"

type Handler struct {
	Health *HealthHandler
	// Auth   *AuthHandler
	// Menu   *MenuHandler
}

func New(db *pgxpool.Pool /*, другие зависимости */) *Handler {
	return &Handler{
		Health: NewHealthHandler(db),
	}
}
