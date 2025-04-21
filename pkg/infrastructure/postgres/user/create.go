package user

import (
	"context"

	"github.com/dkumancev/avito-pvz/pkg/domain"
	"github.com/dkumancev/avito-pvz/pkg/infrastructure/postgres/models"
	"github.com/google/uuid"
)

// Create создает нового пользователя в базе данных
func (r *Repository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	query := `
		INSERT INTO users (id, email, password_hash, role, created_at) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, email, password_hash, role, created_at
	`

	id := uuid.New().String()

	var userModel models.UserModel
	err := r.db.QueryRowxContext(
		ctx,
		query,
		id,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.CreatedAt,
	).StructScan(&userModel)

	if err != nil {
		return domain.User{}, err
	}

	return userModel.ToDomain(), nil
}
