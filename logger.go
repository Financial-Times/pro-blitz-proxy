package main

import (
	"log/slog"
	"os"
)

func registerSlogDefaultLogger(appName string, logLevel slog.Level, loggerArgs ...any) {
	hostname, _ := os.Hostname()

	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})
	logger := slog.New(logHandler).
		With("appName", appName).
		With("hostname", hostname).
		With("pid", os.Getpid())
	if len(loggerArgs) > 0 {
		logger = logger.With(loggerArgs...)
	}

	slog.SetDefault(logger)
}
