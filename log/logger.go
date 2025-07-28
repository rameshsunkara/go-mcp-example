package log

import (
	"log/slog"
	"os"
	"strings"
)

// New creates a new logger with specified level and format.
func New(level string, useTextFormat bool) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: ParseLevel(level),
	}

	var handler slog.Handler
	if useTextFormat {
		handler = slog.NewTextHandler(os.Stderr, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stderr, opts)
	}

	return slog.New(handler)
}

// ParseLevel parses a string into a slog.Level.
// It supports "debug", "info", "warn", "error" (case-insensitive) and defaults to INFO.
func ParseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo // Default to INFO if unknown level
	}
}
