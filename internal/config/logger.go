package config

import (
	"log/slog"
	"os"
	"strings"
)

func ParseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func InitLogger() *slog.Logger {
	logLevelENV := os.Getenv("LOG_LEVEL")
	if logLevelENV == "" {
		logLevelENV = "info"
	}

	level := ParseLevel(logLevelENV)
	handlerLogger := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	})

	logger := slog.New(handlerLogger)

	slog.Info("logger инициализирован", "level", level.String())

	return logger
}
