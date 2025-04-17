package migrations

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/dkumancev/avito-pvz/config"
	"github.com/pressly/goose/v3"
)

type Runner struct {
	db  *sql.DB
	cfg *config.Config
}

func NewRunner(db *sql.DB, cfg *config.Config) *Runner {
	return &Runner{
		db:  db,
		cfg: cfg,
	}
}

// RunMigrationsUp выполняет миграции базы данных вверх до последней версии
func (r *Runner) RunMigrationsUp() error {
	log.Println("Запуск миграций базы данных...")
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("ошибка установки диалекта: %w", err)
	}

	if err := goose.Up(r.db, r.cfg.App.MigrationsDir); err != nil {
		return fmt.Errorf("ошибка выполнения миграций up: %w", err)
	}

	log.Println("Миграции успешно выполнены")
	return nil
}

// RunMigrationsDown откатывает все миграции
func (r *Runner) RunMigrationsDown() error {
	log.Println("Откат всех миграций базы данных...")
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("ошибка установки диалекта: %w", err)
	}

	if err := goose.Reset(r.db, r.cfg.App.MigrationsDir); err != nil {
		return fmt.Errorf("ошибка отката миграций: %w", err)
	}

	log.Println("Откат миграций успешно выполнен")
	return nil
}

// RunMigrationToVersion выполняет миграции до указанной версии
func (r *Runner) RunMigrationToVersion(version int64) error {
	log.Printf("Миграция базы данных до версии %d...\n", version)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("ошибка установки диалекта: %w", err)
	}

	if err := goose.UpTo(r.db, r.cfg.App.MigrationsDir, version); err != nil {
		return fmt.Errorf("ошибка миграции до версии %d: %w", version, err)
	}

	log.Printf("Миграция до версии %d успешно выполнена\n", version)
	return nil
}

// GetCurrentVersion возвращает текущую версию миграции
func (r *Runner) GetCurrentVersion() (int64, error) {
	if err := goose.SetDialect("postgres"); err != nil {
		return 0, fmt.Errorf("ошибка установки диалекта: %w", err)
	}

	// В goose v3 для программного получения текущей версии
	// нужно напрямую запросить её из базы данных
	var version int64
	tableName := goose.TableName()
	query := fmt.Sprintf("SELECT version_id FROM %s ORDER BY id DESC LIMIT 1", tableName)

	err := r.db.QueryRow(query).Scan(&version)
	if err != nil {
		if err == sql.ErrNoRows {
			// База данных существует, но миграции ещё не выполнялись
			return 0, nil
		}
		return 0, fmt.Errorf("ошибка получения версии БД: %w", err)
	}

	return version, nil
}
