package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
)

// RecoveryMiddleware - middleware для восстановления после паники
func RecoveryMiddleware(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// Получаем стек вызовов
					stack := debug.Stack()

					// Логируем ошибку с контекстом HTTP запроса
					log.Error("PANIC RECOVERED",
						"error", fmt.Sprintf("%v", err),
						"path", r.URL.Path,
						"method", r.Method,
						"stack", string(stack))

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error": "Internal server error"}`))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
