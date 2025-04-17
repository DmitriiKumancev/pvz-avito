package repositories

import (
	"context"

	"github.com/dkumancev/avito-pvz/pkg/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	Exists(ctx context.Context, email string) (bool, error)
}
