package user

import "context"

// Exists проверяет, существует ли пользователь с данным email
func (r *Repository) Exists(ctx context.Context, email string) (bool, error) {
	query := `
		SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)
	`

	var exists bool
	err := r.db.QueryRowxContext(ctx, query, email).Scan(&exists)

	return exists, err
}
