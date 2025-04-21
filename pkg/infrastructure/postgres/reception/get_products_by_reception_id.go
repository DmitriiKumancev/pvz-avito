package reception

import (
	"context"
	"fmt"

	"github.com/dkumancev/avito-pvz/pkg/domain"
	"github.com/dkumancev/avito-pvz/pkg/infrastructure/postgres/models"
)

// getProductsByReceptionID получает список товаров для конкретной приемки
func (r *Repository) getProductsByReceptionID(ctx context.Context, receptionID string) ([]domain.Product, error) {
	query := `
		SELECT p.id, p.date_time, p.type, p.reception_id
		FROM product p
		JOIN product_sequence ps ON p.id = ps.product_id
		WHERE p.reception_id = $1
		ORDER BY ps.id  -- Сортируем по порядку добавления (для LIFO)
	`

	var productModels []models.ProductModel
	err := r.db.SelectContext(ctx, &productModels, query, receptionID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении товаров для приемки: %w", err)
	}

	result := make([]domain.Product, 0, len(productModels))
	for _, model := range productModels {
		product := model.ToEntity()
		result = append(result, *product)
	}

	return result, nil
}
