package authgrpc

import (
	"context"
	"fmt"

	"github.com/Sanchir01/auth-micro/internal/features/user"
	authv1 "github.com/Sanchir01/auth-proto/gen/go/auth"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LoginInput struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=1"`
}
type RegisterInput struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=1"`
	Phone    string `validate:"required,phone"`
	Title    string `validate:"required,min=1"`
}
type ServerApi struct {
	authv1.UnimplementedAuthServer
	auth Auth
}
type Auth interface {
	Login(ctx context.Context, email string, password string) (*user.User, error)
	Registrations(ctx context.Context, password, phone, title, email string) error
	ConfirmRegister(ctx context.Context, password, phone, title, email, code string, tx pgx.Tx) (*user.User, error)
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
		return nil, status.Errorf(codes.Internal, "invalid input: %v", err.Error())
	}
	fmt.Println(user)
	return &authv1.LoginResponse{
		Phone: user.Phone,
		Role:  authv1.Role_USER,
		Title: user.Title,
		Email: user.Email,
	}, nil
}

func (s *ServerApi) ConfirmRegister(ctx context.Context, request *authv1.ConfirmRegisterRequest) (*authv1.ConfirmRegisterResponse, error) {

	if err := s.auth.Registrations(ctx, request.GetPassword(), request.GetPhone(), request.GetTitle(), request.GetEmail()); err != nil {
		return nil, status.Errorf(codes.Internal, "invalid input: %v", err.Error())
	}

	return &authv1.ConfirmRegisterResponse{}, nil
}

func (s *ServerApi) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	return nil, nil
}
