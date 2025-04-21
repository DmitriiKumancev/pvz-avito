package reception

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dkumancev/avito-pvz/pkg/domain"
	"github.com/dkumancev/avito-pvz/pkg/infrastructure/postgres/models"
)

// GetLastActiveByPVZID получает последнюю активную приемку для конкретного ПВЗ
func (r *Repository) GetLastActiveByPVZID(ctx context.Context, pvzID string) (*domain.Reception, error) {
	query := `
		SELECT id, date_time, pvz_id, status 
		FROM reception 
		WHERE pvz_id = $1 AND status = $2 
		ORDER BY date_time DESC 
		LIMIT 1
	`

	model := &models.ReceptionModel{}
	err := r.db.GetContext(ctx, model, query, pvzID, domain.ReceptionStatusInProgress)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("активная приемка для ПВЗ %s не найдена", pvzID)
		}
		return nil, fmt.Errorf("ошибка при получении активной приемки: %w", err)
	}

	reception := model.ToEntity()
	products, err := r.getProductsByReceptionID(ctx, reception.ID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении товаров для приемки: %w", err)
	}

	reception.Products = products
	return reception, nil
}
