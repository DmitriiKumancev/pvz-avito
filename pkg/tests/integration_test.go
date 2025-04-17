package tests

import (
	"context"
	"testing"

	"github.com/dkumancev/avito-pvz/pkg/application/services"
	"github.com/dkumancev/avito-pvz/pkg/domain"
)

// Интеграционный тест для проверки полного процесса:
// 1. Создание ПВЗ
// 2. Добавление приемки
// 3. Добавление 50 товаров
// 4. Закрытие приемки
func TestFullReceptionProcess(t *testing.T) {
	if testing.Short() {
		t.Skip("Скип интеграционного теста в коротком режиме")
	}

	ctx := context.Background()

	mockPVZRepo := NewMockPVZRepository()
	mockReceptionRepo := NewMockReceptionRepository()
	mockProductRepo := NewMockProductRepository()

	pvzService := services.NewPVZService(mockPVZRepo)
	receptionService := services.NewReceptionService(mockPVZRepo, mockReceptionRepo, mockProductRepo)

	// Act & Assert

	// 1. Создание ПВЗ
	pvz, err := pvzService.CreatePVZ(ctx, "Москва")
	if err != nil {
		t.Fatalf("Failed to create PVZ: %v", err)
	}
	if pvz == nil {
		t.Fatal("Expected PVZ to be created, got nil")
	}

	// 2. Добавление приемки
	reception, err := receptionService.CreateReception(ctx, pvz.ID)
	if err != nil {
		t.Fatalf("Failed to create reception: %v", err)
	}
	if reception == nil {
		t.Fatal("Expected reception to be created, got nil")
	}
	if reception.Status != domain.ReceptionStatusInProgress {
		t.Errorf("Expected reception status to be %s, got %s", domain.ReceptionStatusInProgress, reception.Status)
	}

	// 3. Добавление 50 товаров
	productTypes := []string{
		domain.ProductTypeElectronics,
		domain.ProductTypeClothes,
		domain.ProductTypeShoes,
	}

	for i := 0; i < 50; i++ {
		productType := productTypes[i%len(productTypes)] // чережование типов товаров

		product, err := receptionService.AddProduct(ctx, pvz.ID, productType)
		if err != nil {
			t.Fatalf("Failed to add product #%d: %v", i+1, err)
		}
		if product == nil {
			t.Fatalf("Expected product #%d to be created, got nil", i+1)
		}
	}

	products, _ := mockProductRepo.GetByReceptionID(ctx, reception.ID)
	if len(products) != 50 {
		t.Errorf("Expected 50 products to be added, got %d", len(products))
	}

	// 4. Закрытие приемки
	closedReception, err := receptionService.CloseReception(ctx, pvz.ID)
	if err != nil {
		t.Fatalf("Failed to close reception: %v", err)
	}
	if closedReception == nil {
		t.Fatal("Expected reception to be returned after closing, got nil")
	}
	if closedReception.Status != domain.ReceptionStatusClosed {
		t.Errorf("Expected reception status to be %s after closing, got %s", domain.ReceptionStatusClosed, closedReception.Status)
	}

	// попытка добавить товар после закрытия приемки  (должно вернуть error)
	_, err = receptionService.AddProduct(ctx, pvz.ID, domain.ProductTypeElectronics)
	if err == nil {
		t.Error("Expected error when adding product to closed reception, got nil")
	}
}
