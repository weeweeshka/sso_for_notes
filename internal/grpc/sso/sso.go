package sso

import (
	"context"
	authGrpc "github.com/weeweeshka/notes_auth/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	RegisterNewUser(ctx context.Context, email string, password string) (int64, error)
	Login(ctx context.Context, email string, password string, appID int) (string, error)
}

type serverApi struct {
	authGrpc.UnimplementedNoteAuthServer
	auth Auth
}

func RegisterServer(server *grpc.Server, auth Auth) {
	authGrpc.RegisterNoteAuthServer(server, &serverApi{auth: auth})
}

func (s *serverApi) Register(ctx context.Context, req *authGrpc.RegisterRequest) (*authGrpc.RegisterResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "Email is required")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "Password is required")
	}

	id, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &authGrpc.RegisterResponse{Id: id}, nil
}

func (s *serverApi) Login(ctx context.Context, req *authGrpc.LoginRequest) (*authGrpc.LoginResponse, error) {

	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "Email is required")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "Password is required")
	}

	if req.GetAppId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "App ID is required")
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &authGrpc.LoginResponse{Token: token}, nil
}
