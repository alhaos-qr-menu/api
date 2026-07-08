package config

import (
	"log/slog"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

// Config holds all application settings, populated from environment
// variables (see the `env` struct tags). In production these variables
// are provided by systemd via EnvironmentFile=, no .env file is present.
type Config struct {
	Port            string        `env:"PORT" env-default:"8080"`
	DatabaseURL     string        `env:"DATABASE_URL" env-required:"true"`
	JWTSecret       string        `env:"JWT_SECRET" env-required:"true"`
	AllowedOrigin   string        `env:"ALLOWED_ORIGIN" env-default:"http://localhost:5173"`
	AccessTokenTTL  time.Duration `env:"ACCESS_TOKEN_TTL" env-default:"15m"`
	RefreshTokenTTL time.Duration `env:"REFRESH_TOKEN_TTL" env-default:"720h"`

	LogPath       string `env:"LOG_PATH" env-default:"./logs/app.log"`
	LogMaxSizeMB  int    `env:"LOG_MAX_SIZE_MB" env-default:"100"`
	LogMaxBackups int    `env:"LOG_MAX_BACKUPS" env-default:"5"`
	LogMaxAgeDays int    `env:"LOG_MAX_AGE_DAYS" env-default:"30"`
	LogCompress   bool   `env:"LOG_COMPRESS" env-default:"true"`
}

// Load reads a local .env file into the process environment (if present),
// then populates Config from environment variables.
func Load() (*Config, error) {
	// Missing .env is expected in production, where systemd supplies
	// the environment directly, so we don't treat it as fatal.
	if err := godotenv.Load(); err != nil {
		slog.Debug("no .env file found, relying on real environment variables")
	}

	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
