package product

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dkumancev/avito-pvz/pkg/domain"
	"github.com/dkumancev/avito-pvz/pkg/infrastructure/postgres/models"
)

// GetByID получает товар по его ID
func (r *Repository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	query := `SELECT id, date_time, type, reception_id FROM product WHERE id = $1`

	model := &models.ProductModel{}
	err := r.db.GetContext(ctx, model, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("товар с ID %s не найден", id)
		}
		return nil, fmt.Errorf("ошибка при получении товара: %w", err)
	}

	// преобразуем модельку обратно в доменную сущность
	result := model.ToEntity()
	return result, nil
}
