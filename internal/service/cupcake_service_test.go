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

func TestCreateCupcake(t *testing.T) {
	tests := []struct {
		name             string
		request          *models.CreateCupcakeRequest
		expectedError    string
		validateResponse func(t *testing.T, cupcake *models.Cupcake)
	}{
		{
			name: "success with defaults",
			request: &models.CreateCupcakeRequest{
				Name:       " Brigadeiro ",
				Flavor:     " Chocolate ",
				PriceCents: 1299,
			},
			validateResponse: func(t *testing.T, cupcake *models.Cupcake) {
				require.Greater(t, cupcake.ID, uint(0))
				require.Equal(t, "Brigadeiro", cupcake.Name)
				require.Equal(t, "Chocolate", cupcake.Flavor)
				require.Equal(t, 1299, cupcake.PriceCents)
				require.True(t, cupcake.IsAvailable)
			},
		},
		{
			name: "success with default availability",
			request: &models.CreateCupcakeRequest{
				Name:       "Vanilla Special",
				Flavor:     "Madagascar Vanilla",
				PriceCents: 1500,
			},
			validateResponse: func(t *testing.T, cupcake *models.Cupcake) {
				require.Greater(t, cupcake.ID, uint(0))
				require.Equal(t, "Vanilla Special", cupcake.Name)
				require.Equal(t, "Madagascar Vanilla", cupcake.Flavor)
				require.Equal(t, 1500, cupcake.PriceCents)
				require.True(t, cupcake.IsAvailable)
			},
		},
		{
			name: "validation error - name too short",
			request: &models.CreateCupcakeRequest{
				Name:       "A",
				Flavor:     "X",
				PriceCents: 1,
			},
			expectedError: "name must have at least 2 characters",
		},
		{
			name: "validation error - empty flavor",
			request: &models.CreateCupcakeRequest{
				Name:       "Valid Name",
				Flavor:     "",
				PriceCents: 1000,
			},
			expectedError: "flavor is required",
		},
		{
			name: "validation error - zero price",
			request: &models.CreateCupcakeRequest{
				Name:       "Valid Name",
				Flavor:     "Valid Flavor",
				PriceCents: 0,
			},
			expectedError: "price must be greater than zero",
		},
		{
			name: "validation error - negative price",
			request: &models.CreateCupcakeRequest{
				Name:       "Valid Name",
				Flavor:     "Valid Flavor",
				PriceCents: -100,
			},
			expectedError: "price must be greater than zero",
		},
		{
			name: "validation error - empty name",
			request: &models.CreateCupcakeRequest{
				Name:       "",
				Flavor:     "Valid Flavor",
				PriceCents: 1000,
			},
			expectedError: "name is required",
		},
		{
			name: "validation error - empty flavor with spaces",
			request: &models.CreateCupcakeRequest{
				Name:       "Valid Name",
				Flavor:     "   ",
				PriceCents: 1000,
			},
			expectedError: "flavor is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := newTestService(t)

			cupcake, err := service.CreateCupcake(tt.request)

			if tt.expectedError != "" {
				require.Error(t, err)
				require.Nil(t, cupcake)
				require.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
				require.NotNil(t, cupcake)
				if tt.validateResponse != nil {
					tt.validateResponse(t, cupcake)
				}
			}
		})
	}
}

func TestGetCupcake(t *testing.T) {
	tests := []struct {
		name             string
		cupcakeID        uint
		setupCupcake     *models.CreateCupcakeRequest
		expectedError    string
		validateResponse func(t *testing.T, cupcake *models.Cupcake)
	}{
		{
			name:      "success - existing cupcake",
			cupcakeID: 1,
			setupCupcake: &models.CreateCupcakeRequest{
				Name:       "Chocolate Special",
				Flavor:     "Belgian Chocolate",
				PriceCents: 1500,
			},
			validateResponse: func(t *testing.T, cupcake *models.Cupcake) {
				require.Equal(t, uint(1), cupcake.ID)
				require.Equal(t, "Chocolate Special", cupcake.Name)
				require.Equal(t, "Belgian Chocolate", cupcake.Flavor)
				require.Equal(t, 1500, cupcake.PriceCents)
			},
		},
		{
			name:          "error - non-existent cupcake",
			cupcakeID:     999,
			expectedError: "record not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := newTestService(t)

			if tt.setupCupcake != nil {
				createdCupcake, err := service.CreateCupcake(tt.setupCupcake)
				require.NoError(t, err)
				tt.cupcakeID = createdCupcake.ID
			}

			cupcake, err := service.GetCupcake(tt.cupcakeID)

			if tt.expectedError != "" {
				require.Error(t, err)
				require.Nil(t, cupcake)
				require.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
				require.NotNil(t, cupcake)
				if tt.validateResponse != nil {
					tt.validateResponse(t, cupcake)
				}
			}
		})
	}
}

