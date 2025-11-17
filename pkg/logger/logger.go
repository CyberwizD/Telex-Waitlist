package logger

import (
	"log/slog"
	"os"
	"strings"
)

// New returns a slog.Logger with JSON output and the requested level.
// Supported levels: debug, info, warn, error. Defaults to info.
func New(level string) *slog.Logger {
	lvl := parseLevel(level)

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	})
	return slog.New(handler)
}

func parseLevel(lvl string) slog.Leveler {
	switch strings.ToLower(strings.TrimSpace(lvl)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
