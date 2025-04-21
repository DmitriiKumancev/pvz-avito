package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dkumancev/avito-pvz/pkg/application/services"
	userServices "github.com/dkumancev/avito-pvz/pkg/application/services/user"
	"github.com/dkumancev/avito-pvz/pkg/domain"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidCredentials = userServices.ErrInvalidCredentials
	ErrUserAlreadyExists  = userServices.ErrUserAlreadyExists
)

type MockUserRepository struct {
	users map[string]domain.User // Хранение пользователей по email
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]domain.User),
	}
}

func (m *MockUserRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	if _, exists := m.users[user.Email]; exists {
		return domain.User{}, errors.New("user already exists")
	}

	// фиксированный ID для предсказуемости тестов
	if user.ID == "" {
		user.ID = "mock-user-id"
	}

	m.users[user.Email] = user
	return user, nil
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	user, exists := m.users[email]
	if !exists {
		return domain.User{}, errors.New("user not found")
	}
	return user, nil
}

func (m *MockUserRepository) Exists(ctx context.Context, email string) (bool, error) {
	_, exists := m.users[email]
	return exists, nil
}

// Тесты

func TestUserService_Register_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := NewMockUserRepository()
	jwtSecret := []byte("test-secret")
	tokenExpiry := 24 * time.Hour

	service := services.NewUserService(mockRepo, jwtSecret, tokenExpiry)

	email := "test@example.com"
	password := "password123"
	role := domain.EmployeeRole

	user, err := service.Register(ctx, email, password, role)

	if err != nil {
		t.Errorf("Ожидалось отсутствие ошибки, получено: %v", err)
	}

	if user.Email != email {
		t.Errorf("Ожидался email %s, получен %s", email, user.Email)
	}

	if user.Role != role {
		t.Errorf("Ожидалась роль %s, получена %s", role, user.Role)
	}

	if user.ID == "" {
		t.Error("ID пользователя не должен быть пустым")
	}

	if user.PasswordHash == password {
		t.Error("Пароль не должен храниться в открытом виде")
	}
}

func TestUserService_Register_UserAlreadyExists(t *testing.T) {
	ctx := context.Background()
	mockRepo := NewMockUserRepository()
	jwtSecret := []byte("test-secret")
	tokenExpiry := 24 * time.Hour

	service := services.NewUserService(mockRepo, jwtSecret, tokenExpiry)

	email := "existing@example.com"
	password := "password123"
	role := domain.EmployeeRole

	_, _ = service.Register(ctx, email, password, role)

	// Пробуем создать пользователя с тем же email
	user, err := service.Register(ctx, email, "different-password", role)

	if err == nil {
		t.Error("Ожидалась ошибка при создании пользователя с существующим email")
	}

	if err != ErrUserAlreadyExists {
		t.Errorf("Ожидалась ошибка %v, получена %v", ErrUserAlreadyExists, err)
	}

	if user != (domain.User{}) {
		t.Error("Возвращенный пользователь должен быть пустым")
	}
}

func TestUserService_Register_InvalidData(t *testing.T) {
	ctx := context.Background()
	mockRepo := NewMockUserRepository()
	jwtSecret := []byte("test-secret")
	tokenExpiry := 24 * time.Hour

	service := services.NewUserService(mockRepo, jwtSecret, tokenExpiry)

	testCases := []struct {
		name     string
		email    string
		password string
		role     domain.UserRole
	}{
		{
			name:     "Неверный email",
			email:    "invalid-email",
			password: "password123",
			role:     domain.EmployeeRole,
		},
		{
			name:     "Неверная роль",
			email:    "test@example.com",
			password: "password123",
			role:     "invalid-role",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user, err := service.Register(ctx, tc.email, tc.password, tc.role)

			if err == nil {
				t.Errorf("Ожидалась ошибка при проверке %s, но ошибки не возникло", tc.name)
			}

			if user != (domain.User{}) {
				t.Error("Возвращенный пользователь должен быть пустым")
			}
		})
	}
}

// Отдельный тест для пустого пароля
func TestUserService_Register_EmptyPassword(t *testing.T) {
	_, err := domain.NewUser("test@example.com", "", domain.EmployeeRole)

	if err == nil {
		t.Error("Ожидалась ошибка при создании пользователя с пустым паролем")
	}
}