func TestGetAllCupcakes(t *testing.T) {
	tests := []struct {
		name             string
		setupCupcakes    []*models.CreateCupcakeRequest
		expectedCount    int
		validateResponse func(t *testing.T, cupcakes []models.Cupcake)
	}{
		{
			name:          "empty list",
			setupCupcakes: []*models.CreateCupcakeRequest{},
			expectedCount: 0,
			validateResponse: func(t *testing.T, cupcakes []models.Cupcake) {
				require.Len(t, cupcakes, 0)
			},
		},
		{
			name: "single cupcake",
			setupCupcakes: []*models.CreateCupcakeRequest{
				{
					Name:       "Chocolate",
					Flavor:     "Belgian",
					PriceCents: 1500,
				},
			},
			expectedCount: 1,
			validateResponse: func(t *testing.T, cupcakes []models.Cupcake) {
				require.Len(t, cupcakes, 1)
				require.Equal(t, "Chocolate", cupcakes[0].Name)
				require.Equal(t, "Belgian", cupcakes[0].Flavor)
			},
		},
		{
			name: "multiple cupcakes",
			setupCupcakes: []*models.CreateCupcakeRequest{
				{
					Name:       "Chocolate",
					Flavor:     "Belgian",
					PriceCents: 1500,
				},
				{
					Name:       "Vanilla",
					Flavor:     "Madagascar",
					PriceCents: 1200,
				},
			},
			expectedCount: 2,
			validateResponse: func(t *testing.T, cupcakes []models.Cupcake) {
				require.Len(t, cupcakes, 2)
				require.Equal(t, "Chocolate", cupcakes[0].Name)
				require.Equal(t, "Vanilla", cupcakes[1].Name)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := newTestService(t)

			for _, cupcakeReq := range tt.setupCupcakes {
				_, err := service.CreateCupcake(cupcakeReq)
				require.NoError(t, err)
			}

			cupcakes, err := service.GetAllCupcakes()

			require.NoError(t, err)
			require.Len(t, cupcakes, tt.expectedCount)

			if tt.validateResponse != nil {
				tt.validateResponse(t, cupcakes)
			}
		})
	}
}

func TestUpdateCupcake(t *testing.T) {
	tests := []struct {
		name             string
		cupcakeID        uint
		updateRequest    *models.UpdateCupcakeRequest
		setupCupcake     *models.CreateCupcakeRequest
		expectedError    string
		validateResponse func(t *testing.T, cupcake *models.Cupcake)
	}{
		{
			name:      "success - update all fields",
			cupcakeID: 1,
			setupCupcake: &models.CreateCupcakeRequest{
				Name:       "Original Name",
				Flavor:     "Original Flavor",
				PriceCents: 1000,
			},
			updateRequest: &models.UpdateCupcakeRequest{
				Name:        stringPtr("Updated Name"),
				Flavor:      stringPtr("Updated Flavor"),
				PriceCents:  intPtr(2000),
				IsAvailable: boolPtr(false),
			},
			validateResponse: func(t *testing.T, cupcake *models.Cupcake) {
				require.Equal(t, "Updated Name", cupcake.Name)
				require.Equal(t, "Updated Flavor", cupcake.Flavor)
				require.Equal(t, 2000, cupcake.PriceCents)
				require.False(t, cupcake.IsAvailable)
			},
		},
		{
			name:      "success - partial update",
			cupcakeID: 1,
			setupCupcake: &models.CreateCupcakeRequest{
				Name:       "Original Name",
				Flavor:     "Original Flavor",
				PriceCents: 1000,
			},
			updateRequest: &models.UpdateCupcakeRequest{
				Name: stringPtr("Updated Name Only"),
			},
			validateResponse: func(t *testing.T, cupcake *models.Cupcake) {
				require.Equal(t, "Updated Name Only", cupcake.Name)
				require.Equal(t, "Original Flavor", cupcake.Flavor)
				require.Equal(t, 1000, cupcake.PriceCents)
			},
		},
		{
			name:      "success - update with trimming",
			cupcakeID: 1,
			setupCupcake: &models.CreateCupcakeRequest{
				Name:       "Original Name",
				Flavor:     "Original Flavor",
				PriceCents: 1000,
			},
			updateRequest: &models.UpdateCupcakeRequest{
				Name:   stringPtr("  Updated Name  "),
				Flavor: stringPtr("  Updated Flavor  "),
			},
			validateResponse: func(t *testing.T, cupcake *models.Cupcake) {
				require.Equal(t, "Updated Name", cupcake.Name)
				require.Equal(t, "Updated Flavor", cupcake.Flavor)
			},
		},
		{
			name:          "error - non-existent cupcake",
			cupcakeID:     999,
			updateRequest: &models.UpdateCupcakeRequest{Name: stringPtr("Updated")},
			expectedError: "record not found",
		},
		{
			name:      "validation error - name too short",
			cupcakeID: 1,
			setupCupcake: &models.CreateCupcakeRequest{
				Name:       "Original Name",
				Flavor:     "Original Flavor",
				PriceCents: 1000,
			},
			updateRequest: &models.UpdateCupcakeRequest{
				Name: stringPtr("A"),
			},
			expectedError: "name must have at least 2 characters",
		},
		{
			name:      "validation error - zero price",
			cupcakeID: 1,
			setupCupcake: &models.CreateCupcakeRequest{
				Name:       "Original Name",
				Flavor:     "Original Flavor",
				PriceCents: 1000,
			},
			updateRequest: &models.UpdateCupcakeRequest{
				PriceCents: intPtr(0),
			},
			expectedError: "price must be greater than zero",
		},
		{
			name:      "validation error - negative price",
			cupcakeID: 1,
			setupCupcake: &models.CreateCupcakeRequest{
				Name:       "Original Name",
				Flavor:     "Original Flavor",
				PriceCents: 1000,
			},
			updateRequest: &models.UpdateCupcakeRequest{
				PriceCents: intPtr(-100),
			},
			expectedError: "price must be greater than zero",
		},
		{
			name:      "success - empty flavor with spaces",
			cupcakeID: 1,
			setupCupcake: &models.CreateCupcakeRequest{
				Name:       "Original Name",
				Flavor:     "Original Flavor",
				PriceCents: 1000,
			},
			updateRequest: &models.UpdateCupcakeRequest{
				Flavor: stringPtr("   "),
			},
			validateResponse: func(t *testing.T, cupcake *models.Cupcake) {
				require.Equal(t, "", cupcake.Flavor)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := newTestService(t)

			if tt.setupCupcake != nil {
				createdCupcake, err := service.CreateCupcake(tt.setupCupcake)
				require.NoError(t, err)
				tt.cupcakeID = createdCupcake.ID
			}

			cupcake, err := service.UpdateCupcake(tt.cupcakeID, tt.updateRequest)

			if tt.expectedError != "" {
				require.Error(t, err)
				require.Nil(t, cupcake)
				require.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
				require.NotNil(t, cupcake)
				if tt.validateResponse != nil {
					tt.validateResponse(t, cupcake)
				}
			}
		})
	}
}

