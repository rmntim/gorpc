package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/rmntim/sso/internal/app/grpc"
	"github.com/rmntim/sso/internal/services/auth"
	"github.com/rmntim/sso/internal/storage/postgres"
)

type App struct {
	GRPCApp *grpcapp.App
}

func New(
	logger *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) (*App, error) {
	storage, err := postgres.New(storagePath)
	if err != nil {
		return nil, err
	}

	if err := storage.Migrate(); err != nil {
		return nil, err
	}

	authService := auth.New(logger, storage, storage, storage, tokenTTL)
	grpcApp := grpcapp.New(logger, authService, grpcPort)

	return &App{
		GRPCApp: grpcApp,
	}, nil
}
