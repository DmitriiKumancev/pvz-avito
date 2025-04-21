package service

import (
	"context"
	"testing"
	"time"

	"github.com/dkumancev/avito-pvz/pkg/application/repositories"
	"github.com/dkumancev/avito-pvz/pkg/domain"
	"github.com/dkumancev/avito-pvz/pkg/infrastructure/grpc/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPVZService struct {
	mock.Mock
}

func (m *MockPVZService) CreatePVZ(ctx context.Context, city string) (*domain.PVZ, error) {
	args := m.Called(ctx, city)
	return args.Get(0).(*domain.PVZ), args.Error(1)
}

func (m *MockPVZService) GetPVZByID(ctx context.Context, id string) (*domain.PVZ, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.PVZ), args.Error(1)
}

func (m *MockPVZService) ListPVZs(ctx context.Context, filter repositories.PVZFilter) ([]*domain.PVZ, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*domain.PVZ), args.Error(1)
}

func (m *MockPVZService) ListPVZ(ctx context.Context, filter repositories.PVZFilter) ([]*domain.PVZ, error) {
	return m.ListPVZs(ctx, filter)
}

func TestGetPVZList(t *testing.T) {
	mockService := new(MockPVZService)

	now := time.Now()
	samplePVZs := []*domain.PVZ{
		{
			ID:               "1",
			RegistrationDate: now,
			City:             "Москва",
		},
		{
			ID:               "2",
			RegistrationDate: now,
			City:             "Санкт-Петербург",
		},
	}

	mockService.On("ListPVZs", mock.Anything, mock.MatchedBy(func(filter repositories.PVZFilter) bool {
		return filter.Page == 1 && filter.Limit == 100
	})).Return(samplePVZs, nil)

	grpcService := NewPVZServiceServer(mockService)

	resp, err := grpcService.GetPVZList(context.Background(), &pb.GetPVZListRequest{})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Pvzs, 2)

	// Check the first PVZ
	assert.Equal(t, "1", resp.Pvzs[0].Id)
	assert.Equal(t, "Москва", resp.Pvzs[0].City)

	// Check the second PVZ
	assert.Equal(t, "2", resp.Pvzs[1].Id)
	assert.Equal(t, "Санкт-Петербург", resp.Pvzs[1].City)

	// Verify all expectations were met
	mockService.AssertExpectations(t)
}