func TestDeleteCupcake(t *testing.T) {
	tests := []struct {
		name          string
		cupcakeID     uint
		setupCupcake  *models.CreateCupcakeRequest
		expectedError string
	}{
		{
			name:      "success - existing cupcake",
			cupcakeID: 1,
			setupCupcake: &models.CreateCupcakeRequest{
				Name:       "To Delete",
				Flavor:     "Test Flavor",
				PriceCents: 1000,
			},
		},
		{
			name:          "error - non-existent cupcake",
			cupcakeID:     999,
			expectedError: "record not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := newTestService(t)

			if tt.setupCupcake != nil {
				createdCupcake, err := service.CreateCupcake(tt.setupCupcake)
				require.NoError(t, err)
				tt.cupcakeID = createdCupcake.ID
			}

			err := service.DeleteCupcake(tt.cupcakeID)

			if tt.expectedError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCreateCupcake_RepositoryError(t *testing.T) {
	tests := []struct {
		name          string
		request       *models.CreateCupcakeRequest
		expectedError string
	}{
		{
			name: "repository error handling",
			request: &models.CreateCupcakeRequest{
				Name:       "Valid Name",
				Flavor:     "Valid Flavor",
				PriceCents: 1000,
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := newTestService(t)

			cupcake, err := service.CreateCupcake(tt.request)

			if tt.expectedError != "" {
				require.Error(t, err)
				require.Nil(t, cupcake)
				require.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
				require.NotNil(t, cupcake)
			}
		})
	}
}

func TestUpdateCupcake_RepositoryError(t *testing.T) {
	tests := []struct {
		name          string
		cupcakeID     uint
		updateRequest *models.UpdateCupcakeRequest
		setupCupcake  *models.CreateCupcakeRequest
		expectedError string
	}{
		{
			name:      "repository error handling",
			cupcakeID: 1,
			setupCupcake: &models.CreateCupcakeRequest{
				Name:       "Original Name",
				Flavor:     "Original Flavor",
				PriceCents: 1000,
			},
			updateRequest: &models.UpdateCupcakeRequest{
				Name: stringPtr("Updated Name"),
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := newTestService(t)

			if tt.setupCupcake != nil {
				createdCupcake, err := service.CreateCupcake(tt.setupCupcake)
				require.NoError(t, err)
				tt.cupcakeID = createdCupcake.ID
			}

			cupcake, err := service.UpdateCupcake(tt.cupcakeID, tt.updateRequest)

			if tt.expectedError != "" {
				require.Error(t, err)
				require.Nil(t, cupcake)
				require.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
				require.NotNil(t, cupcake)
			}
		})
	}
}

func TestDeleteCupcake_RepositoryError(t *testing.T) {
	tests := []struct {
		name          string
		cupcakeID     uint
		setupCupcake  *models.CreateCupcakeRequest
		expectedError string
	}{
		{
			name:      "repository error handling",
			cupcakeID: 1,
			setupCupcake: &models.CreateCupcakeRequest{
				Name:       "To Delete",
				Flavor:     "Test Flavor",
				PriceCents: 1000,
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := newTestService(t)

			if tt.setupCupcake != nil {
				createdCupcake, err := service.CreateCupcake(tt.setupCupcake)
				require.NoError(t, err)
				tt.cupcakeID = createdCupcake.ID
			}

			err := service.DeleteCupcake(tt.cupcakeID)

			if tt.expectedError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}
