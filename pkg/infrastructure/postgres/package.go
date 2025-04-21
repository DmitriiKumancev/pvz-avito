package postgres

import (
	"github.com/jmoiron/sqlx"

	"github.com/dkumancev/avito-pvz/pkg/application/repositories"
	"github.com/dkumancev/avito-pvz/pkg/infrastructure/postgres/product"
	"github.com/dkumancev/avito-pvz/pkg/infrastructure/postgres/pvz"
	"github.com/dkumancev/avito-pvz/pkg/infrastructure/postgres/reception"
	"github.com/dkumancev/avito-pvz/pkg/infrastructure/postgres/user"
)

type Repositories struct {
	User      repositories.UserRepository
	PVZ       repositories.PVZRepository
	Reception *reception.Repository
	Product   *product.Repository
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		User:      user.New(db),
		PVZ:       pvz.New(db),
		Reception: reception.New(db),
		Product:   product.New(db),
	}
}
