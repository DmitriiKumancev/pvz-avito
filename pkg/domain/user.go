package domain

import (
	"errors"
	"regexp"
	"time"
)

type UserRole string

const (
	// EmployeeRole сотрудник
	EmployeeRole UserRole = "employee"
	// ModeratorRole привелегированный пользователь (модератор)
	ModeratorRole UserRole = "moderator"
)

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         UserRole  `json:"role"`
	CreatedAt    time.Time `json:"createdAt"`
}

func NewUser(email string, passwordHash string, role UserRole) (User, error) {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return User{}, errors.New("invalid email format")
	}

	if passwordHash == "" {
		return User{}, errors.New("password hash cannot be empty")
	}

	if role != EmployeeRole && role != ModeratorRole {
		return User{}, errors.New("invalid role")
	}

	return User{
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
		CreatedAt:    time.Now(),
	}, nil
}
