package repositories

import (
	"context"

	"github.com/dkumancev/avito-pvz/pkg/domain"
)

type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) (*domain.Product, error)

	GetByID(ctx context.Context, id string) (*domain.Product, error)

	GetByReceptionID(ctx context.Context, receptionID string) ([]*domain.Product, error)

	DeleteLastByReceptionID(ctx context.Context, receptionID string) error
}
