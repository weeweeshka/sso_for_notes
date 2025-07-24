package buildApp

import (
	"fmt"
	runGrpc "github.com/weeweeshka/sso_for_notes/internal/app/grpcApp"
	"github.com/weeweeshka/sso_for_notes/internal/businessLogic/auth"
	"github.com/weeweeshka/sso_for_notes/internal/postgres/postgres"
	"log/slog"
	"time"
)

type Auth struct {
	GRPCServer runGrpc.GrpcApp
}

func NewApp(port int, storagePath string, slog *slog.Logger, tokenTTL time.Duration) (*Auth, error) {
	const op = "app.buildApp.NewApp"

	storage, err := postgres.NewStorage(storagePath, slog)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ssoService := auth.New(slog, storage, storage, storage, tokenTTL)
	grpcAPP := runGrpc.New(port, slog, ssoService)

	return &Auth{GRPCServer: *grpcAPP}, nil

}
