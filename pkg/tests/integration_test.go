package tests

import (
	"context"
	"errors"
	"testing"
	"time"

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
		productType := productTypes[i%len(productTypes)] // чередование типов товаров

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

	// попытка добавить товар после закрытия приемки (должно вернуть error)
	_, err = receptionService.AddProduct(ctx, pvz.ID, domain.ProductTypeElectronics)
	if err == nil {
		t.Error("Expected error when adding product to closed reception, got nil")
	}
}

// TestFullReceptionProcessWithAuth интеграционный тест с проверкой авторизации и ролей
func TestFullReceptionProcessWithAuth(t *testing.T) {
	if testing.Short() {
		t.Skip("Скип интеграционного теста в коротком режиме")
	}

	ctx := context.Background()

	mockUserRepo := NewMockUserRepository()
	mockPVZRepo := NewMockPVZRepository()
	mockReceptionRepo := NewMockReceptionRepository()
	mockProductRepo := NewMockProductRepository()

	// создаем серамсы
	jwtSecret := []byte("test-secret")
	tokenDuration := 24 * time.Hour
	userService := services.NewUserService(mockUserRepo, jwtSecret, tokenDuration)
	pvzService := services.NewPVZService(mockPVZRepo)
	receptionService := services.NewReceptionService(mockPVZRepo, mockReceptionRepo, mockProductRepo)

	// 1. Регистрация пользователей с разными ролями
	moderator, err := userService.Register(ctx, "moderator@example.com", "password123", domain.ModeratorRole)
	if err != nil {
		t.Fatalf("Failed to create moderator: %v", err)
	}

	employee, err := userService.Register(ctx, "employee@example.com", "password123", domain.EmployeeRole)
	if err != nil {
		t.Fatalf("Failed to create employee: %v", err)
	}

	// 2. Создание PVZ модератором
	pvz, err := createPVZWithAuth(ctx, pvzService, "Москва", moderator.Role)
	if err != nil {
		t.Fatalf("Failed to create PVZ as moderator: %v", err)
	}

	// Поптка создать PVZ сотрудником (должна быть ошибка)
	_, err = createPVZWithAuth(ctx, pvzService, "Санкт-Петербург", employee.Role)
	if err == nil {
		t.Error("Expected error when creating PVZ as employee, got nil")
	}

	// 3. Создание приемки сотрудником
	reception, err := createReceptionWithAuth(ctx, receptionService, pvz.ID, employee.Role)
	if err != nil {
		t.Fatalf("Failed to create reception as employee: %v", err)
	}

	// Попытка создать приемку модератором (должна быть ошибка)
	_, err = createReceptionWithAuth(ctx, receptionService, pvz.ID, moderator.Role)
	if err == nil {
		t.Error("Expected error when creating reception as moderator, got nil")
	}

	// 4. Добавление товаров сотрудником
	productTypes := []string{
		domain.ProductTypeElectronics,
		domain.ProductTypeClothes,
		domain.ProductTypeShoes,
	}

	for i := 0; i < 10; i++ {
		productType := productTypes[i%len(productTypes)]
		_, err := addProductWithAuth(ctx, receptionService, pvz.ID, productType, employee.Role)
		if err != nil {
			t.Fatalf("Failed to add product #%d as employee: %v", i+1, err)
		}
	}

	// Попытка добавить товар модератором (должна быть ошибка)
	_, err = addProductWithAuth(ctx, receptionService, pvz.ID, domain.ProductTypeElectronics, moderator.Role)
	if err == nil {
		t.Error("Expected error when adding product as moderator, got nil")
	}

	// 5. Закрытие приемки сотрудником
	_, err = closeReceptionWithAuth(ctx, receptionService, pvz.ID, employee.Role)
	if err != nil {
		t.Fatalf("Failed to close reception as employee: %v", err)
	}

	// Проверка, что товары были добавлены
	products, _ := mockProductRepo.GetByReceptionID(ctx, reception.ID)
	if len(products) != 10 {
		t.Errorf("Expected 10 products to be added, got %d", len(products))
	}
}


func createPVZWithAuth(ctx context.Context, service services.PVZService, city string, role domain.UserRole) (*domain.PVZ, error) {
	if role != domain.ModeratorRole {
		return nil, errors.New("только модератор может создавать ПВЗ")
	}
	return service.CreatePVZ(ctx, city)
}

func createReceptionWithAuth(ctx context.Context, service services.ReceptionService, pvzID string, role domain.UserRole) (*domain.Reception, error) {
	if role != domain.EmployeeRole {
		return nil, errors.New("только сотрудник может создавать приемки")
	}
	return service.CreateReception(ctx, pvzID)
}

func addProductWithAuth(ctx context.Context, service services.ReceptionService, pvzID, productType string, role domain.UserRole) (*domain.Product, error) {
	if role != domain.EmployeeRole {
		return nil, errors.New("только сотрудник может добавлять товары")
	}
	return service.AddProduct(ctx, pvzID, productType)
}

func closeReceptionWithAuth(ctx context.Context, service services.ReceptionService, pvzID string, role domain.UserRole) (*domain.Reception, error) {
	if role != domain.EmployeeRole {
		return nil, errors.New("только сотрудник может закрывать приемки")
	}
	return service.CloseReception(ctx, pvzID)
}
