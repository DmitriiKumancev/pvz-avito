package reception

import (
	"context"
	"fmt"

	"github.com/dkumancev/avito-pvz/pkg/domain"
)

func (s *service) CloseReception(ctx context.Context, pvzID string) (*domain.Reception, error) {
	_, err := s.pvzRepo.GetByID(ctx, pvzID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения ПВЗ: %w", err)
	}

	reception, err := s.getActiveReception(ctx, pvzID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить активную приемку: %w", err)
	}

	err = reception.Close()
	if err != nil {
		return nil, fmt.Errorf("ошибка закрытия приемки: %w", err)
	}

	err = s.receptionRepo.Update(ctx, reception)
	if err != nil {
		return nil, fmt.Errorf("ошибка обновления приемки: %w", err)
	}

	return reception, nil
}
