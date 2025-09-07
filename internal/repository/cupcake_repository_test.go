package repository

import (
	"testing"

	"github.com/julimonteiro/cupcake-store/internal/models"
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

func TestNewCupcakeRepository(t *testing.T) {
	tests := []struct {
		name string
		db   *gorm.DB
	}{
		{
			name: "creates repository with valid DB",
			db:   setupTestDB(t),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewCupcakeRepository(tt.db)
			require.NotNil(t, repo)
			require.Equal(t, tt.db, repo.db)
		})
	}
}

func TestCupcakeRepository_Create(t *testing.T) {
	tests := []struct {
		name           string
		cupcake        *models.Cupcake
		validateResult func(t *testing.T, cupcake *models.Cupcake, db *gorm.DB)
	}{
		{
			name: "creates cupcake successfully",
			cupcake: &models.Cupcake{
				Name:       "Test Cupcake",
				Flavor:     "Vanilla",
				PriceCents: 1000,
			},
			validateResult: func(t *testing.T, cupcake *models.Cupcake, db *gorm.DB) {
				require.True(t, cupcake.ID > 0)
				var createdCupcake models.Cupcake
				db.First(&createdCupcake, cupcake.ID)
				require.Equal(t, cupcake.Name, createdCupcake.Name)
			},
		},
		{
			name: "creates cupcake with all fields",
			cupcake: &models.Cupcake{
				Name:        "Chocolate Special",
				Flavor:      "Belgian Chocolate",
				PriceCents:  1500,
				IsAvailable: false,
			},
			validateResult: func(t *testing.T, cupcake *models.Cupcake, db *gorm.DB) {
				require.True(t, cupcake.ID > 0)
				require.Equal(t, "Chocolate Special", cupcake.Name)
				require.Equal(t, "Belgian Chocolate", cupcake.Flavor)
				require.Equal(t, 1500, cupcake.PriceCents)
				require.False(t, cupcake.IsAvailable)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := NewCupcakeRepository(db)

			err := repo.Create(tt.cupcake)
			require.NoError(t, err)

			if tt.validateResult != nil {
				tt.validateResult(t, tt.cupcake, db)
			}
		})
	}
}

func TestCupcakeRepository_FindByID(t *testing.T) {
	tests := []struct {
		name           string
		cupcakeID      uint
		setupCupcake   *models.Cupcake
		expectedError  string
		validateResult func(t *testing.T, cupcake *models.Cupcake)
	}{
		{
			name:      "finds existing cupcake",
			cupcakeID: 1,
			setupCupcake: &models.Cupcake{
				Name:       "Test Cupcake",
				Flavor:     "Vanilla",
				PriceCents: 1000,
			},
			validateResult: func(t *testing.T, cupcake *models.Cupcake) {
				require.NotNil(t, cupcake)
				require.Equal(t, uint(1), cupcake.ID)
				require.Equal(t, "Test Cupcake", cupcake.Name)
			},
		},
		{
			name:          "returns error for non-existent cupcake",
			cupcakeID:     999,
			expectedError: "record not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := NewCupcakeRepository(db)

			if tt.setupCupcake != nil {
				err := repo.Create(tt.setupCupcake)
				require.NoError(t, err)
				tt.cupcakeID = tt.setupCupcake.ID
			}

			foundCupcake, err := repo.FindByID(tt.cupcakeID)

			if tt.expectedError != "" {
				require.Error(t, err)
				require.Nil(t, foundCupcake)
				require.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
				require.NotNil(t, foundCupcake)
				if tt.validateResult != nil {
					tt.validateResult(t, foundCupcake)
				}
			}
		})
	}
}

func TestCupcakeRepository_FindAll(t *testing.T) {
	tests := []struct {
		name           string
		setupCupcakes  []*models.Cupcake
		expectedCount  int
		validateResult func(t *testing.T, cupcakes []models.Cupcake)
	}{
		{
			name:          "returns empty list when no cupcakes",
			setupCupcakes: []*models.Cupcake{},
			expectedCount: 0,
			validateResult: func(t *testing.T, cupcakes []models.Cupcake) {
				require.Len(t, cupcakes, 0)
			},
		},
		{
			name: "returns single cupcake",
			setupCupcakes: []*models.Cupcake{
				{
					Name:       "C1",
					Flavor:     "F1",
					PriceCents: 100,
				},
			},
			expectedCount: 1,
			validateResult: func(t *testing.T, cupcakes []models.Cupcake) {
				require.Len(t, cupcakes, 1)
				require.Equal(t, "C1", cupcakes[0].Name)
			},
		},
		{
			name: "returns multiple cupcakes",
			setupCupcakes: []*models.Cupcake{
				{
					Name:       "C1",
					Flavor:     "F1",
					PriceCents: 100,
				},
				{
					Name:       "C2",
					Flavor:     "F2",
					PriceCents: 200,
				},
			},
			expectedCount: 2,
			validateResult: func(t *testing.T, cupcakes []models.Cupcake) {
				require.Len(t, cupcakes, 2)
				require.Equal(t, "C1", cupcakes[0].Name)
				require.Equal(t, "C2", cupcakes[1].Name)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := NewCupcakeRepository(db)

			for _, cupcake := range tt.setupCupcakes {
				err := repo.Create(cupcake)
				require.NoError(t, err)
			}

			cupcakes, err := repo.FindAll()
			require.NoError(t, err)
			require.Len(t, cupcakes, tt.expectedCount)

			if tt.validateResult != nil {
				tt.validateResult(t, cupcakes)
			}
		})
	}
}

