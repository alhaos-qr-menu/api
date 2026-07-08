package logging

import (
	"io"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Config controls log file rotation behaviour.
type Config struct {
	Path       string
	MaxSizeMB  int
	MaxBackups int
	MaxAgeDays int
	Compress   bool
}

// New creates a JSON slog.Logger that writes to a rotating log file.
// lumberjack takes care of creating the target directory, rotating by
// size, and pruning old/compressed backups.
func New(cfg Config) *slog.Logger {
	rotator := &lumberjack.Logger{
		Filename:   cfg.Path,
		MaxSize:    cfg.MaxSizeMB,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAgeDays,
		Compress:   cfg.Compress,
	}

	mw := io.MultiWriter(os.Stdout, rotator)

	handler := slog.NewJSONHandler(mw, nil)
	return slog.New(handler)
}
