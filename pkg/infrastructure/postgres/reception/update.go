package reception

import (
	"context"
	"fmt"

	"github.com/dkumancev/avito-pvz/pkg/domain"
	"github.com/dkumancev/avito-pvz/pkg/infrastructure/postgres/models"
)

// Update обновляет информацию о приемке в базе данных
func (r *Repository) Update(ctx context.Context, reception *domain.Reception) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("ошибка начала транзакции: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	model := &models.ReceptionModel{}
	model.FromEntity(reception)

	query := `
		UPDATE reception 
		SET status = :status, date_time = :date_time 
		WHERE id = :id
	`

	stmt, err := tx.PrepareNamedContext(ctx, query)
	if err != nil {
		return fmt.Errorf("ошибка подготовки запроса: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, model)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении приемки: %w", err)
	}

	// check на обновление записи
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества обновленных записей: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("приемка с ID %s не найдена", reception.ID)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("ошибка фиксации транзакции: %w", err)
	}

	return nil
}
