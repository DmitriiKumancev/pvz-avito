package repositories

import (
	"context"
	"time"

	"github.com/dkumancev/avito-pvz/pkg/domain"
)

type PVZRepository interface {
	Create(ctx context.Context, pvz *domain.PVZ) (*domain.PVZ, error)

	GetByID(ctx context.Context, id string) (*domain.PVZ, error)

	List(ctx context.Context, filter PVZFilter) ([]*domain.PVZ, error)
}

// параметры фильтрации для списка ПВЗ
type PVZFilter struct {
	ReceptionStartDate *time.Time
	ReceptionEndDate   *time.Time
	Page               int
	Limit              int
}
