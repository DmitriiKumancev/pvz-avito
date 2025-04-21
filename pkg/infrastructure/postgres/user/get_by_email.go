package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dkumancev/avito-pvz/pkg/domain"
	"github.com/dkumancev/avito-pvz/pkg/infrastructure/postgres/models"
)

// GetByEmail получает пользователя по email
func (r *Repository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	query := `
		SELECT id, email, password_hash, role, created_at
		FROM users 
		WHERE email = $1
	`

	var userModel models.UserModel
	err := r.db.QueryRowxContext(ctx, query, email).StructScan(&userModel)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, errors.New("user not found")
		}
		return domain.User{}, err
	}

	return userModel.ToDomain(), nil
}
