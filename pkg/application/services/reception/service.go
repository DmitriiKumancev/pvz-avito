package reception

import (
	"context"

	"github.com/dkumancev/avito-pvz/pkg/application/repositories"
	"github.com/dkumancev/avito-pvz/pkg/domain"
)

type Service interface {
	// Создание новой приемки товаров на указанном ПВЗ
	CreateReception(ctx context.Context, pvzID string) (*domain.Reception, error)

	// Закрытие последней активной приемки на ПВЗ
	CloseReception(ctx context.Context, pvzID string) (*domain.Reception, error)

	// Добавление товара в рамках активной приемки на ПВЗ
	AddProduct(ctx context.Context, pvzID string, productType string) (*domain.Product, error)

	// Удаление последнего добавленного товара в рамках активной приемки
	RemoveLastProduct(ctx context.Context, pvzID string) error

	// Удаление последнего товара (алиас для RemoveLastProduct)
	DeleteLastProduct(ctx context.Context, pvzID string) error

	// Получение списка приемок по ID ПВЗ
	GetReceptionsByPVZID(ctx context.Context, pvzID string) ([]*domain.Reception, error)

	// Получение товаров по ID приемки
	GetProductsByReceptionID(ctx context.Context, receptionID string) ([]*domain.Product, error)
}

type service struct {
	pvzRepo       repositories.PVZRepository
	receptionRepo repositories.ReceptionRepository
	productRepo   repositories.ProductRepository
}

func New(
	pvzRepo repositories.PVZRepository,
	receptionRepo repositories.ReceptionRepository,
	productRepo repositories.ProductRepository,
) Service {
	return &service{
		pvzRepo:       pvzRepo,
		receptionRepo: receptionRepo,
		productRepo:   productRepo,
	}
}

func (s *service) getActiveReception(ctx context.Context, pvzID string) (*domain.Reception, error) {
	reception, err := s.receptionRepo.GetLastActiveByPVZID(ctx, pvzID)
	if err != nil {
		return nil, err
	}
	return reception, nil
}
