package services

import (
	"context"
	"errors"
	"testing"

	"github.com/dkumancev/avito-pvz/pkg/application/repositories"
	"github.com/dkumancev/avito-pvz/pkg/domain"
)

type MockPVZRepository struct {
	pvzs map[string]*domain.PVZ
}

func NewMockPVZRepository() *MockPVZRepository {
	return &MockPVZRepository{
		pvzs: make(map[string]*domain.PVZ),
	}
}

func (m *MockPVZRepository) Create(ctx context.Context, pvz *domain.PVZ) (*domain.PVZ, error) {
	pvz.ID = "mock-pvz-id"
	m.pvzs[pvz.ID] = pvz
	return pvz, nil
}

func (m *MockPVZRepository) GetByID(ctx context.Context, id string) (*domain.PVZ, error) {
	pvz, ok := m.pvzs[id]
	if !ok {
		return nil, errors.New("pvz not found")
	}
	return pvz, nil
}

func (m *MockPVZRepository) List(ctx context.Context, filter repositories.PVZFilter) ([]*domain.PVZ, error) {
	result := make([]*domain.PVZ, 0, len(m.pvzs))
	for _, pvz := range m.pvzs {
		result = append(result, pvz)
	}
	return result, nil
}

func TestPVZService_CreatePVZ_ValidCity(t *testing.T) {
	ctx := context.Background()
	mockRepo := NewMockPVZRepository()
	service := NewPVZService(mockRepo)
	city := "Москва"

	// Act
	pvz, err := service.CreatePVZ(ctx, city)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if pvz == nil {
		t.Fatal("Expected PVZ to be created, got nil")
	}
	if pvz.City != city {
		t.Errorf("Expected city to be %s, got %s", city, pvz.City)
	}
	if pvz.ID == "" {
		t.Error("Expected ID to be set, got empty string")
	}
}

func TestPVZService_CreatePVZ_InvalidCity(t *testing.T) {
	ctx := context.Background()
	mockRepo := NewMockPVZRepository()
	service := NewPVZService(mockRepo)
	city := "Новосибирск"

	// Act
	pvz, err := service.CreatePVZ(ctx, city)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid city, got nil")
	}
	if pvz != nil {
		t.Errorf("Expected PVZ to be nil for invalid city, got %+v", pvz)
	}
}

func TestPVZService_GetPVZByID_Exists(t *testing.T) {
	ctx := context.Background()
	mockRepo := NewMockPVZRepository()
	service := NewPVZService(mockRepo)

	originalPVZ, _ := service.CreatePVZ(ctx, "Москва")

	// Act
	pvz, err := service.GetPVZByID(ctx, originalPVZ.ID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if pvz == nil {
		t.Fatal("Expected PVZ to be found, got nil")
	}
	if pvz.ID != originalPVZ.ID {
		t.Errorf("Expected ID to be %s, got %s", originalPVZ.ID, pvz.ID)
	}
	if pvz.City != originalPVZ.City {
		t.Errorf("Expected city to be %s, got %s", originalPVZ.City, pvz.City)
	}
}

func TestPVZService_GetPVZByID_NotExists(t *testing.T) {
	ctx := context.Background()
	mockRepo := NewMockPVZRepository()
	service := NewPVZService(mockRepo)

	// Act
	pvz, err := service.GetPVZByID(ctx, "non-existent-id")

	// Assert
	if err == nil {
		t.Error("Expected error for non-existent PVZ, got nil")
	}
	if pvz != nil {
		t.Errorf("Expected PVZ to be nil for non-existent ID, got %+v", pvz)
	}
}

func TestPVZService_ListPVZs(t *testing.T) {
	ctx := context.Background()
	mockRepo := NewMockPVZRepository()
	service := NewPVZService(mockRepo)

	// cоздаем несклько тестовых ПВЗ
	_, _ = service.CreatePVZ(ctx, "Москва")

	pvz2, _ := domain.NewPVZ("Казань")
	pvz2.ID = "mock-pvz-id-2"
	mockRepo.pvzs[pvz2.ID] = pvz2

	// Act
	pvzs, err := service.ListPVZs(ctx, repositories.PVZFilter{
		Page:  1,
		Limit: 10,
	})

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if len(pvzs) != 2 {
		t.Errorf("Expected 2 PVZs, got %d", len(pvzs))
	}
}
