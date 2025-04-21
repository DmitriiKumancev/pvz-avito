package db

import (
	"fmt"

	"github.com/dkumancev/avito-pvz/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func New(cfg config.PostgresConfig) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", cfg.ConnURI())
	if err != nil {
		return nil, fmt.Errorf("error connecting to PostgreSQL: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error checking connection to PostgreSQL: %w", err)
	}

	return db, nil
}
