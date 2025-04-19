package services

import (
	"context"
	"testing"

	"github.com/dkumancev/avito-pvz/pkg/domain"
	"github.com/dkumancev/avito-pvz/pkg/tests"
)

func TestReceptionService_CreateReception(t *testing.T) {
	ctx := context.Background()
	mockPVZRepo := tests.NewMockPVZRepository()
	mockReceptionRepo := tests.NewMockReceptionRepository()
	mockProductRepo := tests.NewMockProductRepository()

	service := NewReceptionService(mockPVZRepo, mockReceptionRepo, mockProductRepo)

	pvz, _ := domain.NewPVZ("Москва")
	pvz.ID = "pvz-123"
	mockPVZRepo.Create(ctx, pvz)

	// Act
	reception, err := service.CreateReception(ctx, pvz.ID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if reception == nil {
		t.Fatal("Expected reception to be created, got nil")
	}
	if reception.PVZID != pvz.ID {
		t.Errorf("Expected PVZID to be %s, got %s", pvz.ID, reception.PVZID)
	}
	if reception.Status != domain.ReceptionStatusInProgress {
		t.Errorf("Expected status to be %s, got %s", domain.ReceptionStatusInProgress, reception.Status)
	}

	// Попытка создать 2-ую активную приемку должна вернуть error
	_, err = service.CreateReception(ctx, pvz.ID)
	if err == nil {
		t.Error("Expected error when creating second active reception, got nil")
	}
}

func TestReceptionService_CloseReception(t *testing.T) {
	ctx := context.Background()
	mockPVZRepo := tests.NewMockPVZRepository()
	mockReceptionRepo := tests.NewMockReceptionRepository()
	mockProductRepo := tests.NewMockProductRepository()

	service := NewReceptionService(mockPVZRepo, mockReceptionRepo, mockProductRepo)

	pvz, _ := domain.NewPVZ("Москва")
	pvz.ID = "pvz-123"
	mockPVZRepo.Create(ctx, pvz)

	_, _ = service.CreateReception(ctx, pvz.ID)

	// Act
	closedReception, err := service.CloseReception(ctx, pvz.ID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if closedReception == nil {
		t.Fatal("Expected reception to be returned, got nil")
	}
	if closedReception.Status != domain.ReceptionStatusClosed {
		t.Errorf("Expected status to be %s, got %s", domain.ReceptionStatusClosed, closedReception.Status)
	}

	// если попытаться закрыть несуществующую приемку - должно вернуть error
	_, err = service.CloseReception(ctx, "non-existent-pvz")
	if err == nil {
		t.Error("Expected error when closing reception for non-existent PVZ, got nil")
	}
}

func TestReceptionService_AddAndRemoveProduct(t *testing.T) {
	ctx := context.Background()
	mockPVZRepo := tests.NewMockPVZRepository()
	mockReceptionRepo := tests.NewMockReceptionRepository()
	mockProductRepo := tests.NewMockProductRepository()

	service := NewReceptionService(mockPVZRepo, mockReceptionRepo, mockProductRepo)

	pvz, _ := domain.NewPVZ("Москва")
	pvz.ID = "pvz-123"
	mockPVZRepo.Create(ctx, pvz)

	reception, _ := service.CreateReception(ctx, pvz.ID)

	product, err := service.AddProduct(ctx, pvz.ID, domain.ProductTypeElectronics)

	// Assert
	if err != nil {
		t.Errorf("Expected no error when adding product, got: %v", err)
	}
	if product == nil {
		t.Fatal("Expected product to be created, got nil")
	}
	if product.Type != domain.ProductTypeElectronics {
		t.Errorf("Expected product type to be %s, got %s", domain.ProductTypeElectronics, product.Type)
	}

	// добавляем второй товар
	_, _ = service.AddProduct(ctx, pvz.ID, domain.ProductTypeClothes)

	// Act - удаляем последний товар (LIFO)
	err = service.RemoveLastProduct(ctx, pvz.ID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error when removing last product, got: %v", err)
	}

	// check что товары в правильном порядке
	products, _ := mockProductRepo.GetByReceptionID(ctx, reception.ID)
	if len(products) != 1 {
		t.Errorf("Expected 1 product after removal, got %d", len(products))
	}
	if len(products) > 0 && products[0].Type != domain.ProductTypeElectronics {
		t.Errorf("Expected remaining product to be %s, got %s", domain.ProductTypeElectronics, products[0].Type)
	}

	_, _ = service.CloseReception(ctx, pvz.ID)

	// попытка добавить товар в закрытую приемку должна вернуть ошибку
	_, err = service.AddProduct(ctx, pvz.ID, domain.ProductTypeShoes)
	if err == nil {
		t.Error("Expected error when adding product to closed reception, got nil")
	}

	err = service.RemoveLastProduct(ctx, pvz.ID)
	if err == nil {
		t.Error("Expected error when removing product from closed reception, got nil")
	}
}

func TestReceptionService_GetReceptionsAndProducts(t *testing.T) {
	ctx := context.Background()
	mockPVZRepo := tests.NewMockPVZRepository()
	mockReceptionRepo := tests.NewMockReceptionRepository()
	mockProductRepo := tests.NewMockProductRepository()

	service := NewReceptionService(mockPVZRepo, mockReceptionRepo, mockProductRepo)

	// Подготавливаем данные
	pvz, _ := domain.NewPVZ("Москва")
	pvz.ID = "pvz-123"
	mockPVZRepo.Create(ctx, pvz)

	reception, _ := service.CreateReception(ctx, pvz.ID)
	_, _ = service.AddProduct(ctx, pvz.ID, domain.ProductTypeElectronics)
	_, _ = service.AddProduct(ctx, pvz.ID, domain.ProductTypeClothes)

	receptions, err := service.GetReceptionsByPVZID(ctx, pvz.ID)
	if err != nil {
		t.Errorf("Expected no error when getting receptions, got: %v", err)
	}
	if len(receptions) != 1 {
		t.Errorf("Expected 1 reception, got %d", len(receptions))
	}

	products, err := service.GetProductsByReceptionID(ctx, reception.ID)
	if err != nil {
		t.Errorf("Expected no error when getting products, got: %v", err)
	}
	if len(products) != 2 {
		t.Errorf("Expected 2 products, got %d", len(products))
	}

	err = service.DeleteLastProduct(ctx, pvz.ID)
	if err != nil {
		t.Errorf("Expected no error when deleting last product, got: %v", err)
	}

	products, _ = service.GetProductsByReceptionID(ctx, reception.ID)
	if len(products) != 1 {
		t.Errorf("Expected 1 product after deletion, got %d", len(products))
	}
}
