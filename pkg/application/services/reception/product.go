package reception

import (
	"context"
	"fmt"

	"github.com/dkumancev/avito-pvz/pkg/domain"
)

func (s *service) AddProduct(ctx context.Context, pvzID string, productType string) (*domain.Product, error) {
	_, err := s.pvzRepo.GetByID(ctx, pvzID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения ПВЗ: %w", err)
	}

	reception, err := s.getActiveReception(ctx, pvzID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить активную приемку: %w", err)
	}

	product, err := domain.NewProduct(productType, reception.ID)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания товара: %w", err)
	}

	err = reception.AddProduct(*product)
	if err != nil {
		return nil, fmt.Errorf("не удалось добавить товар в приемку: %w", err)
	}

	savedProduct, err := s.productRepo.Create(ctx, product, reception.ID)
	if err != nil {
		return nil, fmt.Errorf("ошибка сохранения товара: %w", err)
	}

	return savedProduct, nil
}

func (s *service) RemoveLastProduct(ctx context.Context, pvzID string) error {
	_, err := s.pvzRepo.GetByID(ctx, pvzID)
	if err != nil {
		return fmt.Errorf("ошибка получения ПВЗ: %w", err)
	}

	reception, err := s.getActiveReception(ctx, pvzID)
	if err != nil {
		return fmt.Errorf("не удалось получить активную приемку: %w", err)
	}

	err = reception.RemoveLastProduct()
	if err != nil {
		return fmt.Errorf("ошибка удаления товара: %w", err)
	}

	err = s.productRepo.DeleteLastByReceptionID(ctx, reception.ID)
	if err != nil {
		return fmt.Errorf("ошибка удаления товара из БД: %w", err)
	}

	return nil
}

func (s *service) DeleteLastProduct(ctx context.Context, pvzID string) error {
	return s.RemoveLastProduct(ctx, pvzID)
}
