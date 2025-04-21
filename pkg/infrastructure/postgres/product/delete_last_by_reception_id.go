package product

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// DeleteLastByReceptionID удаляет последний добавленный товар из приемки
func (r *Repository) DeleteLastByReceptionID(ctx context.Context, receptionID string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("ошибка начала транзакции: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var productID string
	query := `
		SELECT ps.product_id
		FROM product_sequence ps
		WHERE ps.reception_id = $1
		ORDER BY ps.id DESC
		LIMIT 1
	`
	err = tx.GetContext(ctx, &productID, query, receptionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("товары для приемки с ID %s не найдены", receptionID)
		}
		return fmt.Errorf("ошибка при получении последнего товара: %w", err)
	}

	// удаляем из очереди
	_, err = tx.ExecContext(ctx, "DELETE FROM product_sequence WHERE product_id = $1", productID)
	if err != nil {
		return fmt.Errorf("ошибка при удалении товара из очереди: %w", err)
	}

	// удаляем сам товар
	_, err = tx.ExecContext(ctx, "DELETE FROM product WHERE id = $1", productID)
	if err != nil {
		return fmt.Errorf("ошибка при удалении товара: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("ошибка фиксации транзакции: %w", err)
	}

	return nil
}
