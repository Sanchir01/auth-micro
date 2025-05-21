package authgrpc

import (
	"context"
	authv1 "github.com/Sanchir01/auth-proto/gen/go/auth"
	"google.golang.org/grpc"
	"log/slog"
)

type ServerApi struct {
	authv1.UnimplementedAuthServer
}

func RegisterServer(gRPC *grpc.Server) {
	authv1.RegisterAuthServer(gRPC, &ServerApi{})
}

func (s *ServerApi) Login(ctx context.Context, request *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	slog.Info("Login", "request", request)
	return &authv1.LoginResponse{
		Phone: "12345",
		Role:  authv1.Role_ADMIN,
		Title: "test",
		Email: "test@gmail.com",
	}, nil
}

func (s *ServerApi) ConfirmRegister(ctx context.Context, request *authv1.ConfirmRegisterRequest) (*authv1.ConfirmRegisterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServerApi) mustEmbedUnimplementedAuthServer() {
	//TODO implement me
	panic("implement me")
}

func (s *ServerApi) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {

	return nil, nil
}
