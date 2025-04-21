package middleware

import (
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

func LoggerMiddleware(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrapped := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			path := r.URL.Path
			query := sanitizeQuery(r.URL.RawQuery)

			log.Debug("Получен HTTP запрос",
				"method", r.Method,
				"path", path,
				"query", query,
				"remote_addr", r.RemoteAddr,
				"user_agent", r.UserAgent(),
			)

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)

			// уровень логирования в зависимости от статус-кода
			var logFunc func(msg string, args ...any)
			switch {
			case wrapped.statusCode >= 500:
				logFunc = log.Error
			case wrapped.statusCode >= 400:
				logFunc = log.Warn
			default:
				logFunc = log.Info
			}

			// Логируем результат запроса
			logFunc("Завершен HTTP запрос",
				"method", r.Method,
				"path", path,
				"query", query,
				"status", wrapped.statusCode,
				"size", wrapped.size,
				"duration", duration.String(),
				"duration_ms", float64(duration.Microseconds())/1000.0,
			)
		})
	}
}

func sanitizeQuery(query string) string {
	if query == "" {
		return ""
	}

	params := strings.Split(query, "&")
	for i, param := range params {
		parts := strings.SplitN(param, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		// Маскирование чувствительных параметров
		if isConfidentialParam(key) {
			params[i] = key + "=[REDACTED]"
		}
	}

	return strings.Join(params, "&")
}

func isConfidentialParam(param string) bool {
	param = strings.ToLower(param)

	sensitiveParams := []string{
		"password", "token", "jwt", "secret", "key", "auth",
		"apikey", "access_token", "refresh_token", "authorization",
		"пароль", "токен", "ключ",
	}

	for _, sensitive := range sensitiveParams {
		if strings.Contains(param, sensitive) {
			return true
		}
	}

	return false
}
