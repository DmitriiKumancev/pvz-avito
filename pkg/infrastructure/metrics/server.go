package metrics

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	httpServer *http.Server
	port       string
	logger     *slog.Logger
}

func NewServer(port string, logger *slog.Logger) *Server {
	return &Server{
		port:   port,
		logger: logger,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	s.httpServer = &http.Server{
		Addr:              ":" + s.port,
		Handler:           mux,
		ReadHeaderTimeout: 3 * time.Second,
	}

	s.logger.Info("Запуск сервера метрик", "port", s.port)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("Ошибка при запуске сервера метрик", "error", err)
		}
	}()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Остановка сервера метрик")
	return s.httpServer.Shutdown(ctx)
}
