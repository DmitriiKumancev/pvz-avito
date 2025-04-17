.PHONY: build test migrate-up migrate-down migrate-status migrate-create migrate-reset help

# Загрузка переменных окружения из .env файла
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Переменные
DB_URL ?= postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSLMODE)
MIGRATIONS_DIR ?= ./migrations
GOOSE_VERSION = v3.14.0
BUILD_DIR = ./build
BINARY_NAME = pvz-api

help: ## Вывод справки по командам
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Сборка приложения
	@echo "Сборка приложения..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/api

test: ## Запуск всех тестов
	@echo "Запуск тестов..."
	@go test -v ./...

test-cover: ## Запуск тестов с покрытием
	@echo "Запуск тестов с покрытием..."
	@go test -cover -v ./...

deps: ## Установка зависимостей
	@echo "Установка зависимостей..."
	@go mod tidy
	@go install github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION)

migrate-up: ## Выполнение всех миграций
	@echo "Выполнение миграций..."
	@goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up

migrate-down: ## Откат последней миграции
	@echo "Откат последней миграции..."
	@goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" down

migrate-reset: ## Откат всех миграций
	@echo "Откат всех миграций..."
	@goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" reset

migrate-status: ## Вывод статуса миграций
	@echo "Статус миграций:"
	@goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" status

migrate-create: ## Создание новой миграции (использование: make migrate-create name=имя_миграции)
	@if [ -z "$(name)" ]; then \
		echo "Укажите имя миграции: make migrate-create name=имя_миграции"; \
		exit 1; \
	fi
	@echo "Создание новой миграции: $(name)..."
	@goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" create $(name) sql

run: ## Запуск API сервера
	@echo "Запуск API сервера..."
	@go run ./cmd/api

docker-build: ## Сборка Docker-образа
	@echo "Сборка Docker-образа..."
	@docker build -t $(BINARY_NAME) .

docker-up: ## Запуск Docker-контейнеров
	@echo "Запуск Docker-контейнеров..."
	@docker compose up -d

docker-down: ## Остановка Docker-контейнеров
	@echo "Остановка Docker-контейнеров..."
	@docker compose down

# По умолчанию выводим справку
.DEFAULT_GOAL := help 