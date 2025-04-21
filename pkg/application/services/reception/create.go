package reception

import (
	"context"
	"errors"
	"fmt"

	"github.com/dkumancev/avito-pvz/pkg/domain"
)

func (s *service) CreateReception(ctx context.Context, pvzID string) (*domain.Reception, error) {
	pvz, err := s.pvzRepo.GetByID(ctx, pvzID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения ПВЗ: %w", err)
	}

	// Проверяем, что нет активной приемки
	existingReception, err := s.receptionRepo.GetLastActiveByPVZID(ctx, pvzID)
	if err == nil && existingReception != nil {
		return nil, errors.New("для данного ПВЗ уже существует активная приемка")
	}

	reception := domain.NewReception(pvz.ID)

	savedReception, err := s.receptionRepo.Create(ctx, reception)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания приемки: %w", err)
	}

	return savedReception, nil
}
