package user

import (
	"context"
	"errors"
	"time"

	"github.com/dkumancev/avito-pvz/pkg/application/repositories"
	"github.com/dkumancev/avito-pvz/pkg/domain"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
)

type Service interface {
	// Регистрация нового пользователя
	Register(ctx context.Context, email, password string, role domain.UserRole) (domain.User, error)

	// Вход пользователя
	Login(ctx context.Context, email, password string) (string, error)

	// Тестовый вход для получения токена с заданной ролью
	DummyLogin(ctx context.Context, role domain.UserRole) (string, error)
}

type service struct {
	userRepo    repositories.UserRepository
	jwtSecret   []byte
	tokenExpiry time.Duration
}

func New(
	userRepo repositories.UserRepository,
	jwtSecret []byte,
	tokenExpiry time.Duration,
) Service {
	return &service{
		userRepo:    userRepo,
		jwtSecret:   jwtSecret,
		tokenExpiry: tokenExpiry,
	}
}

func (s *service) generateToken(user domain.User) (string, error) {
	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(s.tokenExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
