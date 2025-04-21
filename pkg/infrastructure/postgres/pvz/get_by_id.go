package pvz

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dkumancev/avito-pvz/pkg/domain"
	"github.com/dkumancev/avito-pvz/pkg/infrastructure/postgres/models"
)

// GetByID получает ПВЗ по его идентификатору
func (r *Repository) GetByID(ctx context.Context, id string) (*domain.PVZ, error) {
	query := `SELECT id, registration_date, city FROM pvz WHERE id = $1`

	model := &models.PVZModel{}
	err := r.db.GetContext(ctx, model, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("ПВЗ с ID %s не найден", id)
		}
		return nil, fmt.Errorf("ошибка получения ПВЗ: %w", err)
	}

	result := model.ToEntity()
	return result, nil
}