func TestUserService_Login_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := NewMockUserRepository()
	jwtSecret := []byte("test-secret")
	tokenExpiry := 24 * time.Hour

	service := services.NewUserService(mockRepo, jwtSecret, tokenExpiry)

	email := "test@example.com"
	password := "password123"
	role := domain.EmployeeRole

	createdUser, _ := service.Register(ctx, email, password, role)

	// logining
	token, err := service.Login(ctx, email, password)

	// Проверки
	if err != nil {
		t.Errorf("Ожидалось отсутствие ошибки, получено: %v", err)
	}

	if token == "" {
		t.Error("Ожидался непустой JWT токен")
	}

	// Проверка валидности токена
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("неожиданный метод подписи")
		}
		return jwtSecret, nil
	})

	if err != nil {
		t.Errorf("Ошибка при парсинге токена: %v", err)
	}

	if !parsedToken.Valid {
		t.Error("Токен должен быть валидным")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Error("Не удалось получить claims из токена")
	}

	if claims["email"] != email {
		t.Errorf("В токене ожидался email %s, получен %v", email, claims["email"])
	}

	if claims["id"] != createdUser.ID {
		t.Errorf("В токене ожидался id %s, получен %v", createdUser.ID, claims["id"])
	}

	if claims["role"] != string(role) {
		t.Errorf("В токене ожидалась роль %s, получена %v", role, claims["role"])
	}
}

func TestUserService_Login_InvalidCredentials(t *testing.T) {
	ctx := context.Background()
	mockRepo := NewMockUserRepository()
	jwtSecret := []byte("test-secret")
	tokenExpiry := 24 * time.Hour

	service := services.NewUserService(mockRepo, jwtSecret, tokenExpiry)

	email := "test@example.com"
	password := "password123"
	role := domain.EmployeeRole

	_, _ = service.Register(ctx, email, password, role)

	testCases := []struct {
		name          string
		loginEmail    string
		loginPassword string
	}{
		{
			name:          "Неверный email",
			loginEmail:    "wrong@example.com",
			loginPassword: password,
		},
		{
			name:          "Неверный пароль",
			loginEmail:    email,
			loginPassword: "wrong-password",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			token, err := service.Login(ctx, tc.loginEmail, tc.loginPassword)

			if err == nil {
				t.Errorf("Ожидалась ошибка при тесте %s, но ошибки не возникло", tc.name)
			}

			if err != ErrInvalidCredentials {
				t.Errorf("Ожидалась ошибка %v, получена %v", ErrInvalidCredentials, err)
			}

			if token != "" {
				t.Error("Токен должен быть пустым при неверных учетных данных")
			}
		})
	}
}

func TestUserService_DummyLogin(t *testing.T) {
	ctx := context.Background()
	mockRepo := NewMockUserRepository()
	jwtSecret := []byte("test-secret")
	tokenExpiry := 24 * time.Hour

	service := services.NewUserService(mockRepo, jwtSecret, tokenExpiry)

	testCases := []struct {
		name string
		role domain.UserRole
	}{
		{
			name: "Роль сотрудника",
			role: domain.EmployeeRole,
		},
		{
			name: "Роль модератора",
			role: domain.ModeratorRole,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			token, err := service.DummyLogin(ctx, tc.role)

			if err != nil {
				t.Errorf("Ожидалось отсутствие ошибки, получено: %v", err)
			}

			if token == "" {
				t.Error("Ожидался непустой JWT токен")
			}

			// Проверка валидности токена
			parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("неожиданный метод подписи")
				}
				return jwtSecret, nil
			})

			if err != nil {
				t.Errorf("Ошибка при парсинге токена: %v", err)
			}

			if !parsedToken.Valid {
				t.Error("Токен должен быть валидным")
			}

			claims, ok := parsedToken.Claims.(jwt.MapClaims)
			if !ok {
				t.Error("Не удалось получить claims из токена")
			}

			if claims["email"] != "dummy@example.com" {
				t.Errorf("В токене ожидался email %s, получен %v", "dummy@example.com", claims["email"])
			}

			if claims["id"] != "dummy-id" {
				t.Errorf("В токене ожидался id %s, получен %v", "dummy-id", claims["id"])
			}

			if claims["role"] != string(tc.role) {
				t.Errorf("В токене ожидалась роль %s, получена %v", tc.role, claims["role"])
			}
		})
	}
}
