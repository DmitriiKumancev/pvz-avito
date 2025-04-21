package reception

import (
	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB) *Repository {
	return NewRepository(db)
}
