package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dkumancev/avito-pvz/config"
	"github.com/dkumancev/avito-pvz/internal/api"
	"github.com/dkumancev/avito-pvz/pkg/infrastructure/logger"
	"github.com/dkumancev/avito-pvz/pkg/infrastructure/metrics"
	"github.com/dkumancev/avito-pvz/pkg/infrastructure/postgres/db"
)

type Server struct {
	httpServer    *http.Server
	cfg           *config.Config
	logger        *slog.Logger
	metricsServer *metrics.Server
}

func NewServer(cfg *config.Config) *Server {
	appLogger := logger.NewLoggerFromEnvironment(cfg.App.LogLevel)
	return &Server{
		cfg:    cfg,
		logger: appLogger,
	}
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		MaxHeaderBytes:    1 << 20,
		ReadHeaderTimeout: s.cfg.HTTP.Timeout,
		WriteTimeout:      s.cfg.HTTP.Timeout,
		IdleTimeout:       s.cfg.HTTP.Timeout * 3, // обычно IdleTimeout в 2-3 раза больше
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) Start() error {
	s.logger.Info("Запуск сервиса для управления ПВЗ",
		"environment", s.cfg.App.Environment,
		"version", "1.0.0")

	// Запуск сервера метрик
	s.metricsServer = metrics.NewServer(s.cfg.Metrics.Port, s.logger)
	if err := s.metricsServer.Start(); err != nil {
		s.logger.Error("Ошибка запуска сервера метрик",
			"error", logger.SanitizeError(err),
			"port", s.cfg.Metrics.Port)
		return err
	}

	dbConn, err := db.New(s.cfg.Postgres)
	if err != nil {
		s.logger.Error("Ошибка подключения к базе данных",
			"error", logger.SanitizeError(err),
			"host", s.cfg.Postgres.Host,
			"port", s.cfg.Postgres.Port,
			"database", s.cfg.Postgres.DBName)
		return err
	}
	defer func() {
		s.logger.Debug("Закрытие соединения с базой данных")
		dbConn.Close()
	}()

	s.logger.Info("Успешное подключение к базе данных",
		"host", s.cfg.Postgres.Host,
		"port", s.cfg.Postgres.Port,
		"database", s.cfg.Postgres.DBName)

	router := api.NewRouter(dbConn, []byte(s.cfg.Auth.JWTSecret))
	handler := router.Setup()

	go func() {
		s.logger.Info("HTTP сервер запущен",
			"port", s.cfg.HTTP.Port,
			"timeout", s.cfg.HTTP.Timeout)

		if err := s.Run(s.cfg.HTTP.Port, handler); err != nil && err != http.ErrServerClosed {
			s.logger.Error("Критическая ошибка запуска сервера",
				"error", logger.SanitizeError(err))
			os.Exit(1)
		}
	}()

	return s.gracefulShutdown()
}

func (s *Server) gracefulShutdown() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	s.logger.Info("Получен сигнал завершения", "signal", sig.String())
	s.logger.Info("Начинаем корректное завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if s.metricsServer != nil {
		if err := s.metricsServer.Stop(ctx); err != nil {
			s.logger.Error("Ошибка при остановке сервера метрик",
				"error", logger.SanitizeError(err))
		}
	}

	if err := s.Shutdown(ctx); err != nil {
		s.logger.Error("Ошибка при остановке сервера",
			"error", logger.SanitizeError(err),
			"timeout", "10s")
		return err
	}

	s.logger.Info("Сервер успешно остановлен")
	return nil
}
