package pvz

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/dkumancev/avito-pvz/pkg/domain"
	"github.com/dkumancev/avito-pvz/pkg/infrastructure/postgres/models"
)

// Create создает новый ПВЗ в базе данных
func (r *Repository) Create(ctx context.Context, pvz *domain.PVZ) (*domain.PVZ, error) {
	if pvz.ID == "" {
		pvz.ID = uuid.New().String()
	}

	if pvz.RegistrationDate.IsZero() {
		pvz.RegistrationDate = time.Now()
	}

	// Создаём модельку для БД
	model := &models.PVZModel{}
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
