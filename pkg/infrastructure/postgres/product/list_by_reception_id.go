package product

import (
	"context"
	"fmt"

	"github.com/dkumancev/avito-pvz/pkg/domain"
	"github.com/dkumancev/avito-pvz/pkg/infrastructure/postgres/models"
)

// ListByReceptionID получает список всех товаров для приемки
// отсортированный по времени создания
func (r *Repository) ListByReceptionID(ctx context.Context, receptionID string) ([]domain.Product, error) {
	query := `
        SELECT p.id, p.date_time, p.type, p.reception_id
        FROM product p
        WHERE p.reception_id = $1
        ORDER BY p.date_time
    `

	var productModels []models.ProductModel
	err := r.db.SelectContext(ctx, &productModels, query, receptionID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении списка товаров: %w", err)
	}

	result := make([]domain.Product, 0, len(productModels))
	for _, model := range productModels {
		result = append(result, *model.ToEntity())
	}

	return result, nil
}
