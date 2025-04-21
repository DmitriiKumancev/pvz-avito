package reception

import (
	"context"
	"fmt"

	"github.com/dkumancev/avito-pvz/pkg/domain"
	"github.com/dkumancev/avito-pvz/pkg/infrastructure/postgres/models"
)

// GetByPVZID получает список всех приемок для конкретного ПВЗ
func (r *Repository) GetByPVZID(ctx context.Context, pvzID string) ([]*domain.Reception, error) {
	query := `
		SELECT id, date_time, pvz_id, status 
		FROM reception 
		WHERE pvz_id = $1 
		ORDER BY date_time DESC
	`

	var receptionModels []models.ReceptionModel
	err := r.db.SelectContext(ctx, &receptionModels, query, pvzID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении приемок для ПВЗ: %w", err)
	}

	result := make([]*domain.Reception, 0, len(receptionModels))
	for _, model := range receptionModels {
		reception := model.ToEntity()

		// все товары для каждой приемки
		products, err := r.getProductsByReceptionID(ctx, reception.ID)
		if err != nil {
			return nil, fmt.Errorf("ошибка при получении товаров для приемки: %w", err)
		}

		reception.Products = products
		result = append(result, reception)
	}

	return result, nil
}
