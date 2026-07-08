package models

import "time"

// User represents an application user stored in the database.
// PasswordHash is intentionally excluded from JSON output (json:"-")
// so it can never leak into an API response.
type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}
