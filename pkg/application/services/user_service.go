package services

import (
	"context"
	"errors"
	"time"

	"github.com/dkumancev/avito-pvz/pkg/application/repositories"
	"github.com/dkumancev/avito-pvz/pkg/domain"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
)

type UserService interface {
	Register(ctx context.Context, email, password string, role domain.UserRole) (domain.User, error)
	Login(ctx context.Context, email, password string) (string, error)
	DummyLogin(ctx context.Context, role domain.UserRole) (string, error)
}

type userService struct {
	userRepo    repositories.UserRepository
	jwtSecret   []byte
	tokenExpiry time.Duration
}

func NewUserService(repo repositories.UserRepository, jwtSecret []byte, tokenExpiry time.Duration) UserService {
	return &userService{
		userRepo:    repo,
		jwtSecret:   jwtSecret,
		tokenExpiry: tokenExpiry,
	}
}

func (s *userService) Register(ctx context.Context, email, password string, role domain.UserRole) (domain.User, error) {
	exists, err := s.userRepo.Exists(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	if exists {
		return domain.User{}, ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, err
	}

	user, err := domain.NewUser(email, string(hashedPassword), role)
	if err != nil {
		return domain.User{}, err
	}

	return s.userRepo.Create(ctx, user)
}

func (s *userService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", ErrInvalidCredentials
	}

	// Generate JWT token
	return s.generateToken(user)
}

func (s *userService) DummyLogin(ctx context.Context, role domain.UserRole) (string, error) {
	// For testing purposes only - creates a dummy token with the given role
	dummyUser := domain.User{
		ID:    "dummy-id",
		Email: "dummy@example.com",
		Role:  role,
	}

	return s.generateToken(dummyUser)
}

func (s *userService) generateToken(user domain.User) (string, error) {
	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(s.tokenExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
