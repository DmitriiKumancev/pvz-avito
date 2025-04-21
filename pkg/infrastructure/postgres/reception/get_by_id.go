package reception

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dkumancev/avito-pvz/pkg/domain"
	"github.com/dkumancev/avito-pvz/pkg/infrastructure/postgres/models"
)

// GetByID получает приемку по её идентификатору
func (r *Repository) GetByID(ctx context.Context, id string) (*domain.Reception, error) {
	query := `SELECT id, date_time, pvz_id, status FROM reception WHERE id = $1`

	model := &models.ReceptionModel{}
	err := r.db.GetContext(ctx, model, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("приемка с ID %s не найдена", id)
		}
		return nil, fmt.Errorf("ошибка при получении приемки: %w", err)
	}

	// get все товары для этой приемки
	reception := model.ToEntity()
	products, err := r.getProductsByReceptionID(ctx, reception.ID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении товаров для приемки: %w", err)
	}

	reception.Products = products
	return reception, nil
}
