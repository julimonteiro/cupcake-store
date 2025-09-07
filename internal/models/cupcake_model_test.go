package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCupcake_TableName(t *testing.T) {
	tests := []struct {
		name     string
		cupcake  Cupcake
		expected string
	}{
		{
			name:     "returns correct table name",
			cupcake:  Cupcake{},
			expected: "cupcakes",
		},
		{
			name: "returns correct table name for populated cupcake",
			cupcake: Cupcake{
				ID:          1,
				Name:        "Test Cupcake",
				Flavor:      "Test Flavor",
				PriceCents:  1000,
				IsAvailable: true,
			},
			expected: "cupcakes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cupcake.TableName()
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestCreateCupcakeRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request CreateCupcakeRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			request: CreateCupcakeRequest{
				Name:       "Chocolate",
				Flavor:     "Dark",
				PriceCents: 100,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			request: CreateCupcakeRequest{
				Name:       "",
				Flavor:     "Dark",
				PriceCents: 100,
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "name too short",
			request: CreateCupcakeRequest{
				Name:       "A",
				Flavor:     "Dark",
				PriceCents: 100,
			},
			wantErr: true,
			errMsg:  "name must have at least 2 characters",
		},
		{
			name: "name with only spaces",
			request: CreateCupcakeRequest{
				Name:       "   ",
				Flavor:     "Dark",
				PriceCents: 100,
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "empty flavor",
			request: CreateCupcakeRequest{
				Name:       "Chocolate",
				Flavor:     "",
				PriceCents: 100,
			},
			wantErr: true,
			errMsg:  "flavor is required",
		},
		{
			name: "flavor with only spaces",
			request: CreateCupcakeRequest{
				Name:       "Chocolate",
				Flavor:     "   ",
				PriceCents: 100,
			},
			wantErr: true,
			errMsg:  "flavor is required",
		},
		{
			name: "zero price",
			request: CreateCupcakeRequest{
				Name:       "Chocolate",
				Flavor:     "Dark",
				PriceCents: 0,
			},
			wantErr: true,
			errMsg:  "price must be greater than zero",
		},
		{
			name: "negative price",
			request: CreateCupcakeRequest{
				Name:       "Chocolate",
				Flavor:     "Dark",
				PriceCents: -10,
			},
			wantErr: true,
			errMsg:  "price must be greater than zero",
		},
		{
			name: "valid request with trimmed fields",
			request: CreateCupcakeRequest{
				Name:       "  Chocolate  ",
				Flavor:     "  Dark  ",
				PriceCents: 100,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NotNil(t, tt.request)
			require.Equal(t, tt.request.Name, tt.request.Name)
			require.Equal(t, tt.request.Flavor, tt.request.Flavor)
			require.Equal(t, tt.request.PriceCents, tt.request.PriceCents)
		})
	}
}

func TestUpdateCupcakeRequest_OptionalFields(t *testing.T) {
	tests := []struct {
		name              string
		request           UpdateCupcakeRequest
		expectedName      *string
		expectedFlavor    *string
		expectedPrice     *int
		expectedAvailable *bool
	}{
		{
			name: "all fields set",
			request: UpdateCupcakeRequest{
				Name:        stringPtr("New Name"),
				Flavor:      stringPtr("New Flavor"),
				PriceCents:  intPtr(200),
				IsAvailable: boolPtr(false),
			},
			expectedName:      stringPtr("New Name"),
			expectedFlavor:    stringPtr("New Flavor"),
			expectedPrice:     intPtr(200),
			expectedAvailable: boolPtr(false),
		},
		{
			name: "partial fields set",
			request: UpdateCupcakeRequest{
				Name:   stringPtr("New Name"),
				Flavor: stringPtr("New Flavor"),
			},
			expectedName:      stringPtr("New Name"),
			expectedFlavor:    stringPtr("New Flavor"),
			expectedPrice:     nil,
			expectedAvailable: nil,
		},
		{
			name: "only name set",
			request: UpdateCupcakeRequest{
				Name: stringPtr("New Name"),
			},
			expectedName:      stringPtr("New Name"),
			expectedFlavor:    nil,
			expectedPrice:     nil,
			expectedAvailable: nil,
		},
		{
			name: "only price set",
			request: UpdateCupcakeRequest{
				PriceCents: intPtr(300),
			},
			expectedName:      nil,
			expectedFlavor:    nil,
			expectedPrice:     intPtr(300),
			expectedAvailable: nil,
		},
		{
			name: "only availability set",
			request: UpdateCupcakeRequest{
				IsAvailable: boolPtr(true),
			},
			expectedName:      nil,
			expectedFlavor:    nil,
			expectedPrice:     nil,
			expectedAvailable: boolPtr(true),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expectedName, tt.request.Name)
			require.Equal(t, tt.expectedFlavor, tt.request.Flavor)
			require.Equal(t, tt.expectedPrice, tt.request.PriceCents)
		})
	}
}

func TestUpdateCupcakeRequest_NilFields(t *testing.T) {
	tests := []struct {
		name    string
		request UpdateCupcakeRequest
	}{
		{
			name:    "empty request",
			request: UpdateCupcakeRequest{},
		},
		{
			name: "request with nil pointers",
			request: UpdateCupcakeRequest{
				Name:        nil,
				Flavor:      nil,
				PriceCents:  nil,
				IsAvailable: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Nil(t, tt.request.Name)
			require.Nil(t, tt.request.Flavor)
			require.Nil(t, tt.request.PriceCents)
			require.Nil(t, tt.request.IsAvailable)
		})
	}
}

func TestCupcake_Fields(t *testing.T) {
	tests := []struct {
		name              string
		cupcake           Cupcake
		expectedID        uint
		expectedName      string
		expectedFlavor    string
		expectedPrice     int
		expectedAvailable bool
		expectedCreatedAt time.Time
		expectedUpdatedAt time.Time
	}{
		{
			name: "fully populated cupcake",
			cupcake: Cupcake{
				ID:          1,
				Name:        "Test Cupcake",
				Flavor:      "Test Flavor",
				PriceCents:  1000,
				IsAvailable: true,
				CreatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedID:        1,
			expectedName:      "Test Cupcake",
			expectedFlavor:    "Test Flavor",
			expectedPrice:     1000,
			expectedAvailable: true,
			expectedCreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			expectedUpdatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "cupcake with zero values",
			cupcake: Cupcake{
				ID:          0,
				Name:        "",
				Flavor:      "",
				PriceCents:  0,
				IsAvailable: false,
				CreatedAt:   time.Time{},
				UpdatedAt:   time.Time{},
			},
			expectedID:        0,
			expectedName:      "",
			expectedFlavor:    "",
			expectedPrice:     0,
			expectedAvailable: false,
			expectedCreatedAt: time.Time{},
			expectedUpdatedAt: time.Time{},
		},
		{
			name: "cupcake with high values",
			cupcake: Cupcake{
				ID:          999999,
				Name:        "Very Long Cupcake Name That Exceeds Normal Length",
				Flavor:      "Very Long Flavor Name That Exceeds Normal Length",
				PriceCents:  999999,
				IsAvailable: true,
				CreatedAt:   time.Date(2023, 12, 31, 23, 59, 59, 999999999, time.UTC),
				UpdatedAt:   time.Date(2023, 12, 31, 23, 59, 59, 999999999, time.UTC),
			},
			expectedID:        999999,
			expectedName:      "Very Long Cupcake Name That Exceeds Normal Length",
			expectedFlavor:    "Very Long Flavor Name That Exceeds Normal Length",
			expectedPrice:     999999,
			expectedAvailable: true,
			expectedCreatedAt: time.Date(2023, 12, 31, 23, 59, 59, 999999999, time.UTC),
			expectedUpdatedAt: time.Date(2023, 12, 31, 23, 59, 59, 999999999, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expectedID, tt.cupcake.ID)
			require.Equal(t, tt.expectedName, tt.cupcake.Name)
			require.Equal(t, tt.expectedFlavor, tt.cupcake.Flavor)
			require.Equal(t, tt.expectedPrice, tt.cupcake.PriceCents)
			require.Equal(t, tt.expectedAvailable, tt.cupcake.IsAvailable)
			require.Equal(t, tt.expectedCreatedAt, tt.cupcake.CreatedAt)
			require.Equal(t, tt.expectedUpdatedAt, tt.cupcake.UpdatedAt)
		})
	}
}

func TestCreateCupcakeRequest_Fields(t *testing.T) {
	tests := []struct {
		name              string
		request           CreateCupcakeRequest
		expectedName      string
		expectedFlavor    string
		expectedPrice     int
		expectedAvailable bool
	}{
		{
			name: "fully populated request",
			request: CreateCupcakeRequest{
				Name:       "Test Cupcake",
				Flavor:     "Test Flavor",
				PriceCents: 1000,
			},
			expectedName:      "Test Cupcake",
			expectedFlavor:    "Test Flavor",
			expectedPrice:     1000,
			expectedAvailable: true,
		},
		{
			name: "request with default values",
			request: CreateCupcakeRequest{
				Name:       "Test Cupcake",
				Flavor:     "Test Flavor",
				PriceCents: 1000,
			},
			expectedName:      "Test Cupcake",
			expectedFlavor:    "Test Flavor",
			expectedPrice:     1000,
			expectedAvailable: false,
		},
		{
			name: "request with zero values",
			request: CreateCupcakeRequest{
				Name:       "",
				Flavor:     "",
				PriceCents: 0,
			},
			expectedName:      "",
			expectedFlavor:    "",
			expectedPrice:     0,
			expectedAvailable: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expectedName, tt.request.Name)
			require.Equal(t, tt.expectedFlavor, tt.request.Flavor)
			require.Equal(t, tt.expectedPrice, tt.request.PriceCents)
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
