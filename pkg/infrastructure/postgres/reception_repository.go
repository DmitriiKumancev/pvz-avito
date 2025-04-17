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

type ReceptionRepository struct {
	db *sqlx.DB
}

func NewReceptionRepository(db *sqlx.DB) *ReceptionRepository {
	return &ReceptionRepository{
		db: db,
	}
}

func (r *ReceptionRepository) Create(ctx context.Context, reception *domain.Reception) (*domain.Reception, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка начала транзакции: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	model := &ReceptionModel{}
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

func (r *ReceptionRepository) GetByID(ctx context.Context, id string) (*domain.Reception, error) {
	query := `SELECT id, date_time, pvz_id, status FROM reception WHERE id = $1`

	model := &ReceptionModel{}
	err := r.db.GetContext(ctx, model, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("приемка с ID %s не найдена", id)
		}
		return nil, fmt.Errorf("ошибка при получении приемки: %w", err)
	}

	// Получаем все товары для этой приемки
	reception := model.ToEntity()
	products, err := r.getProductsByReceptionID(ctx, reception.ID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении товаров для приемки: %w", err)
	}

	reception.Products = products
	return reception, nil
}

func (r *ReceptionRepository) Update(ctx context.Context, reception *domain.Reception) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("ошибка начала транзакции: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	model := &ReceptionModel{}
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

func (r *ReceptionRepository) GetLastActiveByPVZID(ctx context.Context, pvzID string) (*domain.Reception, error) {
	query := `
		SELECT id, date_time, pvz_id, status 
		FROM reception 
		WHERE pvz_id = $1 AND status = $2 
		ORDER BY date_time DESC 
		LIMIT 1
	`

	model := &ReceptionModel{}
	err := r.db.GetContext(ctx, model, query, pvzID, domain.ReceptionStatusInProgress)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("активная приемка для ПВЗ %s не найдена", pvzID)
		}
		return nil, fmt.Errorf("ошибка при получении активной приемки: %w", err)
	}

	reception := model.ToEntity()
	products, err := r.getProductsByReceptionID(ctx, reception.ID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении товаров для приемки: %w", err)
	}

	reception.Products = products
	return reception, nil
}

func (r *ReceptionRepository) GetByPVZID(ctx context.Context, pvzID string) ([]*domain.Reception, error) {
	query := `
		SELECT id, date_time, pvz_id, status 
		FROM reception 
		WHERE pvz_id = $1 
		ORDER BY date_time DESC
	`

	var models []ReceptionModel
	err := r.db.SelectContext(ctx, &models, query, pvzID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении приемок для ПВЗ: %w", err)
	}

	result := make([]*domain.Reception, 0, len(models))
	for _, model := range models {
		reception := model.ToEntity()

		// все товары для каждой приемки
		products, err := r.getProductsByReceptionID(ctx, reception.ID)
		if err != nil {
			return nil, fmt.Errorf("ошибка при получении товаров для приемки: %w", err)
		}

		reception.Products = products
		result = append(result, reception)
	}

	return result, nil
}

func (r *ReceptionRepository) getProductsByReceptionID(ctx context.Context, receptionID string) ([]domain.Product, error) {
	query := `
		SELECT p.id, p.date_time, p.type, p.reception_id
		FROM product p
		JOIN product_sequence ps ON p.id = ps.product_id
		WHERE p.reception_id = $1
		ORDER BY ps.id  -- Сортируем по порядку добавления (для LIFO)
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
