package authgrpc

import (
	"context"
	"github.com/Sanchir01/auth-micro/internal/features/user"
	authv1 "github.com/Sanchir01/auth-proto/gen/go/auth"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"runtime/debug"
	"strings"
)

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=1"`
}

type RegisterInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=1"`
	Phone    string `json:"phone" validate:"required"`
	Title    string `json:"title" validate:"required,min=1"`
}

type ConfirmRegisterInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=1"`
	Phone    string `json:"phone" validate:"required"`
	Title    string `json:"title" validate:"required,min=1"`
	Code     string `json:"code" validate:"required,min=1"`
}
type ServerApi struct {
	authv1.UnimplementedAuthServer
	auth Auth
}
type Auth interface {
	Login(ctx context.Context, email string, password string) (*user.User, error)
	Registrations(ctx context.Context, password, phone, title, email string) error
	ConfirmRegister(ctx context.Context, password, phone, title, email, code string) (*user.User, error)
	UserById(ctx context.Context, id uuid.UUID) (*user.User, error)
}

func RegisterServer(gRPC *grpc.Server, auth Auth) {
	authv1.RegisterAuthServer(gRPC, &ServerApi{auth: auth})
}
func (s *ServerApi) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	input := LoginInput{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
	var validate = validator.New()
	if err := validate.Struct(input); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid input: %v", err.Error())
	}
	user, err := s.auth.Login(ctx, input.Email, input.Password)
	if err != nil {
		slog.Error("Login error", "error", err)
		return nil, status.Errorf(codes.Internal, "invalid input: %v", err)
	}
	slog.Warn("Login", "user", user)
	return &authv1.LoginResponse{
		Phone: user.Phone,
		Role:  mapRole(user.Role),
		Title: user.Title,
		Email: user.Email,
	}, nil
}

func (s *ServerApi) Register(ctx context.Context, request *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	input := RegisterInput{
		Email:    request.GetEmail(),
		Phone:    request.GetPhone(),
		Title:    request.GetTitle(),
		Password: request.GetPassword(),
	}
	defer func() {
		if r := recover(); r != nil {
			slog.Error("panic recovered in ConfirmRegister", "reason", r, "stack", string(debug.Stack()))
		}
	}()
	if err := validator.New().Struct(input); err != nil {
		slog.Error("invalid input", "error", err)
		return nil, status.Errorf(codes.InvalidArgument, "invalid input: %v", err.Error())
	}
	if err := s.auth.Registrations(ctx, input.Password, input.Phone, input.Title, input.Email); err != nil {
		slog.Error("registration failed", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to register: %v", err.Error())
	}
	return &authv1.RegisterResponse{
		Ok: "success",
	}, nil
}
func (s *ServerApi) ConfirmRegister(ctx context.Context, request *authv1.ConfirmRegisterRequest) (*authv1.ConfirmRegisterResponse, error) {

	input := ConfirmRegisterInput{
		Email:    request.GetEmail(),
		Phone:    request.GetPhone(),
		Title:    request.GetTitle(),
		Password: request.GetPassword(),
		Code:     request.GetCode(),
	}

	if err := validator.New().Struct(input); err != nil {
		slog.Error("invalid input", "error", err)
		return nil, status.Errorf(codes.InvalidArgument, "invalid input: %v", err.Error())
	}

	confirmRegister, err := s.auth.ConfirmRegister(ctx, input.Password, input.Phone, input.Title, input.Email, input.Code)
	if err != nil {
		slog.Error("registration failed", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to register: %v", err.Error())
	}

	return &authv1.ConfirmRegisterResponse{
		Phone: confirmRegister.Phone,
		Role:  authv1.Role_USER,
		Title: confirmRegister.Title,
		Email: confirmRegister.Email,
	}, nil
}
func mapRole(role string) authv1.Role {
	switch strings.ToLower(role) {
	case "admin":
		return authv1.Role_ADMIN
	case "user":
		return authv1.Role_USER
	default:
		return authv1.Role_ROLE_UNSPECIFIED
	}
}
