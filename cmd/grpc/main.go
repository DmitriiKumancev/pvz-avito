package main

import (
	"log"

	"github.com/dkumancev/avito-pvz/config"
	"github.com/dkumancev/avito-pvz/pkg/application/services/pvz"
	"github.com/dkumancev/avito-pvz/pkg/infrastructure/grpc/server"
	"github.com/dkumancev/avito-pvz/pkg/infrastructure/postgres/db"
	pgzvrepository "github.com/dkumancev/avito-pvz/pkg/infrastructure/postgres/pvz"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	dbConn, err := db.New(cfg.Postgres)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer dbConn.Close()

	pvzRepo := pgzvrepository.NewRepository(dbConn)

	pvzService := pvz.New(pvzRepo)

	port := 50051
	grpcServer := server.NewGRPCServer(pvzService, port)
	log.Printf("Запуск gRPC сервера на порту %d...", port)
	if err := grpcServer.Start(); err != nil {
		log.Fatalf("Ошибка запуска gRPC сервера: %v", err)
	}
}
 