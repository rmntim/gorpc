package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/rmntim/sso/internal/app"
	"github.com/rmntim/sso/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := newLogger(cfg.Env)
	log.Info("starting application", slog.Any("config", cfg))

	application := app.New(log, cfg.Grpc.Port, cfg.StoragePath, cfg.TokenTTL)

	go application.GRPCApp.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop
	log.Info("received signal", slog.String("signal", sign.String()))

	application.GRPCApp.Stop()
	log.Info("application stopped")
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
