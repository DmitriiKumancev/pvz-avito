package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

type Config struct {
	Level     string
	AddSource bool
	Output    io.Writer
	Format    string // json или text
}

var defaultConfig = Config{
	Level:     LevelInfo,
	AddSource: false,
	Output:    os.Stdout,
	Format:    "text",
}

func NewLogger(cfg Config) *slog.Logger {
	if cfg.Level == "" {
		cfg.Level = defaultConfig.Level
	}

	if cfg.Output == nil {
		cfg.Output = defaultConfig.Output
	}

	if cfg.Format == "" {
		cfg.Format = defaultConfig.Format
	}

	var level slog.Level
	switch strings.ToLower(cfg.Level) {
	case LevelDebug:
		level = slog.LevelDebug
	case LevelWarn:
		level = slog.LevelWarn
	case LevelError:
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		AddSource: cfg.AddSource,
		Level:     level,
	}

	var handler slog.Handler
	if strings.ToLower(cfg.Format) == "json" {
		handler = slog.NewJSONHandler(cfg.Output, opts)
	} else {
		handler = slog.NewTextHandler(cfg.Output, opts)
	}

	return slog.New(handler)
}

// SanitizeCredentials удаляет чувствительные данные из строки или заменяет их на маску
func SanitizeCredentials(value string) string {
	if len(value) <= 0 {
		return value
	}
	return "[REDACTED]"
}

// SanitizeError маскирует чувствительные данные в сообщении об ошибке
func SanitizeError(err error) string {
	if err == nil {
		return ""
	}

	errMsg := err.Error()

	// маскируем данные, которые могут содержать пароли и секреты
	lowered := strings.ToLower(errMsg)

	if strings.Contains(lowered, "password") ||
		strings.Contains(lowered, "secret") ||
		strings.Contains(lowered, "token") ||
		strings.Contains(lowered, "jwt") ||
		strings.Contains(lowered, "auth") {
		return "Ошибка содержит чувствительные данные"
	}

	return errMsg
}

func NewLoggerFromEnvironment(logLevel string) *slog.Logger {
	return NewLogger(Config{
		Level:     logLevel,
		AddSource: logLevel == LevelDebug, // В режиме отладки добавляем информацию об источнике
		Format:    "text",                 // По умолчанию текстовый формат
	})
}
