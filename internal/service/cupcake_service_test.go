package service

import (
	"testing"

	"github.com/julimonteiro/cupcake-store/internal/models"
	"github.com/julimonteiro/cupcake-store/internal/repository"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Cupcake{})
	require.NoError(t, err)

	return db
}

func newTestService(t *testing.T) *CupcakeService {
	t.Helper()

	db := setupTestDB(t)
	repo := repository.NewCupcakeRepository(db)
	return NewCupcakeService(repo)
}

func TestCreate_SuccessAndDefaults(t *testing.T) {
	service := newTestService(t)

	req := &models.CreateCupcakeRequest{
		Name:       " Brigadeiro ",
		Flavor:     " Chocolate ",
		PriceCents: 1299,
	}

	cupcake, err := service.CreateCupcake(req)

	require.NoError(t, err)
	require.NotNil(t, cupcake)
	require.Greater(t, cupcake.ID, uint(0))
	require.Equal(t, "Brigadeiro", cupcake.Name)
	require.Equal(t, "Chocolate", cupcake.Flavor)
	require.Equal(t, 1299, cupcake.PriceCents)
	require.True(t, cupcake.IsAvailable)
}

func TestCreate_ValidationErrors(t *testing.T) {
	service := newTestService(t)

	tests := []struct {
		name        string
		req         *models.CreateCupcakeRequest
		expectedErr string
	}{
		{
			name: "name too short",
			req: &models.CreateCupcakeRequest{
				Name:       "A",
				Flavor:     "X",
				PriceCents: 1,
			},
			expectedErr: "name must have at least 2 characters",
		},
		{
			name: "empty flavor",
			req: &models.CreateCupcakeRequest{
				Name:       "Ok",
				Flavor:     "",
				PriceCents: 1,
			},
			expectedErr: "flavor is required",
		},
		{
			name: "invalid price",
			req: &models.CreateCupcakeRequest{
				Name:       "Ok",
				Flavor:     "Vanilla",
				PriceCents: 0,
			},
			expectedErr: "price must be greater than zero",
		},
		{
			name: "empty name",
			req: &models.CreateCupcakeRequest{
				Name:       "",
				Flavor:     "Chocolate",
				PriceCents: 1000,
			},
			expectedErr: "name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cupcake, err := service.CreateCupcake(tt.req)

			require.Error(t, err)
			require.Nil(t, cupcake)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestList_ReturnsInsertedItems(t *testing.T) {
	service := newTestService(t)

	req1 := &models.CreateCupcakeRequest{
		Name:       "Chocolate",
		Flavor:     "Belgian",
		PriceCents: 1500,
	}
	cupcake1, err := service.CreateCupcake(req1)
	require.NoError(t, err)

	req2 := &models.CreateCupcakeRequest{
		Name:       "Vanilla",
		Flavor:     "Madagascar",
		PriceCents: 1200,
	}
	cupcake2, err := service.CreateCupcake(req2)
	require.NoError(t, err)

	cupcakes, err := service.GetAllCupcakes()

	require.NoError(t, err)
	require.Len(t, cupcakes, 2)

	require.Equal(t, cupcake1.ID, cupcakes[0].ID)
	require.Equal(t, cupcake2.ID, cupcakes[1].ID)
	require.Equal(t, "Chocolate", cupcakes[0].Name)
	require.Equal(t, "Vanilla", cupcakes[1].Name)
}

func TestGet_Update_Delete(t *testing.T) {
	service := newTestService(t)

	req := &models.CreateCupcakeRequest{
		Name:       "Original",
		Flavor:     "Original",
		PriceCents: 1000,
	}
	created, err := service.CreateCupcake(req)
	require.NoError(t, err)

	retrieved, err := service.GetCupcake(created.ID)
	require.NoError(t, err)
	require.Equal(t, created.ID, retrieved.ID)
	require.Equal(t, "Original", retrieved.Name)
	require.Equal(t, "Original", retrieved.Flavor)
	require.Equal(t, 1000, retrieved.PriceCents)
	require.True(t, retrieved.IsAvailable)

	updateReq := &models.UpdateCupcakeRequest{
		Name:       stringPtr("Updated Name"),
		Flavor:     stringPtr("Updated Flavor"),
		PriceCents: intPtr(2000),
	}

	updated, err := service.UpdateCupcake(created.ID, updateReq)
	require.NoError(t, err)
	require.Equal(t, created.ID, updated.ID)
	require.Equal(t, "Updated Name", updated.Name)
	require.Equal(t, "Updated Flavor", updated.Flavor)
	require.Equal(t, 2000, updated.PriceCents)
	require.True(t, updated.IsAvailable)

	err = service.DeleteCupcake(created.ID)
	require.NoError(t, err)

	_, err = service.GetCupcake(created.ID)
	require.Error(t, err)
	require.Contains(t, err.Error(), "record not found")
}

func TestUpdate_ValidationErrors(t *testing.T) {
	service := newTestService(t)

	req := &models.CreateCupcakeRequest{
		Name:       "Test",
		Flavor:     "Test",
		PriceCents: 1000,
	}
	created, err := service.CreateCupcake(req)
	require.NoError(t, err)

	tests := []struct {
		name        string
		updateReq   *models.UpdateCupcakeRequest
		expectedErr string
	}{
		{
			name: "name too short",
			updateReq: &models.UpdateCupcakeRequest{
				Name: stringPtr("A"),
			},
			expectedErr: "name must have at least 2 characters",
		},
		{
			name: "invalid price",
			updateReq: &models.UpdateCupcakeRequest{
				PriceCents: intPtr(0),
			},
			expectedErr: "price must be greater than zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.UpdateCupcake(created.ID, tt.updateReq)

			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestDelete_NotFound(t *testing.T) {
	service := newTestService(t)

	err := service.DeleteCupcake(999)
	require.Error(t, err)
	require.Contains(t, err.Error(), "cupcake not found")
}

func TestGet_NotFound(t *testing.T) {
	service := newTestService(t)

	_, err := service.GetCupcake(999)
	require.Error(t, err)
	require.Contains(t, err.Error(), "record not found")
}

func TestUpdate_NotFound(t *testing.T) {
	service := newTestService(t)

	updateReq := &models.UpdateCupcakeRequest{
		Name: stringPtr("New Name"),
	}

	_, err := service.UpdateCupcake(999, updateReq)
	require.Error(t, err)
	require.Contains(t, err.Error(), "cupcake not found")
}

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
