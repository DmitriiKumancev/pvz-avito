package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dkumancev/avito-pvz/config"
	"github.com/dkumancev/avito-pvz/pkg/infrastructure/postgres"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	log.Printf("Сервис для управления ПВЗ (Среда: %s)", cfg.App.Environment)

	db, err := postgres.NewPostgresDB(cfg.Postgres)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	log.Println("Успешное подключение к базе данных")

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status": "ok", "db": "connected"}`)
	})

	server := &http.Server{
		Addr:         ":" + cfg.HTTP.Port,
		Handler:      mux,
		ReadTimeout:  cfg.HTTP.Timeout,
		WriteTimeout: cfg.HTTP.Timeout,
		IdleTimeout:  cfg.HTTP.Timeout,
	}

	go func() {
		log.Printf("HTTP сервер запущен на порту: %s", cfg.HTTP.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка при остановке сервера: %v", err)
	}

	log.Println("Сервер успешно остановлен")
}
