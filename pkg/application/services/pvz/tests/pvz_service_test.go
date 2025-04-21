package tests

import (
	"context"
	"testing"

	"github.com/dkumancev/avito-pvz/pkg/application/repositories"
	"github.com/dkumancev/avito-pvz/pkg/application/services"
	"github.com/dkumancev/avito-pvz/pkg/tests"
)

func TestPVZService_CreatePVZ(t *testing.T) {
	ctx := context.Background()
	mockRepo := tests.NewMockPVZRepository()
	service := services.NewPVZService(mockRepo)

	// Valid city
	pvz, err := service.CreatePVZ(ctx, "Москва")
	if err != nil {
		t.Errorf("Expected no error for valid city, got: %v", err)
	}
	if pvz == nil {
		t.Fatal("Expected PVZ to be created, got nil")
	}
	if pvz.City != "Москва" {
		t.Errorf("Expected city to be 'Москва', got %s", pvz.City)
	}

	// Invalid city
	_, err = service.CreatePVZ(ctx, "InvalidCity")
	if err == nil {
		t.Error("Expected error for invalid city, got nil")
	}
}

func TestPVZService_GetPVZ(t *testing.T) {
	ctx := context.Background()
	mockRepo := tests.NewMockPVZRepository()
	service := services.NewPVZService(mockRepo)

	// Create a PVZ first
	createdPVZ, _ := service.CreatePVZ(ctx, "Москва")

	// Valid ID
	pvz, err := service.GetPVZByID(ctx, createdPVZ.ID)
	if err != nil {
		t.Errorf("Expected no error for valid ID, got: %v", err)
	}
	if pvz == nil {
		t.Fatal("Expected PVZ to be returned, got nil")
	}
	if pvz.City != "Москва" {
		t.Errorf("Expected city to be 'Москва', got %s", pvz.City)
	}

	// Invalid ID
	_, err = service.GetPVZByID(ctx, "non-existent-id")
	if err == nil {
		t.Error("Expected error for non-existent ID, got nil")
	}
}

func TestPVZService_ListPVZs(t *testing.T) {
	ctx := context.Background()
	mockRepo := tests.NewMockPVZRepository()
	service := services.NewPVZService(mockRepo)

	// Create a PVZ - в текущей реализации мока, Create всегда использует
	// фиксированный ID "mock-pvz-id", так что второй вызов перезапишет первый
	_, _ = service.CreatePVZ(ctx, "Москва")

	// List all PVZs
	filter := repositories.PVZFilter{}
	pvzs, err := service.ListPVZs(ctx, filter)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if len(pvzs) != 1 {
		t.Errorf("Expected 1 PVZ, got %d", len(pvzs))
	}
}

// Тест для алиаса ListPVZ
func TestPVZService_ListPVZ(t *testing.T) {
	ctx := context.Background()
	mockRepo := tests.NewMockPVZRepository()
	service := services.NewPVZService(mockRepo)

	// Create a PVZ - в текущей реализации мока, Create всегда использует
	// фиксированный ID "mock-pvz-id", так что второй вызов перезапишет первый
	_, _ = service.CreatePVZ(ctx, "Москва")

	// Test alias
	filter := repositories.PVZFilter{}
	pvzs, err := service.ListPVZ(ctx, filter)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if len(pvzs) != 1 {
		t.Errorf("Expected 1 PVZ, got %d", len(pvzs))
	}
}
