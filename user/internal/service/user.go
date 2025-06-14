package service

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"shopy/internal/models"
	"shopy/internal/types"
	"shopy/pkg/encrypt"
	"shopy/pkg/token"
	"strconv"
	"time"
)

type Repository interface {
	AddUser(ctx context.Context, params types.UserParams) (*models.User, error)
	DelUser(ctx context.Context, email string) error
	GetUser(ctx context.Context, email string) (*models.User, error)
}

type User struct {
	logger     *slog.Logger
	jwt        *token.JWT
	expiresAt  time.Duration
	repository Repository
}

func NewUser(logger *slog.Logger, repository Repository) *User {
	hours, err := strconv.Atoi(os.Getenv("TOKEN_EXP"))
	if err != nil {
		logger.Error("error parsing token expiration time", "error", err)
		hours = 1 // default to 1 hour
	}

	return &User{
		logger:     logger,
		repository: repository,
		expiresAt:  time.Hour * time.Duration(hours),
		jwt:        token.NewJWT(os.Getenv("TOKEN_KEY")),
	}
}

func (u *User) LoginUser(ctx context.Context, email, password string) (string, error) {
	user, err := u.repository.GetUser(ctx, email)
	if err != nil {
		u.logger.Error("error getting user", "error", err)
		return "", types.ErrUnauthorized
	}

	if !encrypt.VerifyPassword(password, user.Password) {
		return "", types.ErrUnauthorized
	}

	return u.jwt.Generate(ctx, u.expiresAt)
}

func (u *User) AddUser(ctx context.Context, params types.UserParams) (*models.User, error) {
	hash, err := encrypt.HashPassword(params.Password)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	params.Password = hash
	return u.repository.AddUser(ctx, params)
}

func (u *User) DelUser(ctx context.Context, email string) error {
	return u.repository.DelUser(ctx, email)
}
