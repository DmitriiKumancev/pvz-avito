package pvz

import (
	"context"
	"fmt"

	"github.com/dkumancev/avito-pvz/pkg/domain"
)

func (s *service) CreatePVZ(ctx context.Context, city string) (*domain.PVZ, error) {
	pvz, err := domain.NewPVZ(city)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания ПВЗ: %w", err)
	}

	savedPVZ, err := s.pvzRepo.Create(ctx, pvz)
	if err != nil {
		return nil, fmt.Errorf("ошибка сохранения ПВЗ: %w", err)
	}

	return savedPVZ, nil
}
