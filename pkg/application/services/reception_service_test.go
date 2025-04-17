package services

import (
	"context"
	"errors"
	"testing"

	"github.com/dkumancev/avito-pvz/pkg/domain"
)

type MockReceptionRepository struct {
	receptions map[string]*domain.Reception
}

func NewMockReceptionRepository() *MockReceptionRepository {
	return &MockReceptionRepository{
		receptions: make(map[string]*domain.Reception),
	}
}

func (m *MockReceptionRepository) Create(ctx context.Context, reception *domain.Reception) (*domain.Reception, error) {
	reception.ID = "mock-reception-id"
	m.receptions[reception.ID] = reception
	return reception, nil
}

func (m *MockReceptionRepository) GetByID(ctx context.Context, id string) (*domain.Reception, error) {
	reception, ok := m.receptions[id]
	if !ok {
		return nil, errors.New("reception not found")
	}
	return reception, nil
}

func (m *MockReceptionRepository) Update(ctx context.Context, reception *domain.Reception) error {
	if _, ok := m.receptions[reception.ID]; !ok {
		return errors.New("reception not found")
	}
	m.receptions[reception.ID] = reception
	return nil
}

func (m *MockReceptionRepository) GetLastActiveByPVZID(ctx context.Context, pvzID string) (*domain.Reception, error) {
	for _, reception := range m.receptions {
		if reception.PVZID == pvzID && reception.IsActive() {
			return reception, nil
		}
	}
	return nil, errors.New("active reception not found")
}

func (m *MockReceptionRepository) GetByPVZID(ctx context.Context, pvzID string) ([]*domain.Reception, error) {
	var result []*domain.Reception
	for _, reception := range m.receptions {
		if reception.PVZID == pvzID {
			result = append(result, reception)
		}
	}
	return result, nil
}

type MockProductRepository struct {
	products map[string]*domain.Product
	// Для каждой приемки список ID товаров в порядке добавления (для LIFO)
	receptionProducts map[string][]string
}

func NewMockProductRepository() *MockProductRepository {
	return &MockProductRepository{
		products:          make(map[string]*domain.Product),
		receptionProducts: make(map[string][]string),
	}
}

func (m *MockProductRepository) Create(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	product.ID = "mock-product-id-" + product.Type
	m.products[product.ID] = product

	m.receptionProducts[product.ReceptionID] = append(m.receptionProducts[product.ReceptionID], product.ID)

	return product, nil
}

func (m *MockProductRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	product, ok := m.products[id]
	if !ok {
		return nil, errors.New("product not found")
	}
	return product, nil
}

func (m *MockProductRepository) GetByReceptionID(ctx context.Context, receptionID string) ([]*domain.Product, error) {
	var result []*domain.Product
	for _, productID := range m.receptionProducts[receptionID] {
		if product, ok := m.products[productID]; ok {
			result = append(result, product)
		}
	}
	return result, nil
}

func (m *MockProductRepository) DeleteLastByReceptionID(ctx context.Context, receptionID string) error {
	productIDs := m.receptionProducts[receptionID]
	if len(productIDs) == 0 {
		return errors.New("no products to delete")
	}

	lastProductID := productIDs[len(productIDs)-1]

	m.receptionProducts[receptionID] = productIDs[:len(productIDs)-1]

	delete(m.products, lastProductID)

	return nil
}

func TestReceptionService_CreateReception(t *testing.T) {
	ctx := context.Background()
	mockPVZRepo := NewMockPVZRepository()
	mockReceptionRepo := NewMockReceptionRepository()
	mockProductRepo := NewMockProductRepository()

	service := NewReceptionService(mockPVZRepo, mockReceptionRepo, mockProductRepo)

	pvz, _ := domain.NewPVZ("Москва")
	pvz.ID = "pvz-123"
	mockPVZRepo.pvzs[pvz.ID] = pvz

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
	mockPVZRepo := NewMockPVZRepository()
	mockReceptionRepo := NewMockReceptionRepository()
	mockProductRepo := NewMockProductRepository()

	service := NewReceptionService(mockPVZRepo, mockReceptionRepo, mockProductRepo)

	pvz, _ := domain.NewPVZ("Москва")
	pvz.ID = "pvz-123"
	mockPVZRepo.pvzs[pvz.ID] = pvz

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
	mockPVZRepo := NewMockPVZRepository()
	mockReceptionRepo := NewMockReceptionRepository()
	mockProductRepo := NewMockProductRepository()

	service := NewReceptionService(mockPVZRepo, mockReceptionRepo, mockProductRepo)

	pvz, _ := domain.NewPVZ("Москва")
	pvz.ID = "pvz-123"
	mockPVZRepo.pvzs[pvz.ID] = pvz

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
