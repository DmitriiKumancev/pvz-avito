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
	"github.com/google/uuid"
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
	if pvz.ID == "" {
		pvz.ID = uuid.New().String()
	}

	if pvz.RegistrationDate.IsZero() {
		pvz.RegistrationDate = time.Now()
	}

	// Создаём модельку для БД
	model := &PVZModel{}
	model.FromEntity(pvz)

	query := `
		INSERT INTO pvz (id, registration_date, city)
		VALUES (:id, :registration_date, :city)
		RETURNING id, registration_date, city
	`

	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка подготовки запроса: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowxContext(ctx, model).StructScan(model)
	if err != nil {
		if err.Error() == "ERROR: новое значение для \"city\" нарушает ограничение-проверку \"pvz_valid_cities\" (SQLSTATE 23514)" {
			return nil, errors.New("город не поддерживается: разрешены только Москва, Санкт-Петербург и Казань")
		}
		return nil, fmt.Errorf("ошибка создания ПВЗ: %w", err)
	}

	result := model.ToEntity()
	return result, nil
}

func (r *PVZRepository) GetByID(ctx context.Context, id string) (*domain.PVZ, error) {
	query := `SELECT id, registration_date, city FROM pvz WHERE id = $1`

	model := &PVZModel{}
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

func (r *PVZRepository) List(ctx context.Context, filter repositories.PVZFilter) ([]*domain.PVZ, error) {
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

	var models []PVZModel
	err := r.db.SelectContext(ctx, &models, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении списка ПВЗ: %w", err)
	}

	result := make([]*domain.PVZ, 0, len(models))
	for _, model := range models {
		model := model // локальная копия перемнной для безопаснрго использования в замыкании
		result = append(result, model.ToEntity())
	}

	return result, nil
}
