package pvz

import (
	"github.com/jmoiron/sqlx"

	"github.com/dkumancev/avito-pvz/pkg/application/repositories"
)

func New(db *sqlx.DB) repositories.PVZRepository {
	return NewRepository(db)
}
