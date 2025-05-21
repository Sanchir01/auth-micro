package app

import "github.com/Sanchir01/auth-micro/internal/features/user"

type Services struct {
	UserService *user.Service
}

func NewServices(repos *Repository) *Services {
	return &Services{
		UserService: user.NewService(repos.UserRepository),
	}
}
