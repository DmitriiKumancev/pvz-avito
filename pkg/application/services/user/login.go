package user

import (
	"context"

	"github.com/dkumancev/avito-pvz/pkg/domain"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", ErrInvalidCredentials
	}

	// Генерируем JWT токен
	return s.generateToken(user)
}

// DummyLogin создает тестовый токен с указанной ролью
func (s *service) DummyLogin(ctx context.Context, role domain.UserRole) (string, error) {
	//  фиктивный пользователь для тестирования
	dummyUser := domain.User{
		ID:    "dummy-id",
		Email: "dummy@example.com",
		Role:  role,
	}

	return s.generateToken(dummyUser)
}
