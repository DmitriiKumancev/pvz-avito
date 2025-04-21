package pvz

import (
	"context"
	"fmt"

	"github.com/dkumancev/avito-pvz/pkg/application/repositories"
	"github.com/dkumancev/avito-pvz/pkg/domain"
	"github.com/dkumancev/avito-pvz/pkg/infrastructure/postgres/models"
)

// List возвращает список ПВЗ с возможностью фильтрации
func (r *Repository) List(ctx context.Context, filter repositories.PVZFilter) ([]*domain.PVZ, error) {
	baseQuery := `
		SELECT p.id, p.registration_date, p.city
		FROM pvz p
	`

	// JOIN с таблицей reception, если указаны даты фильтрации
	var whereClause string
	var args []interface{}
	var joinClause string

	if filter.ReceptionStartDate != nil || filter.ReceptionEndDate != nil {
		joinClause = " JOIN reception r ON p.id = r.pvz_id"
		whereClause = " WHERE 1=1"

		if filter.ReceptionStartDate != nil {
			whereClause += " AND r.date_time >= $1"
			args = append(args, filter.ReceptionStartDate)
		}

		if filter.ReceptionEndDate != nil {
			paramIdx := len(args) + 1
			whereClause += fmt.Sprintf(" AND r.date_time <= $%d", paramIdx)
			args = append(args, filter.ReceptionEndDate)
		}
	}

	//  пагинация
	offset := (filter.Page - 1) * filter.Limit
	limitClause := fmt.Sprintf(" ORDER BY p.registration_date DESC LIMIT $%d OFFSET $%d",
		len(args)+1, len(args)+2)
	args = append(args, filter.Limit, offset)

	query := baseQuery + joinClause + whereClause + limitClause

	var pvzModels []models.PVZModel
	err := r.db.SelectContext(ctx, &pvzModels, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении списка ПВЗ: %w", err)
	}

	result := make([]*domain.PVZ, 0, len(pvzModels))
	for _, model := range pvzModels {
		model := model // локальная копия перемнной для безопаснрго использования в замыкании
		result = append(result, model.ToEntity())
	}

	return result, nil
}
