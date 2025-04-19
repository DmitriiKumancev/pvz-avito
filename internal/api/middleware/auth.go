package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/dkumancev/avito-pvz/pkg/domain"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	userContextKey contextKey = "user"
)

func GetUserFromContext(ctx context.Context) (*domain.User, error) {
	user, ok := ctx.Value(userContextKey).(*domain.User)
	if !ok {
		return nil, errors.New("пользователь не найден в контексте")
	}
	return user, nil
}

func AuthMiddleware(jwtSecret []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"message":"Отсутствует токен авторизации"}`, http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, `{"message":"Неверный формат токена авторизации"}`, http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// Парсим токен
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("неподдерживаемый метод подписи")
			}
			return jwtSecret, nil
		})

		if err != nil {
			http.Error(w, `{"message":"Недействительный токен"}`, http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, `{"message":"Недействительный токен"}`, http.StatusUnauthorized)
			return
		}

		// Извлекаем данные пользователя из токена
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, `{"message":"Недействительный токен"}`, http.StatusUnauthorized)
			return
		}

		userID, ok := claims["id"].(string)
		if !ok {
			http.Error(w, `{"message":"Недействительный токен"}`, http.StatusUnauthorized)
			return
		}

		email, ok := claims["email"].(string)
		if !ok {
			http.Error(w, `{"message":"Недействительный токен"}`, http.StatusUnauthorized)
			return
		}

		roleStr, ok := claims["role"].(string)
		if !ok {
			http.Error(w, `{"message":"Недействительный токен"}`, http.StatusUnauthorized)
			return
		}

		var role domain.UserRole
		switch roleStr {
		case string(domain.EmployeeRole):
			role = domain.EmployeeRole
		case string(domain.ModeratorRole):
			role = domain.ModeratorRole
		default:
			http.Error(w, `{"message":"Неизвестная роль пользователя"}`, http.StatusUnauthorized)
			return
		}

		user := &domain.User{
			ID:    userID,
			Email: email,
			Role:  role,
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RoleMiddleware(allowedRoles []domain.UserRole, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := GetUserFromContext(r.Context())
		if err != nil {
			http.Error(w, `{"message":"Ошибка авторизации"}`, http.StatusUnauthorized)
			return
		}

		allowed := false
		for _, role := range allowedRoles {
			if user.Role == role {
				allowed = true
				break
			}
		}

		if !allowed {
			http.Error(w, `{"message":"Недостаточно прав для выполнения операции"}`, http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
