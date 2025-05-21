package app

import "github.com/Sanchir01/auth-micro/internal/features/user"

type Repository struct {
	UserRepository *user.RepositoryUser
}

func NewRepository(databases *Database) *Repository {
	return &Repository{
		UserRepository: user.NewRepositoryUser(databases.PrimaryDB, databases.RedisDB),
	}
}
