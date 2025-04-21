package pvz

import (
	"context"
	"fmt"

	"github.com/dkumancev/avito-pvz/pkg/application/repositories"
	"github.com/dkumancev/avito-pvz/pkg/domain"
)

func (s *service) GetPVZByID(ctx context.Context, id string) (*domain.PVZ, error) {
	pvz, err := s.pvzRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения ПВЗ: %w", err)
	}
	return pvz, nil
}

func (s *service) ListPVZs(ctx context.Context, filter repositories.PVZFilter) ([]*domain.PVZ, error) {
	// дефолтные значения для пагинации, если не указаны
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	}

	pvzs, err := s.pvzRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка ПВЗ: %w", err)
	}

	return pvzs, nil
}

// ListPVZ алиас для ListPVZs
func (s *service) ListPVZ(ctx context.Context, filter repositories.PVZFilter) ([]*domain.PVZ, error) {
	return s.ListPVZs(ctx, filter)
}
