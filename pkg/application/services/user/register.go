package user

import (
	"context"
	"fmt"

	"github.com/dkumancev/avito-pvz/pkg/domain"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) Register(ctx context.Context, email, password string, role domain.UserRole) (domain.User, error) {
	exists, err := s.userRepo.Exists(ctx, email)
	if err != nil {
		return domain.User{}, fmt.Errorf("ошибка проверки существования пользователя: %w", err)
	}
	if exists {
		return domain.User{}, ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, fmt.Errorf("ошибка хеширования пароля: %w", err)
	}

	user, err := domain.NewUser(email, string(hashedPassword), role)
	if err != nil {
		return domain.User{}, fmt.Errorf("ошибка создания пользователя: %w", err)
	}

	return s.userRepo.Create(ctx, user)
}
