package pvz

import (
	"context"

	"github.com/dkumancev/avito-pvz/pkg/application/repositories"
	"github.com/dkumancev/avito-pvz/pkg/domain"
)

type Service interface {
	// Создание нового ПВЗ
	CreatePVZ(ctx context.Context, city string) (*domain.PVZ, error)

	// Получение ПВЗ по ID
	GetPVZByID(ctx context.Context, id string) (*domain.PVZ, error)

	// Получение списка ПВЗ с фильтрацией
	ListPVZs(ctx context.Context, filter repositories.PVZFilter) ([]*domain.PVZ, error)

	// Алиас для ListPVZs
	ListPVZ(ctx context.Context, filter repositories.PVZFilter) ([]*domain.PVZ, error)
}

type service struct {
	pvzRepo repositories.PVZRepository
}

func New(pvzRepo repositories.PVZRepository) Service {
	return &service{
		pvzRepo: pvzRepo,
	}
}
