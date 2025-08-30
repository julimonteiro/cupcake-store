package service

import (
	"testing"

	"github.com/julimonteiro/cupcake-store/internal/models"
	"github.com/julimonteiro/cupcake-store/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCupcakeRepository struct {
	mock.Mock
}

var _ repository.CupcakeRepositoryInterface = (*MockCupcakeRepository)(nil)

func (m *MockCupcakeRepository) Create(cupcake *models.Cupcake) error {
	args := m.Called(cupcake)
	return args.Error(0)
}

func (m *MockCupcakeRepository) FindByID(id uint) (*models.Cupcake, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Cupcake), args.Error(1)
}

func (m *MockCupcakeRepository) FindAll() ([]models.Cupcake, error) {
	args := m.Called()
	return args.Get(0).([]models.Cupcake), args.Error(1)
}

func (m *MockCupcakeRepository) Update(cupcake *models.Cupcake) error {
	args := m.Called(cupcake)
	return args.Error(0)
}

func (m *MockCupcakeRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockCupcakeRepository) Exists(id uint) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}

func TestCreateCupcake_ValidRequest(t *testing.T) {
	mockRepo := new(MockCupcakeRepository)
	service := NewCupcakeService(mockRepo)

	req := &models.CreateCupcakeRequest{
		Name:       "Special Chocolate",
		Flavor:     "Belgian Chocolate",
		PriceCents: 1500,
	}

	expectedCupcake := &models.Cupcake{
		Name:        "Special Chocolate",
		Flavor:      "Belgian Chocolate",
		PriceCents:  1500,
		IsAvailable: true,
	}

	mockRepo.On("Create", mock.AnythingOfType("*models.Cupcake")).Return(nil)

	result, err := service.CreateCupcake(req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCupcake.Name, result.Name)
	assert.Equal(t, expectedCupcake.Flavor, result.Flavor)
	assert.Equal(t, expectedCupcake.PriceCents, result.PriceCents)
	assert.Equal(t, expectedCupcake.IsAvailable, result.IsAvailable)
	mockRepo.AssertExpectations(t)
}

func TestCreateCupcake_InvalidName(t *testing.T) {
	mockRepo := new(MockCupcakeRepository)
	service := NewCupcakeService(mockRepo)

	req := &models.CreateCupcakeRequest{
		Name:       "A",
		Flavor:     "Belgian Chocolate",
		PriceCents: 1500,
	}

	result, err := service.CreateCupcake(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "name must have at least 2 characters")
}

func TestCreateCupcake_EmptyName(t *testing.T) {
	mockRepo := new(MockCupcakeRepository)
	service := NewCupcakeService(mockRepo)

	req := &models.CreateCupcakeRequest{
		Name:       "",
		Flavor:     "Belgian Chocolate",
		PriceCents: 1500,
	}

	result, err := service.CreateCupcake(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "name is required")
}

func TestCreateCupcake_EmptyFlavor(t *testing.T) {
	mockRepo := new(MockCupcakeRepository)
	service := NewCupcakeService(mockRepo)

	req := &models.CreateCupcakeRequest{
		Name:       "Special Chocolate",
		Flavor:     "",
		PriceCents: 1500,
	}

	result, err := service.CreateCupcake(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "flavor is required")
}

func TestCreateCupcake_InvalidPrice(t *testing.T) {
	mockRepo := new(MockCupcakeRepository)
	service := NewCupcakeService(mockRepo)

	req := &models.CreateCupcakeRequest{
		Name:       "Special Chocolate",
		Flavor:     "Belgian Chocolate",
		PriceCents: 0,
	}

	result, err := service.CreateCupcake(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "price must be greater than zero")
}

func TestGetCupcake_Success(t *testing.T) {
	mockRepo := new(MockCupcakeRepository)
	service := NewCupcakeService(mockRepo)

	expectedCupcake := &models.Cupcake{
		ID:          1,
		Name:        "Special Chocolate",
		Flavor:      "Belgian Chocolate",
		PriceCents:  1500,
		IsAvailable: true,
	}

	mockRepo.On("FindByID", uint(1)).Return(expectedCupcake, nil)

	result, err := service.GetCupcake(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCupcake.ID, result.ID)
	assert.Equal(t, expectedCupcake.Name, result.Name)
	mockRepo.AssertExpectations(t)
}

func TestGetAllCupcakes_Success(t *testing.T) {
	mockRepo := new(MockCupcakeRepository)
	service := NewCupcakeService(mockRepo)

	expectedCupcakes := []models.Cupcake{
		{ID: 1, Name: "Special Chocolate", Flavor: "Belgian Chocolate", PriceCents: 1500, IsAvailable: true},
		{ID: 2, Name: "Vanilla", Flavor: "Vanilla", PriceCents: 1200, IsAvailable: true},
	}

	mockRepo.On("FindAll").Return(expectedCupcakes, nil)

	result, err := service.GetAllCupcakes()

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, expectedCupcakes[0].Name, result[0].Name)
	assert.Equal(t, expectedCupcakes[1].Name, result[1].Name)
	mockRepo.AssertExpectations(t)
}

func TestDeleteCupcake_Success(t *testing.T) {
	mockRepo := new(MockCupcakeRepository)
	service := NewCupcakeService(mockRepo)

	mockRepo.On("Exists", uint(1)).Return(true, nil)
	mockRepo.On("Delete", uint(1)).Return(nil)

	err := service.DeleteCupcake(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteCupcake_NotFound(t *testing.T) {
	mockRepo := new(MockCupcakeRepository)
	service := NewCupcakeService(mockRepo)

	mockRepo.On("Exists", uint(999)).Return(false, nil)

	err := service.DeleteCupcake(999)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cupcake not found")
	mockRepo.AssertExpectations(t)
}
