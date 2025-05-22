package user

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"math/rand"
	"time"
)

type Service struct {
	RepositoryUser *RepositoryUser
	primaryDB      *pgxpool.Pool
}

func NewService(RepositoryUser *RepositoryUser, primarydb *pgxpool.Pool) *Service {
	return &Service{
		RepositoryUser: RepositoryUser,
		primaryDB:      primarydb,
	}
}

func (s *Service) UserById(ctx context.Context, id uuid.UUID) (*User, error) {
	user, err := s.RepositoryUser.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (s *Service) Login(ctx context.Context, email string, password string) (*User, error) {
	usersdb, err := s.RepositoryUser.GetByEmail(ctx, email)
	if err != nil {

		return nil, err
	}
	decodepass, err := base64.StdEncoding.DecodeString(usersdb.Password)
	if err != nil {
		return nil, fmt.Errorf("Неправильный пароль")
	}
	verifypass := VerifyPassword(
		decodepass,
		password,
	)
	if verifypass == false {
		return nil, fmt.Errorf("Неправильный пароль")
	}

	return usersdb, nil
}

func (s *Service) UserByPhone(ctx context.Context, phone string) (*User, error) {
	usersdb, err := s.RepositoryUser.GetByPhone(ctx, phone)
	if err != nil {

		return nil, err
	}

	return usersdb, nil
}

func (s *Service) Registrations(ctx context.Context, password, phone, title, email string) error {

	_, err := s.RepositoryUser.GetByPhone(ctx, phone)
	if err == nil {
		slog.Error("User with this phone already exists")
		return errors.New("user with this phone already exists")
	}

	_, err = s.RepositoryUser.GetByEmail(ctx, email)
	if err == nil {
		slog.Error("User with this email already exists")
		return errors.New("user with this email already exists")
	}

	rand.Seed(time.Now().UnixNano())

	randomNumber := rand.Intn(900000) + 100000

	if err := s.RepositoryUser.SetConfirmationCode(ctx, email, randomNumber); err != nil {
		return err
	}
	if err := SendMail(randomNumber); err != nil {
		return err
	}
	return nil
}
func (s *Service) ConfirmRegister(ctx context.Context, password, phone, title, email, code string) (*User, error) {
	conn, err := s.primaryDB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				err = errors.Join(err, rollbackErr)
				return
			}
		}
	}()
	oldcode, err := s.RepositoryUser.GetUserCodeByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if oldcode != code {
		return nil, errors.New("confirmation code is invalid")
	}
	if err := s.RepositoryUser.DeleteUserCodeByEmail(ctx, email); err != nil {
		return nil, err
	}
	hashedPassword, err := GeneratePasswordHash(password)
	if err != nil {
		return nil, err
	}
	user, err := s.RepositoryUser.CreateUser(ctx, title, phone, email, "user", hashedPassword, tx)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return user, nil
}
