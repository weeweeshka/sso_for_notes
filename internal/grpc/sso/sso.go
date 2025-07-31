package sso

import (
	"context"
	authGrpc "github.com/weeweeshka/auth_proto/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	RegisterNewUser(ctx context.Context, email string, password string) (int64, error)
	Login(ctx context.Context, email string, password string, appID int32) (string, error)
	RegisterApp(ctx context.Context, appName string, secret string) (int32, error)
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

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), req.GetAppId())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &authGrpc.LoginResponse{Token: token}, nil
}

func (s *serverApi) RegisterApp(ctx context.Context, req *authGrpc.AppRequest) (*authGrpc.AppResponse, error) {

	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "App name is required")
	}
	if req.GetSecret() == "" {
		return nil, status.Error(codes.InvalidArgument, "Secret is required")
	}

	appID, err := s.auth.RegisterApp(ctx, req.GetName(), req.GetSecret())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authGrpc.AppResponse{AppId: appID}, nil
}
