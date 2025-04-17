package domain

import (
	"strings"
	"testing"
	"time"
)

func TestNewUser_ValidInputs(t *testing.T) {
	// Проверка валидных вариантов
	testCases := []struct {
		name         string
		email        string
		passwordHash string
		role         UserRole
	}{
		{
			name:         "Сотрудник с обычным email",
			email:        "employee@example.com",
			passwordHash: "hashed_password",
			role:         EmployeeRole,
		},
		{
			name:         "Модератор с обычным email",
			email:        "moderator@example.com",
			passwordHash: "hashed_password",
			role:         ModeratorRole,
		},
		{
			name:         "Email с поддоменом",
			email:        "user@sub.example.com",
			passwordHash: "hashed_password",
			role:         EmployeeRole,
		},
		{
			name:         "Email с плюсом",
			email:        "user+tag@example.com",
			passwordHash: "hashed_password",
			role:         EmployeeRole,
		},
		{
			name:         "Email с точкой в имени",
			email:        "first.last@example.com",
			passwordHash: "hashed_password",
			role:         EmployeeRole,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user, err := NewUser(tc.email, tc.passwordHash, tc.role)

			if err != nil {
				t.Errorf("Ожидалось отсутствие ошибки для валидных данных, получено: %v", err)
			}

			// Проверка полей
			if user.Email != tc.email {
				t.Errorf("Ожидался email %s, получен %s", tc.email, user.Email)
			}

			if user.PasswordHash != tc.passwordHash {
				t.Errorf("Ожидался passwordHash %s, получен %s", tc.passwordHash, user.PasswordHash)
			}

			if user.Role != tc.role {
				t.Errorf("Ожидалась роль %s, получена %s", tc.role, user.Role)
			}

			// Проверка времени создания
			if user.CreatedAt.IsZero() {
				t.Error("Ожидалось заполненное время создания, получено пустое значение")
			}

			if time.Since(user.CreatedAt) > time.Second {
				t.Error("Ожидалось время создания близкое к текущему времени")
			}
		})
	}
}

func TestNewUser_InvalidEmail(t *testing.T) {
	invalidEmails := []string{
		"plainaddress",           // email без @
		"@missingusername.com",   // отсутствует имя пользователя
		"user@.com",              // отсутствует домен
		"user@domain",            // отсутствует домен верхнего уровня
		"user@domain.",           // неполный домен верхнего уровня
		"",                       // пустая строка
		"user@domain@domain.com", // несколько символов @
		" user@domain.com",       // начинается с пробела
		"user@domain.com ",       // заканчивается пробелом
	}

	passwordHash := "hashed_password"
	role := EmployeeRole

	for _, email := range invalidEmails {
		t.Run(email, func(t *testing.T) {
			user, err := NewUser(email, passwordHash, role)

			if err == nil {
				t.Errorf("Ожидалась ошибка для невалидного email %q, но ее не было", email)
			}

			if err != nil && !strings.Contains(err.Error(), "invalid email format") {
				t.Errorf("Ожидалась ошибка с текстом 'invalid email format', получено: %v", err)
			}

			if user != (User{}) {
				t.Errorf("Ожидался пустой пользователь для невалидного email, получено: %+v", user)
			}
		})
	}
}

func TestNewUser_EmptyPasswordHash(t *testing.T) {
	email := "valid@example.com"
	passwordHash := ""
	role := EmployeeRole

	user, err := NewUser(email, passwordHash, role)

	if err == nil {
		t.Error("Ожидалась ошибка для пустого пароля, но ее не было")
	}

	if err != nil && !strings.Contains(err.Error(), "password hash cannot be empty") {
		t.Errorf("Ожидалась ошибка с текстом 'password hash cannot be empty', получено: %v", err)
	}

	if user != (User{}) {
		t.Errorf("Ожидался пустой пользователь для пустого пароля, получено: %+v", user)
	}
}

func TestNewUser_InvalidRole(t *testing.T) {
	email := "valid@example.com"
	passwordHash := "hashed_password"
	invalidRoles := []UserRole{
		"",
		"admin",
		"superuser",
		"guest",
	}

	for _, role := range invalidRoles {
		t.Run(string(role), func(t *testing.T) {
			user, err := NewUser(email, passwordHash, role)

			if err == nil {
				t.Errorf("Ожидалась ошибка для невалидной роли %q, но ее не было", role)
			}

			if err != nil && !strings.Contains(err.Error(), "invalid role") {
				t.Errorf("Ожидалась ошибка с текстом 'invalid role', получено: %v", err)
			}

			if user != (User{}) {
				t.Errorf("Ожидался пустой пользователь для невалидной роли, получено: %+v", user)
			}
		})
	}
}
