package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dkumancev/avito-pvz/pkg/application/repositories"
	"github.com/dkumancev/avito-pvz/pkg/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) repositories.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	query := `
		INSERT INTO users (id, email, password_hash, role, created_at) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, email, password_hash, role, created_at
	`

	id := uuid.New().String()

	var userModel UserModel
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

func (r *userRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	query := `
		SELECT id, email, password_hash, role, created_at
		FROM users 
		WHERE email = $1
	`

	var userModel UserModel
	err := r.db.QueryRowxContext(ctx, query, email).StructScan(&userModel)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, errors.New("user not found")
		}
		return domain.User{}, err
	}

	return userModel.ToDomain(), nil
}

func (r *userRepository) Exists(ctx context.Context, email string) (bool, error) {
	query := `
		SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)
	`

	var exists bool
	err := r.db.QueryRowxContext(ctx, query, email).Scan(&exists)

	return exists, err
}
