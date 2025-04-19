package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/dkumancev/avito-pvz/pkg/domain"
)

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) Create(ctx context.Context, product *domain.Product, receptionID string) (*domain.Product, error) {
	// транзакция для атомарной операции
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
	model := &ProductModel{}
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

func (r *ProductRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	query := `SELECT id, date_time, type, reception_id FROM product WHERE id = $1`

	model := &ProductModel{}
	err := r.db.GetContext(ctx, model, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("товар с ID %s не найден", id)
		}
		return nil, fmt.Errorf("ошибка при получении товара: %w", err)
	}

	// преобразуем модельку обратно в доменную сущность
	result := model.ToEntity()
	return result, nil
}

func (r *ProductRepository) DeleteByID(ctx context.Context, id string) error {
	//  транзакция для атомарной операции
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

func (r *ProductRepository) GetByReceptionID(ctx context.Context, receptionID string) ([]*domain.Product, error) {
	query := `
		SELECT p.id, p.date_time, p.type, p.reception_id
		FROM product p
		JOIN product_sequence ps ON p.id = ps.product_id
		WHERE p.reception_id = $1
		ORDER BY ps.id
	`

	var models []ProductModel
	err := r.db.SelectContext(ctx, &models, query, receptionID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении товаров для приемки: %w", err)
	}

	result := make([]*domain.Product, 0, len(models))
	for _, model := range models {
		product := model.ToEntity()
		result = append(result, product)
	}

	return result, nil
}

func (r *ProductRepository) DeleteLastByReceptionID(ctx context.Context, receptionID string) error {
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

	// удаляем товар из очереди товаров
	_, err = tx.ExecContext(ctx, "DELETE FROM product_sequence WHERE product_id = $1", productID)
	if err != nil {
		return fmt.Errorf("ошибка при удалении товара из очереди: %w", err)
	}

	// Удаляем сам товар
	_, err = tx.ExecContext(ctx, "DELETE FROM product WHERE id = $1", productID)
	if err != nil {
		return fmt.Errorf("ошибка при удалении товара: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("ошибка фиксации транзакции: %w", err)
	}

	return nil
}

func (r *ProductRepository) ListByReceptionID(ctx context.Context, receptionID string) ([]domain.Product, error) {
	query := `
		SELECT p.id, p.date_time, p.type, p.reception_id
		FROM product p
		JOIN product_sequence ps ON p.id = ps.product_id
		WHERE p.reception_id = $1
		ORDER BY ps.id
	`

	var models []ProductModel
	err := r.db.SelectContext(ctx, &models, query, receptionID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении товаров для приемки: %w", err)
	}

	result := make([]domain.Product, 0, len(models))
	for _, model := range models {
		product := model.ToEntity()
		result = append(result, *product)
	}

	return result, nil
}

func (r *ProductRepository) GetLastAddedProduct(ctx context.Context, receptionID string) (*domain.Product, error) {
	query := `
		SELECT p.id, p.date_time, p.type, p.reception_id
		FROM product p
		JOIN product_sequence ps ON p.id = ps.product_id
		WHERE p.reception_id = $1
		ORDER BY ps.id DESC
		LIMIT 1
	`

	model := &ProductModel{}
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
