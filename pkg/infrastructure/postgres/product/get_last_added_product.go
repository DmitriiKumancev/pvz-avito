package product

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dkumancev/avito-pvz/pkg/domain"
	"github.com/dkumancev/avito-pvz/pkg/infrastructure/postgres/models"
)

// GetLastAddedProduct получает последний добавленный товар в приемку
func (r *Repository) GetLastAddedProduct(ctx context.Context, receptionID string) (*domain.Product, error) {
	query := `
        SELECT p.id, p.date_time, p.type, p.reception_id
        FROM product p
        JOIN product_sequence ps ON p.id = ps.product_id
        WHERE p.reception_id = $1
        ORDER BY ps.id DESC
        LIMIT 1
    `

	model := &models.ProductModel{}
	err := r.db.GetContext(ctx, model, query, receptionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("товары для приемки с ID %s не найдены", receptionID)
		}
		return nil, fmt.Errorf("ошибка при получении последнего товара: %w", err)
	}

	result := model.ToEntity()
	return result, nil
}
