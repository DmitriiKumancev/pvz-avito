package repositories

import (
	"context"

	"github.com/dkumancev/avito-pvz/pkg/domain"
)

type ReceptionRepository interface {
	Create(ctx context.Context, reception *domain.Reception) (*domain.Reception, error)

	GetByID(ctx context.Context, id string) (*domain.Reception, error)

	Update(ctx context.Context, reception *domain.Reception) error

	GetLastActiveByPVZID(ctx context.Context, pvzID string) (*domain.Reception, error)

	GetByPVZID(ctx context.Context, pvzID string) ([]*domain.Reception, error)
}
