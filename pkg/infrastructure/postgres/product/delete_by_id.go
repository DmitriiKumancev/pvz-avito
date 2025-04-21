package product

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// DeleteByID удаляет товар по его ID
func (r *Repository) DeleteByID(ctx context.Context, id string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("ошибка начала транзакции: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// инфо о товаре, чтобы узнать, в какой приемке он находится
	var receptionID string
	err = tx.GetContext(ctx, &receptionID, "SELECT reception_id FROM product WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("товар с ID %s не найден", id)
		}
		return fmt.Errorf("ошибка при получении информации о товаре: %w", err)
	}

	// 1. удаляем из очереди товаров
	_, err = tx.ExecContext(ctx, "DELETE FROM product_sequence WHERE product_id = $1", id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении товара из очереди: %w", err)
	}

	// 2. удаляем сам товар
	result, err := tx.ExecContext(ctx, "DELETE FROM product WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении товара: %w", err)
	}

	// check на удаление записи
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества удаленных записей: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("товар с ID %s не найден", id)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("ошибка фиксации транзакции: %w", err)
	}

	return nil
}
