package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

const (
	envFormat = "LOG_FORMAT"
	envLevel  = "LOG_LEVEL"

	defaultFormat = "text"
	defaultLevel  = "info"
)

func Init() {
	format := strings.ToLower(os.Getenv(envFormat))
	if format == "" {
		format = defaultFormat
	}

	levelStr := strings.ToLower(os.Getenv(envLevel))
	if levelStr == "" {
		levelStr = defaultLevel
	}

	var level slog.Level
	switch levelStr {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: level}

	var handler slog.Handler
	switch format {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	default:
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	slog.SetDefault(slog.New(handler))
}

func Discard() {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
}
