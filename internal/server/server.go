package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dkumancev/avito-pvz/config"
	"github.com/dkumancev/avito-pvz/internal/api"
	"github.com/dkumancev/avito-pvz/pkg/infrastructure/postgres/db"
)

type Server struct {
	httpServer *http.Server
	cfg        *config.Config
	logger     *log.Logger
}

func NewServer(cfg *config.Config) *Server {
	logger := log.New(os.Stdout, "[SERVER] ", log.LstdFlags)
	return &Server{
		cfg:    cfg,
		logger: logger,
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
	s.logger.Println("Сервис для управления ПВЗ (Среда:", s.cfg.App.Environment, ")")

	dbConn, err := db.New(s.cfg.Postgres)
	if err != nil {
		s.logger.Printf("Ошибка подключения к базе данных: %v", err)
		return err
	}
	defer dbConn.Close()

	s.logger.Println("Успешное подключение к базе данных")

	router := api.NewRouter(dbConn, []byte(s.cfg.Auth.JWTSecret))
	handler := router.Setup()

	go func() {
		s.logger.Printf("HTTP сервер запущен на порту: %s", s.cfg.HTTP.Port)
		if err := s.Run(s.cfg.HTTP.Port, handler); err != nil && err != http.ErrServerClosed {
			s.logger.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	return s.gracefulShutdown()
}

func (s *Server) gracefulShutdown() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	s.logger.Println("Завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		s.logger.Printf("Ошибка при остановке сервера: %v", err)
		return err
	}

	s.logger.Println("Сервер успешно остановлен")
	return nil
}
