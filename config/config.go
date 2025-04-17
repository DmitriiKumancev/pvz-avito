package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Postgres PostgresConfig
	HTTP     HTTPConfig
	GRPC     GRPCConfig
	Metrics  MetricsConfig
	Auth     AuthConfig
}

// общие настройки приложения
type AppConfig struct {
	Environment   string // development, staging, production
	LogLevel      string // debug, info, warn, error
	MigrationsDir string // путь к директории с миграциями
}

type PostgresConfig struct {
	Host            string
	Port            string
	Username        string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	MaxIdleConns    int
	ConnMaxIdleTime time.Duration
}

func (c PostgresConfig) ConnURI() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.Username, c.Password, c.DBName, c.SSLMode)
}

type HTTPConfig struct {
	Port    string
	Timeout time.Duration
}

type GRPCConfig struct {
	Port string
}

type MetricsConfig struct {
	Port string
}

type AuthConfig struct {
	JWTSecret string
	TokenTTL  time.Duration
}


func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		dir, err := os.Getwd()
		if err != nil {
			log.Printf("Не удалось получить текущую директорию: %v", err)
			return
		}

		for {
			parent := filepath.Dir(dir)
			if parent == dir {
				break 
			}
			dir = parent

			envPath := filepath.Join(dir, ".env")
			if _, err := os.Stat(envPath); err == nil {
				if err := godotenv.Load(envPath); err == nil {
					log.Printf("Загружены переменные окружения из %s", envPath)
					return
				}
			}
		}

		log.Printf("Внимание: .env файл не найден, используются переменные окружения системы")
	} else {
		log.Printf("Загружены переменные окружения из .env файла")
	}
}

func NewConfig() (*Config, error) {
	LoadEnv()

	// Настройки приложения
	environment := getEnv("APP_ENVIRONMENT", "development")
	logLevel := getEnv("APP_LOG_LEVEL", "info")
	migrationsDir := getEnv("APP_MIGRATIONS_DIR", "./migrations")

	// Настройки PostgreSQL
	pgHost := getEnv("POSTGRES_HOST", "localhost")
	pgPort := getEnv("POSTGRES_PORT", "5432")
	pgUser := getEnv("POSTGRES_USER", "postgres")
	pgPassword := getEnv("POSTGRES_PASSWORD", "postgres")
	pgDBName := getEnv("POSTGRES_DB", "pvz")
	pgSSLMode := getEnv("POSTGRES_SSLMODE", "disable")

	pgMaxOpenConns, err := strconv.Atoi(getEnv("POSTGRES_MAX_OPEN_CONNS", "60"))
	if err != nil {
		log.Printf("Неверное значение POSTGRES_MAX_OPEN_CONNS, используется значение по умолчанию: %v", err)
		pgMaxOpenConns = 60
	}

	pgConnMaxLifetime, err := time.ParseDuration(getEnv("POSTGRES_CONN_MAX_LIFETIME", "120s"))
	if err != nil {
		log.Printf("Неверное значение POSTGRES_CONN_MAX_LIFETIME, используется значение по умолчанию: %v", err)
		pgConnMaxLifetime = 120 * time.Second
	}

	pgMaxIdleConns, err := strconv.Atoi(getEnv("POSTGRES_MAX_IDLE_CONNS", "30"))
	if err != nil {
		log.Printf("Неверное значение POSTGRES_MAX_IDLE_CONNS, используется значение по умолчанию: %v", err)
		pgMaxIdleConns = 30
	}

	pgConnMaxIdleTime, err := time.ParseDuration(getEnv("POSTGRES_CONN_MAX_IDLE_TIME", "20s"))
	if err != nil {
		log.Printf("Неверное значение POSTGRES_CONN_MAX_IDLE_TIME, используется значение по умолчанию: %v", err)
		pgConnMaxIdleTime = 20 * time.Second
	}

	// Настройки HTTP сервера
	httpPort := getEnv("HTTP_PORT", "8080")
	httpTimeout, err := time.ParseDuration(getEnv("HTTP_TIMEOUT", "30s"))
	if err != nil {
		log.Printf("Неверное значение HTTP_TIMEOUT, используется значение по умолчанию: %v", err)
		httpTimeout = 30 * time.Second
	}

	// Настройки GRPC сервера
	grpcPort := getEnv("GRPC_PORT", "3000")

	// Настройки метрик
	metricsPort := getEnv("METRICS_PORT", "9000")

	// Настройки авторизации
	jwtSecret := getEnv("JWT_SECRET", "your-secret-key") // TODO: в продакшене secret key должен быть задан
	tokenTTL, err := time.ParseDuration(getEnv("TOKEN_TTL", "24h"))
	if err != nil {
		log.Printf("Неверное значение TOKEN_TTL, используется значение по умолчанию: %v", err)
		tokenTTL = 24 * time.Hour
	}

	return &Config{
		App: AppConfig{
			Environment:   environment,
			LogLevel:      logLevel,
			MigrationsDir: migrationsDir,
		},
		Postgres: PostgresConfig{
			Host:            pgHost,
			Port:            pgPort,
			Username:        pgUser,
			Password:        pgPassword,
			DBName:          pgDBName,
			SSLMode:         pgSSLMode,
			MaxOpenConns:    pgMaxOpenConns,
			ConnMaxLifetime: pgConnMaxLifetime,
			MaxIdleConns:    pgMaxIdleConns,
			ConnMaxIdleTime: pgConnMaxIdleTime,
		},
		HTTP: HTTPConfig{
			Port:    httpPort,
			Timeout: httpTimeout,
		},
		GRPC: GRPCConfig{
			Port: grpcPort,
		},
		Metrics: MetricsConfig{
			Port: metricsPort,
		},
		Auth: AuthConfig{
			JWTSecret: jwtSecret,
			TokenTTL:  tokenTTL,
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
