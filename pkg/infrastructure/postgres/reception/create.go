package reception

import (
	"context"
	"fmt"
	"time"

	"github.com/dkumancev/avito-pvz/pkg/domain"
	"github.com/dkumancev/avito-pvz/pkg/infrastructure/postgres/models"
)

// Create создает новую приемку товаров в базе данных
func (r *Repository) Create(ctx context.Context, reception *domain.Reception) (*domain.Reception, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка начала транзакции: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	model := &models.ReceptionModel{}
	model.FromEntity(reception)

	if model.DateTime.IsZero() {
		model.DateTime = time.Now()
	}

	query := `
		INSERT INTO reception (date_time, pvz_id, status) 
		VALUES (:date_time, :pvz_id, :status) 
		RETURNING id, date_time, pvz_id, status
	`

	stmt, err := tx.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка подготовки запроса: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowxContext(ctx, model).StructScan(model)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании приемки: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("ошибка фиксации транзакции: %w", err)
	}

	result := model.ToEntity()
	return result, nil
}
