package user

import (
	"github.com/dkumancev/avito-pvz/pkg/application/repositories"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) repositories.UserRepository {
	return &Repository{
		db: db,
	}
}
