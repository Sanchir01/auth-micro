package user

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"log"
	"log/slog"
	"time"
)

type RepositoryUser struct {
	primaryDB *pgxpool.Pool
	redisDB   *redis.Client
}

func NewRepositoryUser(primaryDB *pgxpool.Pool, redisDB *redis.Client) *RepositoryUser {
	return &RepositoryUser{
		primaryDB,
		redisDB,
	}
}
func (r *RepositoryUser) GetByPhone(ctx context.Context, phone string) (*User, error) {
	conn, err := r.primaryDB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query, arg, err := sq.
		Select("id ,title, phone, created_at, updated_at, version, role, password,email").
		From("public.users").
		Where(sq.Eq{"phone": phone}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	var user DBUser
	if err := conn.QueryRow(ctx, query, arg...).Scan(
		&user.ID, &user.Title, &user.Phone, &user.CreatedAt, &user.UpdatedAt, &user.Version, &user.Role); err != nil {
		return nil, err
	}

	return &User{
		ID:        user.ID,
		Title:     user.Title,
		Phone:     user.Phone,
		Role:      user.Role,
		Email:     user.Email,
		Version:   user.Version,
		Password:  "",
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
func (r *RepositoryUser) GetByEmail(ctx context.Context, email string) (*User, error) {
	
	conn, err := r.primaryDB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query, args, err := sq.
		Select("id, title, phone, created_at, updated_at, version, role, password, email").
		From("public.users").
		Where(sq.Eq{"email": email}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var user DBUser
	err = conn.QueryRow(ctx, query, args...).Scan(
		&user.ID, &user.Title, &user.Phone, &user.CreatedAt,
		&user.UpdatedAt, &user.Version, &user.Role,
		&user.Password, &user.Email,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("Неправильный логин или пароль")
		}
		return nil, err
	}
	userJson, err := json.Marshal(user)
	if err != nil {
		slog.Error("failed cashed user", err.Error())
	}
	//todo update seconds in 2 week
	if err := r.redisDB.Set(ctx, "userbyemail", userJson, 2*time.Second).Err(); err != nil {
		log.Printf("Failed to cache candles: %v", err)
	}
	result := &User{
		ID:        user.ID,
		Title:     user.Title,
		Phone:     user.Phone,
		Role:      user.Role,
		Email:     user.Email,
		Version:   user.Version,
		Password:  base64.StdEncoding.EncodeToString(user.Password),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return result, nil
}
func (r *RepositoryUser) GetBySlug(ctx context.Context, slug string) (*User, error) {
	conn, err := r.primaryDB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query, arg, err := sq.
		Select("id ,title, phone, created_at, updated_at, version, role, password,email").
		From("public.users").
		Where(sq.Eq{"slug": slug}).
		PlaceholderFormat(sq.Dollar).ToSql()
	var user DBUser

	if err := conn.QueryRow(ctx, query, arg...).Scan(
		&user.ID, &user.Title, &user.Phone, &user.CreatedAt, &user.UpdatedAt, &user.Version, &user.Role, &user.Password, &user.Email); err != nil {
		return nil, err
	}

	return &User{
		ID:        user.ID,
		Title:     user.Title,
		Phone:     user.Phone,
		Role:      user.Role,
		Email:     user.Email,
		Version:   user.Version,
		Password:  "",
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (r *RepositoryUser) GetById(ctx context.Context, userId uuid.UUID) (*User, error) {
	conn, err := r.primaryDB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query, arg, err := sq.
		Select("id ,title, phone, created_at, updated_at, version, role, password, email").
		From("public.users").
		Where(sq.Eq{"id": userId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	var user DBUser
	if err := conn.QueryRow(ctx, query, arg...).Scan(
		&user.ID, &user.Title, &user.Phone, &user.CreatedAt, &user.UpdatedAt, &user.Version, &user.Role, &user.Password, &user.Email); err != nil {
		return nil, err
	}

	return &User{
		ID:        user.ID,
		Title:     user.Title,
		Phone:     user.Phone,
		Role:      user.Role,
		Email:     user.Email,
		Version:   user.Version,
		Password:  "",
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (r *RepositoryUser) CreateUser(ctx context.Context, title, phone, email, role string, password []byte, tx pgx.Tx) (*User, error) {
	query, arg, err := sq.
		Insert("users").
		Columns("title", "phone", "email", "role", "password").
		Values(title, phone, email, role, password).
		Suffix("RETURNING id, phone, role, email").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}
	var users DBUser

	if err := tx.QueryRow(ctx, query, arg...).Scan(&users.ID, &users.Phone, &users.Role, &users.Email); err != nil {
		if err == pgx.ErrTxCommitRollback {
			return nil, fmt.Errorf("ошибка при создании пользователя")
		}
		return nil, err
	}

	return &User{
		ID:        users.ID,
		Title:     users.Title,
		Phone:     users.Phone,
		Role:      users.Role,
		Email:     users.Email,
		Version:   users.Version,
		Password:  "",
		CreatedAt: users.CreatedAt,
		UpdatedAt: users.UpdatedAt,
	}, nil
}

func (r *RepositoryUser) SetConfirmationCode(ctx context.Context, email string, code int) error {

	if err := r.redisDB.Set(ctx, "verify:"+email, code, 2*time.Minute).Err(); err != nil {
		return err
	}

	return nil
}
func (r *RepositoryUser) GetUserCodeByEmail(ctx context.Context, email string) (string, error) {
	code, err := r.redisDB.Get(ctx, "verify:"+email).Result()
	if err != nil {
		return "", err
	}
	return code, nil
}
func (r *RepositoryUser) DeleteUserCodeByEmail(ctx context.Context, email string) error {
	_, err := r.redisDB.Del(ctx, "verify:"+email).Result()
	if err != nil {
		return err
	}
	return nil
}
