.PHONY: build test test-unit test-integration test-services test-domain test-cover test-cover-domain test-cover-services test-cover-func test-race migrate-up migrate-down migrate-status migrate-create migrate-reset help run-grpc

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
COVERAGE_DIR = ./coverage
COVERAGE_FILE = $(COVERAGE_DIR)/coverage.out
COVERAGE_FILE_DOMAIN = $(COVERAGE_DIR)/coverage_domain.out
COVERAGE_FILE_SERVICES = $(COVERAGE_DIR)/coverage_services.out


help: ## Вывод справки по командам
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Сборка приложения
	@echo "Сборка приложения..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/api

test: ## Запуск всех тестов
	@echo "Запуск тестов..."
	@go test -v ./...

test-unit: ## Запуск только юнит-тестов (без интеграционных)
	@echo "Запуск юнит-тестов..."
	@go test -v $(shell go list ./... | grep -v /tests$$) -short

test-integration: ## Запуск только интеграционных тестов
	@echo "Запуск интеграционных тестов..."
	@go test -v ./pkg/tests

test-services: ## Запуск тестов сервисов
	@echo "Запуск тестов сервисов..."
	@go test -v ./pkg/application/services/...

test-domain: ## Запуск тестов доменных моделей
	@echo "Запуск тестов доменных моделей..."
	@go test -v ./pkg/domain/...

test-cover: ## Запуск всех тестов с покрытием
	@echo "Запуск тестов с покрытием..."
	@mkdir -p $(COVERAGE_DIR)
	@go test -coverprofile=$(COVERAGE_FILE) -coverpkg=./pkg/... -covermode=atomic -v ./...

test-cover-domain: ## Запуск тестов доменных моделей с покрытием
	@echo "Запуск тестов домена с покрытием..."
	@mkdir -p $(COVERAGE_DIR)
	@go test -coverprofile=$(COVERAGE_FILE_DOMAIN) -coverpkg=./pkg/domain/... -covermode=atomic -v ./pkg/domain/...
	@go tool cover -func=$(COVERAGE_FILE_DOMAIN)

test-cover-services: ## Запуск тестов сервисов с покрытием
	@echo "Запуск тестов сервисов с покрытием..."
	@mkdir -p $(COVERAGE_DIR)
	@go test -coverprofile=$(COVERAGE_FILE_SERVICES) -coverpkg=./pkg/application/services/... -covermode=atomic -v ./pkg/application/services/...
	@go tool cover -func=$(COVERAGE_FILE_SERVICES)

test-race: ## Запуск тестов с проверкой на race condition
	@echo "Запуск тестов с проверкой на race condition..."
	@go test -race -v ./...

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

run-grpc: ## Запуск gRPC сервера
	@echo "Запуск gRPC сервера..."
	@go run ./cmd/grpc

docker-build: ## Сборка Docker-образа
	@echo "Сборка Docker-образа..."
	@docker build -t $(BINARY_NAME) .

docker-up: ## Запуск Docker-контейнеров
	@echo "Запуск Docker-контейнеров..."
	@docker compose up -d

docker-down: ## Остановка Docker-контейнеров
	@echo "Остановка Docker-контейнеров..."
	@docker compose down

clean: ## Очистка артефактов сборки и тестирования
	@echo "Очистка артефактов..."
	@rm -rf $(BUILD_DIR) $(COVERAGE_DIR)

# По умолчанию выводим справку
.DEFAULT_GOAL := help 