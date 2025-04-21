package main

import (
	"log"

	"github.com/dkumancev/avito-pvz/config"
	"github.com/dkumancev/avito-pvz/internal/server"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	srv := server.NewServer(cfg)
	if err := srv.Start(); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
