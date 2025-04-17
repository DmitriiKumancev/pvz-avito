# ПВЗ Авито

Сервис для управления пунктами выдачи заказов и приемки товаров.

## Структура проекта

```
.
├── cmd/                  # Исполняемые файлы
│   └── api/              # API-сервер
├── migrations/           # SQL-миграции базы данных
├── pkg/                  # Исходный код пакетов
│   ├── application/      # Слой бизнес-логики
│   │   ├── repositories/ # Интерфейсы репозиториев
│   │   └── services/     # Сервисы бизнес-логики
│   ├── domain/           # Доменные модели
│   └── infrastructure/   # Внешние адаптеры
│       ├── migrations/   # Утилиты для программного управления миграциями
│       └── postgres/     # Реализация репозиториев в PostgreSQL
├── docker-compose.yml    # Конфигурация Docker Compose
├── Dockerfile            # Файл сборки Docker-образа
├── go.mod                # Зависимости Go
├── go.sum                # Хеши зависимостей Go
└── Makefile              # Команды для сборки и управления
```

## Предварительные требования

- Go 1.22+
- PostgreSQL 14+
- Docker & Docker Compose (опционально)
- Make

## Установка зависимостей

```bash
make deps
```

Эта команда установит необходимые зависимости, включая утилиту goose для
миграций.

## Миграции базы данных

Проект использует [goose](https://github.com/pressly/goose/v3) для управления
миграциями базы данных.

### Запуск миграций с помощью Makefile

В проекте есть несколько команд для управления миграциями:

```bash
# Применить все миграции
make migrate-up

# Откатить последнюю миграцию
make migrate-down

# Откатить все миграции
make migrate-reset

# Проверить статус миграций
make migrate-status

# Создать новую миграцию
make migrate-create name=имя_миграции
```

По умолчанию, команды используют подключение к базе данных из переменной
`DB_URL` в Makefile. Вы можете переопределить настройки подключения, запустив
команду с параметром:

```bash
make migrate-up DB_URL=postgres://user:password@host:port/database
```

### Программное управление миграциями

В `pkg/infrastructure/migrations/migrations.go` реализованы функции для
программного управления миграциями, которые можно использовать при запуске
приложения:

```go
// Пример применения миграций программно
db, err := sql.Open("postgres", "postgres://user:password@host:port/database")
if err != nil {
    log.Fatal(err)
}
defer db.Close()

// Применить все миграции
if err := migrations.RunMigrationsUp(db, migrationsDir); err != nil {
    log.Fatal(err)
}

// Получить текущую версию миграций
version, err := migrations.GetCurrentVersion(db, migrationsDir)
if err != nil {
    log.Fatal(err)
}
log.Printf("Текущая версия базы данных: %d", version)
```

## Запуск в Docker

```bash
# Сборка образа
make docker-build

# Запуск контейнеров
make docker-up

# Остановка контейнеров
make docker-down
```

## Разработка

```bash
# Сборка приложения
make build

# Запуск API сервера
make run

# Запуск тестов
make test

# Запуск тестов с покрытием
make test-cover
```

## Справка по командам

Для просмотра всех доступных команд выполните:

```bash
make help
```
