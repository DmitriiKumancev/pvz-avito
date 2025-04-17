package services

import (
	"context"
	"errors"

	"github.com/dkumancev/avito-pvz/pkg/application/repositories"
	"github.com/dkumancev/avito-pvz/pkg/domain"
)

type PVZService interface {
	CreatePVZ(ctx context.Context, city string) (*domain.PVZ, error)

	GetPVZByID(ctx context.Context, id string) (*domain.PVZ, error)

	ListPVZs(ctx context.Context, filter repositories.PVZFilter) ([]*domain.PVZ, error)
}

type PVZServiceImpl struct {
	pvzRepo repositories.PVZRepository
}

func NewPVZService(pvzRepo repositories.PVZRepository) PVZService {
	return &PVZServiceImpl{
		pvzRepo: pvzRepo,
	}
}

func (s *PVZServiceImpl) CreatePVZ(ctx context.Context, city string) (*domain.PVZ, error) {
	pvz, err := domain.NewPVZ(city)
	if err != nil {
		return nil, err
	}

	return s.pvzRepo.Create(ctx, pvz)
}

func (s *PVZServiceImpl) GetPVZByID(ctx context.Context, id string) (*domain.PVZ, error) {
	if id == "" {
		return nil, errors.New("ID ПВЗ не может быть пустым")
	}

	return s.pvzRepo.GetByID(ctx, id)
}

func (s *PVZServiceImpl) ListPVZs(ctx context.Context, filter repositories.PVZFilter) ([]*domain.PVZ, error) {
	// значения по умолчанию для пагинации
	if filter.Page <= 0 {
		filter.Page = 1
	}

	if filter.Limit <= 0 {
		filter.Limit = 10
	} else if filter.Limit > 30 {
		filter.Limit = 30
	}

	return s.pvzRepo.List(ctx, filter)
}
