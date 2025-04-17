package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/dkumancev/avito-pvz/pkg/application/repositories"
	"github.com/dkumancev/avito-pvz/pkg/domain"
)

type PVZRepository struct {
	db *sqlx.DB
}

func NewPVZRepository(db *sqlx.DB) *PVZRepository {
	return &PVZRepository{
		db: db,
	}
}

func (r *PVZRepository) Create(ctx context.Context, pvz *domain.PVZ) (*domain.PVZ, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка начала транзакции: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	model := &PVZModel{}
	model.FromEntity(pvz)

	query := `
		INSERT INTO pvz (city, registration_date) 
		VALUES (:city, :registration_date) 
		RETURNING id, city, registration_date
	`

	// Если дата регистрации не указана, тогда текущее время
	if model.RegistrationDate.IsZero() {
		model.RegistrationDate = time.Now()
	}

	// NamedExec для безопасной подстановки параметров и получения результата
	stmt, err := tx.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка подготовки запроса: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowxContext(ctx, model).StructScan(model)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании ПВЗ: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("ошибка фиксации транзакции: %w", err)
	}

	result := model.ToEntity()
	return result, nil
}

func (r *PVZRepository) GetByID(ctx context.Context, id string) (*domain.PVZ, error) {
	query := `SELECT id, city, registration_date FROM pvz WHERE id = $1`

	model := &PVZModel{}
	err := r.db.GetContext(ctx, model, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("ПВЗ с ID %s не найден", id)
		}
		return nil, fmt.Errorf("ошибка при получении ПВЗ: %w", err)
	}

	return model.ToEntity(), nil
}

func (r *PVZRepository) List(ctx context.Context, filter repositories.PVZFilter) ([]*domain.PVZ, error) {
	baseQuery := `
		SELECT DISTINCT p.id, p.city, p.registration_date
		FROM pvz p
	`

	// JOIN с таблицей reception, если указаны даты фильтрации
	var whereClause string
	var args []interface{}
	var joinClause string

	if filter.StartDate != "" || filter.EndDate != "" {
		joinClause = " JOIN reception r ON p.id = r.pvz_id"
		whereClause = " WHERE 1=1"

		if filter.StartDate != "" {
			whereClause += " AND r.date_time >= $1"
			args = append(args, filter.StartDate)
		}

		if filter.EndDate != "" {
			paramIdx := len(args) + 1
			whereClause += fmt.Sprintf(" AND r.date_time <= $%d", paramIdx)
			args = append(args, filter.EndDate)
		}
	}

	//  пагинация
	offset := (filter.Page - 1) * filter.Limit
	limitClause := fmt.Sprintf(" ORDER BY p.registration_date DESC LIMIT $%d OFFSET $%d",
		len(args)+1, len(args)+2)
	args = append(args, filter.Limit, offset)

	// формирование итогового запроса
	query := baseQuery + joinClause + whereClause + limitClause

	var models []PVZModel
	err := r.db.SelectContext(ctx, &models, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении списка ПВЗ: %w", err)
	}

	//  преобразование моделей в доменные сущности
	result := make([]*domain.PVZ, 0, len(models))
	for _, model := range models {
		model := model // локальная копия перемнной для безопаснрго использования в замыкании
		result = append(result, model.ToEntity())
	}

	return result, nil
}
