package tests

import (
	"context"
	"errors"

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

// Update обновляет информацию о приемке
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

	// get ID last added product
	lastProductID := productIDs[len(productIDs)-1]

	// delete last product from reception
	m.receptionProducts[receptionID] = productIDs[:len(productIDs)-1]

	// Удаляем из общего списка товаров
	delete(m.products, lastProductID)

	return nil
}
