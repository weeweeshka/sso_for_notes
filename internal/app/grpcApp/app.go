package grpcApp

import (
	"fmt"
	"github.com/weeweeshka/sso_for_notes/internal/grpc/sso"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type GrpcApp struct {
	gRPCServer *grpc.Server
	slog       slog.Logger
	port       int
}

func New(port int, logger *slog.Logger, authService sso.Auth) *GrpcApp {
	grpcServer := grpc.NewServer()
	sso.RegisterServer(grpcServer, authService)
	return &GrpcApp{
		gRPCServer: grpcServer,
		slog:       *logger,
		port:       port,
	}
}

func (a *GrpcApp) Run() error {
	const op = "app.grpcApp.Run"
	log := a.slog.With(slog.String("op", op), slog.Int("port", a.port))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("starting gRPC server")
	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}

func (a *GrpcApp) MustRun() {
	_ = a.Run()
}

func (a *GrpcApp) GracefulStop() {
	const op = "app.grpcApp.GracefulStop"
	a.slog.With(slog.String("op", op))

	a.gRPCServer.GracefulStop()
}
