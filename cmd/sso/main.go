package main

import (
	"log/slog"
	"os"

	"github.com/rmntim/sso/internal/app"
	"github.com/rmntim/sso/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	config := config.MustLoad()

	logger := newLogger(config.Env)
	logger.Info("starting application", slog.Any("config", config))

	application := app.New(logger, config.Grpc.Port, config.StoragePath, config.TokenTTL)

	application.GRPCApp.MustRun()
}

func newLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return logger
}