func TestCupcakeRepository_Update(t *testing.T) {
	tests := []struct {
		name            string
		originalCupcake *models.Cupcake
		updatedCupcake  *models.Cupcake
		validateResult  func(t *testing.T, cupcake *models.Cupcake, db *gorm.DB)
	}{
		{
			name: "updates cupcake successfully",
			originalCupcake: &models.Cupcake{
				Name:       "Old Name",
				Flavor:     "Old Flavor",
				PriceCents: 100,
			},
			updatedCupcake: &models.Cupcake{
				Name:       "New Name",
				Flavor:     "New Flavor",
				PriceCents: 200,
			},
			validateResult: func(t *testing.T, cupcake *models.Cupcake, db *gorm.DB) {
				var updatedCupcake models.Cupcake
				db.First(&updatedCupcake, cupcake.ID)
				require.Equal(t, "New Name", updatedCupcake.Name)
				require.Equal(t, "New Flavor", updatedCupcake.Flavor)
				require.Equal(t, 200, updatedCupcake.PriceCents)
			},
		},
		{
			name: "updates partial fields",
			originalCupcake: &models.Cupcake{
				Name:       "Original Name",
				Flavor:     "Original Flavor",
				PriceCents: 100,
			},
			updatedCupcake: &models.Cupcake{
				Name:       "Updated Name Only",
				Flavor:     "Original Flavor",
				PriceCents: 100,
			},
			validateResult: func(t *testing.T, cupcake *models.Cupcake, db *gorm.DB) {
				var updatedCupcake models.Cupcake
				db.First(&updatedCupcake, cupcake.ID)
				require.Equal(t, "Updated Name Only", updatedCupcake.Name)
				require.Equal(t, "Original Flavor", updatedCupcake.Flavor)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := NewCupcakeRepository(db)

			err := repo.Create(tt.originalCupcake)
			require.NoError(t, err)

			tt.updatedCupcake.ID = tt.originalCupcake.ID
			err = repo.Update(tt.updatedCupcake)
			require.NoError(t, err)

			if tt.validateResult != nil {
				tt.validateResult(t, tt.updatedCupcake, db)
			}
		})
	}
}

func TestCupcakeRepository_Delete(t *testing.T) {
	tests := []struct {
		name           string
		cupcakeID      uint
		setupCupcake   *models.Cupcake
		expectedError  string
		validateResult func(t *testing.T, db *gorm.DB, cupcakeID uint)
	}{
		{
			name:      "deletes cupcake successfully",
			cupcakeID: 1,
			setupCupcake: &models.Cupcake{
				Name:       "To Delete",
				Flavor:     "Test",
				PriceCents: 100,
			},
			validateResult: func(t *testing.T, db *gorm.DB, cupcakeID uint) {
				var deletedCupcake models.Cupcake
				result := db.First(&deletedCupcake, cupcakeID)
				require.Error(t, result.Error)
				require.Contains(t, result.Error.Error(), "record not found")
			},
		},
		{
			name:          "returns error for non-existent cupcake",
			cupcakeID:     999,
			expectedError: "record not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := NewCupcakeRepository(db)

			if tt.setupCupcake != nil {
				err := repo.Create(tt.setupCupcake)
				require.NoError(t, err)
				tt.cupcakeID = tt.setupCupcake.ID
			}

			err := repo.Delete(tt.cupcakeID)

			if tt.expectedError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
				if tt.validateResult != nil {
					tt.validateResult(t, db, tt.cupcakeID)
				}
			}
		})
	}
}

func TestCupcakeRepository_Exists(t *testing.T) {
	tests := []struct {
		name           string
		cupcakeID      uint
		setupCupcake   *models.Cupcake
		expectedExists bool
		expectedError  string
	}{
		{
			name:      "returns true for existing cupcake",
			cupcakeID: 1,
			setupCupcake: &models.Cupcake{
				Name:       "Exists",
				Flavor:     "Test",
				PriceCents: 100,
			},
			expectedExists: true,
		},
		{
			name:           "returns false for non-existent cupcake",
			cupcakeID:      999,
			expectedExists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := NewCupcakeRepository(db)

			if tt.setupCupcake != nil {
				err := repo.Create(tt.setupCupcake)
				require.NoError(t, err)
				tt.cupcakeID = tt.setupCupcake.ID
			}

			exists, err := repo.Exists(tt.cupcakeID)

			if tt.expectedError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedExists, exists)
			}
		})
	}
}
