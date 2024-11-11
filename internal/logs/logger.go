package logs

import (
	"log/slog"
	"os"
)

var logLevel = os.Getenv("LOG_LEVEL")

func SetLogger() {
	levels := map[string]slog.Level{
		"DEBUG":  slog.LevelDebug,
		"WARN": slog.LevelWarn,
		"INFO": slog.LevelInfo,
		"ERROR": slog.LevelError,
	}

	if logLevel == "" {
		logLevel = "INFO"
	}

	opts := &slog.HandlerOptions{
		AddSource: true,
        Level: levels[logLevel],
		
    }
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(logger)
	slogLevel := new(slog.LevelVar)
	slogLevel.UnmarshalText([]byte(logLevel))
}