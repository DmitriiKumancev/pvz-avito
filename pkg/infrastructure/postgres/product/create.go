package product

import (
	"context"
	"fmt"
	"time"

	"github.com/dkumancev/avito-pvz/pkg/domain"
	"github.com/dkumancev/avito-pvz/pkg/infrastructure/postgres/models"
)

// Create создает новый товар и добавляет его в очередь приемки
func (r *Repository) Create(ctx context.Context, product *domain.Product, receptionID string) (*domain.Product, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка начала транзакции: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// create model for db
	model := &models.ProductModel{}
	model.FromEntity(product)
	model.ReceptionID = receptionID

	// Если дата не указана, тогда юзаем текущее время
	if model.DateTime.IsZero() {
		model.DateTime = time.Now()
	}

	insertProductQuery := `
		INSERT INTO product (date_time, type, reception_id) 
		VALUES (:date_time, :type, :reception_id) 
		RETURNING id, date_time, type, reception_id
	`

	stmt, err := tx.PrepareNamedContext(ctx, insertProductQuery)
	if err != nil {
		return nil, fmt.Errorf("ошибка подготовки запроса: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowxContext(ctx, model).StructScan(model)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании товара: %w", err)
	}

	// Добавление товара в очередь товаров для приемки
	insertSequenceQuery := `
		INSERT INTO product_sequence (product_id, reception_id) 
		VALUES ($1, $2)
	`
	_, err = tx.ExecContext(ctx, insertSequenceQuery, model.ID, receptionID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при добавлении товара в очередь: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("ошибка фиксации транзакции: %w", err)
	}

	// преобразуем модельку обратно в доменную сущность
	result := model.ToEntity()
	return result, nil
}
