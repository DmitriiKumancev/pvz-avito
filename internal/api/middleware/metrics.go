package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/dkumancev/avito-pvz/pkg/infrastructure/metrics"
)

func MetricsMiddleware(m *metrics.HTTPMetrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrapped := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start).Seconds()

			path := normalizePath(r.URL.Path)

			status := strconv.Itoa(wrapped.statusCode)

			m.RequestsTotal.WithLabelValues(r.Method, path, status).Inc()

			m.RequestDuration.WithLabelValues(r.Method, path, status).Observe(duration)
		})
	}
}

func normalizePath(path string) string {
	// простая реализация, которую можно расширить для конкретных маршрутов
	return path
}
