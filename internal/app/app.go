package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/rmntim/sso/internal/app/grpc"
)

type App struct {
	GRPCApp *grpcapp.App
}

func New(
	logger *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	// TODO: init storage

	// TODO: init auth service

	grpcApp := grpcapp.New(logger, grpcPort)

	return &App{
		GRPCApp: grpcApp,
	}
}
