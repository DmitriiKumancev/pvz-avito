package services

import (
	"context"
	"errors"

	"github.com/dkumancev/avito-pvz/pkg/application/repositories"
	"github.com/dkumancev/avito-pvz/pkg/domain"
)

type ReceptionService interface {
	CreateReception(ctx context.Context, pvzID string) (*domain.Reception, error)

	CloseReception(ctx context.Context, pvzID string) (*domain.Reception, error)

	AddProduct(ctx context.Context, pvzID string, productType string) (*domain.Product, error)

	RemoveLastProduct(ctx context.Context, pvzID string) error

	// Алиас для DeleteLastProduct - для совместимости с хендлерами
	DeleteLastProduct(ctx context.Context, pvzID string) error

	GetReceptionsByPVZID(ctx context.Context, pvzID string) ([]*domain.Reception, error)

	GetProductsByReceptionID(ctx context.Context, receptionID string) ([]*domain.Product, error)
}

type ReceptionServiceImpl struct {
	pvzRepo       repositories.PVZRepository
	receptionRepo repositories.ReceptionRepository
	productRepo   repositories.ProductRepository
}

func NewReceptionService(
	pvzRepo repositories.PVZRepository,
	receptionRepo repositories.ReceptionRepository,
	productRepo repositories.ProductRepository,
) ReceptionService {
	return &ReceptionServiceImpl{
		pvzRepo:       pvzRepo,
		receptionRepo: receptionRepo,
		productRepo:   productRepo,
	}
}

func (s *ReceptionServiceImpl) CreateReception(ctx context.Context, pvzID string) (*domain.Reception, error) {
	pvz, err := s.pvzRepo.GetByID(ctx, pvzID)
	if err != nil {
		return nil, errors.New("ПВЗ не найден")
	}

	// Проверяем, есть ли активная приемка для этого ПВЗ
	activeReception, err := s.receptionRepo.GetLastActiveByPVZID(ctx, pvzID)
	if err == nil && activeReception != nil {
		return nil, errors.New("для данного ПВЗ уже есть активная приемка товаров")
	}

	reception := domain.NewReception(pvz.ID)
	return s.receptionRepo.Create(ctx, reception)
}

func (s *ReceptionServiceImpl) CloseReception(ctx context.Context, pvzID string) (*domain.Reception, error) {
	_, err := s.pvzRepo.GetByID(ctx, pvzID)
	if err != nil {
		return nil, errors.New("ПВЗ не найден")
	}

	reception, err := s.receptionRepo.GetLastActiveByPVZID(ctx, pvzID)
	if err != nil || reception == nil {
		return nil, errors.New("активная приемка не найдена")
	}

	// Закрываем приемку
	err = reception.Close()
	if err != nil {
		return nil, err
	}

	// Обновляем в репо
	err = s.receptionRepo.Update(ctx, reception)
	if err != nil {
		return nil, err
	}

	return reception, nil
}

func (s *ReceptionServiceImpl) AddProduct(ctx context.Context, pvzID string, productType string) (*domain.Product, error) {
	_, err := s.pvzRepo.GetByID(ctx, pvzID)
	if err != nil {
		return nil, errors.New("ПВЗ не найден")
	}

	reception, err := s.receptionRepo.GetLastActiveByPVZID(ctx, pvzID)
	if err != nil || reception == nil {
		return nil, errors.New("активная приемка не найдена")
	}

	product, err := domain.NewProduct(productType, reception.ID)
	if err != nil {
		return nil, err
	}

	err = reception.AddProduct(*product)
	if err != nil {
		return nil, err
	}

	product, err = s.productRepo.Create(ctx, product, reception.ID)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ReceptionServiceImpl) RemoveLastProduct(ctx context.Context, pvzID string) error {
	_, err := s.pvzRepo.GetByID(ctx, pvzID)
	if err != nil {
		return errors.New("ПВЗ не найден")
	}

	reception, err := s.receptionRepo.GetLastActiveByPVZID(ctx, pvzID)
	if err != nil || reception == nil {
		return errors.New("активная приемка не найдена")
	}

	if !reception.IsActive() {
		return errors.New("нельзя удалить товар из закрытой приемки")
	}

	err = s.productRepo.DeleteLastByReceptionID(ctx, reception.ID)
	if err != nil {
		return err
	}

	return reception.RemoveLastProduct()
}

func (s *ReceptionServiceImpl) DeleteLastProduct(ctx context.Context, pvzID string) error {
	return s.RemoveLastProduct(ctx, pvzID)
}

func (s *ReceptionServiceImpl) GetReceptionsByPVZID(ctx context.Context, pvzID string) ([]*domain.Reception, error) {
	_, err := s.pvzRepo.GetByID(ctx, pvzID)
	if err != nil {
		return nil, errors.New("ПВЗ не найден")
	}

	return s.receptionRepo.GetByPVZID(ctx, pvzID)
}

func (s *ReceptionServiceImpl) GetProductsByReceptionID(ctx context.Context, receptionID string) ([]*domain.Product, error) {
	return s.productRepo.GetByReceptionID(ctx, receptionID)
}
