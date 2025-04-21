package reception

import (
	"context"
	"fmt"

	"github.com/dkumancev/avito-pvz/pkg/domain"
)

func (s *service) GetReceptionsByPVZID(ctx context.Context, pvzID string) ([]*domain.Reception, error) {
	_, err := s.pvzRepo.GetByID(ctx, pvzID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения ПВЗ: %w", err)
	}

	receptions, err := s.receptionRepo.GetByPVZID(ctx, pvzID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения приемок: %w", err)
	}

	return receptions, nil
}

func (s *service) GetProductsByReceptionID(ctx context.Context, receptionID string) ([]*domain.Product, error) {
	_, err := s.receptionRepo.GetByID(ctx, receptionID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения приемки: %w", err)
	}

	products, err := s.productRepo.GetByReceptionID(ctx, receptionID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения товаров: %w", err)
	}

	return products, nil
}
