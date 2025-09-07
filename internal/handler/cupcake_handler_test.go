package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/julimonteiro/cupcake-store/internal/models"
	"github.com/julimonteiro/cupcake-store/internal/repository"
	"github.com/julimonteiro/cupcake-store/internal/service"
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

func newHandler(t *testing.T) *CupcakeHandler {
	t.Helper()

	db := setupTestDB(t)
	repo := repository.NewCupcakeRepository(db)
	svc := service.NewCupcakeService(repo)
	return NewCupcakeHandler(svc)
}

func newTestRouter(t *testing.T) chi.Router {
	t.Helper()

	handler := newHandler(t)
	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/cupcakes", func(r chi.Router) {
			r.Post("/", handler.CreateCupcake)
			r.Get("/", handler.GetAllCupcakes)
			r.Get("/{id}", handler.GetCupcake)
			r.Put("/{id}", handler.UpdateCupcake)
			r.Delete("/{id}", handler.DeleteCupcake)
		})
	})

	return r
}

func TestCreateCupcake(t *testing.T) {
	tests := []struct {
		name             string
		payload          map[string]interface{}
		expectedStatus   int
		expectedError    string
		validateResponse func(t *testing.T, response models.Cupcake)
	}{
		{
			name: "valid payload returns 201",
			payload: map[string]interface{}{
				"name":        "Chocolate Special",
				"flavor":      "Belgian Chocolate",
				"price_cents": 1500,
			},
			expectedStatus: http.StatusCreated,
			validateResponse: func(t *testing.T, response models.Cupcake) {
				require.Greater(t, response.ID, uint(0))
				require.Equal(t, "Chocolate Special", response.Name)
				require.Equal(t, "Belgian Chocolate", response.Flavor)
				require.Equal(t, 1500, response.PriceCents)
				require.True(t, response.IsAvailable)
			},
		},
		{
			name: "invalid payload - name too short",
			payload: map[string]interface{}{
				"name":        "A",
				"flavor":      "X",
				"price_cents": 1,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "name must have at least 2 characters",
		},
		{
			name: "invalid payload - empty flavor",
			payload: map[string]interface{}{
				"name":        "Valid Name",
				"flavor":      "",
				"price_cents": 1000,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "flavor is required",
		},
		{
			name: "invalid payload - zero price",
			payload: map[string]interface{}{
				"name":        "Valid Name",
				"flavor":      "Valid Flavor",
				"price_cents": 0,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "price must be greater than zero",
		},
		{
			name: "invalid payload - negative price",
			payload: map[string]interface{}{
				"name":        "Valid Name",
				"flavor":      "Valid Flavor",
				"price_cents": -100,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "price must be greater than zero",
		},
		{
			name: "invalid payload - empty name",
			payload: map[string]interface{}{
				"name":        "",
				"flavor":      "Valid Flavor",
				"price_cents": 1000,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "name is required",
		},
		{
			name: "invalid payload - missing required fields",
			payload: map[string]interface{}{
				"name": "Valid Name",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "flavor is required",
		},
		{
			name:           "invalid payload - empty object",
			payload:        map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := newTestRouter(t)

			jsonBody, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			req := httptest.NewRequest("POST", "/api/v1/cupcakes", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedStatus, w.Code)
			require.Equal(t, "application/json", w.Header().Get("Content-Type"))

			if tt.expectedError != "" {
				require.Contains(t, w.Body.String(), tt.expectedError)
			}

			if tt.validateResponse != nil {
				var response models.Cupcake
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				tt.validateResponse(t, response)
			}
		})
	}
}

func TestCreateCupcake_InvalidJSON(t *testing.T) {
	tests := []struct {
		name           string
		payload        string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "malformed JSON",
			payload:        `{"name":"Test", "flavor":"Test", "price_cents":1000, "extra_field": "invalid"`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Error decoding request",
		},
		{
			name:           "invalid JSON syntax",
			payload:        `{"name": "Test", "flavor": "Test", "price_cents": 1000,}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Error decoding request",
		},
		{
			name:           "empty string",
			payload:        "",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Error decoding request",
		},
		{
			name:           "non-JSON string",
			payload:        "invalid json",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Error decoding request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := newTestRouter(t)

			req := httptest.NewRequest("POST", "/api/v1/cupcakes", bytes.NewBufferString(tt.payload))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedStatus, w.Code)
			require.Contains(t, w.Body.String(), tt.expectedError)
		})
	}
}

func TestListCupcakes(t *testing.T) {
	tests := []struct {
		name             string
		setupCupcakes    []map[string]interface{}
		expectedStatus   int
		expectedCount    int
		validateResponse func(t *testing.T, response []models.Cupcake)
	}{
		{
			name:           "empty list returns 200 with empty array",
			setupCupcakes:  []map[string]interface{}{},
			expectedStatus: http.StatusOK,
			expectedCount:  0,
			validateResponse: func(t *testing.T, response []models.Cupcake) {
				require.Len(t, response, 0)
			},
		},
		{
			name: "single cupcake returns 200 with one item",
			setupCupcakes: []map[string]interface{}{
				{
					"name":        "Chocolate",
					"flavor":      "Belgian",
					"price_cents": 1500,
				},
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
			validateResponse: func(t *testing.T, response []models.Cupcake) {
				require.Len(t, response, 1)
				require.Equal(t, "Chocolate", response[0].Name)
				require.Equal(t, "Belgian", response[0].Flavor)
				require.Equal(t, 1500, response[0].PriceCents)
			},
		},
		{
			name: "multiple cupcakes returns 200 with all items",
			setupCupcakes: []map[string]interface{}{
				{
					"name":        "Chocolate",
					"flavor":      "Belgian",
					"price_cents": 1500,
				},
				{
					"name":        "Vanilla",
					"flavor":      "Madagascar",
					"price_cents": 1200,
				},
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
			validateResponse: func(t *testing.T, response []models.Cupcake) {
				require.Len(t, response, 2)
				require.Equal(t, "Chocolate", response[0].Name)
				require.Equal(t, "Vanilla", response[1].Name)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := newTestRouter(t)

			for _, cupcakeData := range tt.setupCupcakes {
				jsonBody, _ := json.Marshal(cupcakeData)
				req := httptest.NewRequest("POST", "/api/v1/cupcakes", bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				require.Equal(t, http.StatusCreated, w.Code)
			}

			req := httptest.NewRequest("GET", "/api/v1/cupcakes", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedStatus, w.Code)
			require.Equal(t, "application/json", w.Header().Get("Content-Type"))

			var response []models.Cupcake
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			require.Len(t, response, tt.expectedCount)

			if tt.validateResponse != nil {
				tt.validateResponse(t, response)
			}
		})
	}
}

func TestGetCupcake(t *testing.T) {
	tests := []struct {
		name             string
		cupcakeID        string
		setupCupcake     map[string]interface{}
		expectedStatus   int
		expectedError    string
		validateResponse func(t *testing.T, response models.Cupcake)
	}{
		{
			name:      "valid ID returns 200",
			cupcakeID: "1",
			setupCupcake: map[string]interface{}{
				"name":        "Chocolate Special",
				"flavor":      "Belgian",
				"price_cents": 1500,
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, response models.Cupcake) {
				require.Equal(t, uint(1), response.ID)
				require.Equal(t, "Chocolate Special", response.Name)
				require.Equal(t, "Belgian", response.Flavor)
				require.Equal(t, 1500, response.PriceCents)
			},
		},
		{
			name:           "non-existent ID returns 404",
			cupcakeID:      "9999",
			expectedStatus: http.StatusNotFound,
			expectedError:  "cupcake not found",
		},
		{
			name:           "invalid ID format returns 400",
			cupcakeID:      "invalid",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid ID",
		},
		{
			name:           "zero ID returns 400",
			cupcakeID:      "0",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := newTestRouter(t)

			if tt.setupCupcake != nil {
				jsonBody, _ := json.Marshal(tt.setupCupcake)
				req := httptest.NewRequest("POST", "/api/v1/cupcakes", bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				require.Equal(t, http.StatusCreated, w.Code)
			}

			req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/cupcakes/%s", tt.cupcakeID), nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedStatus, w.Code)
			require.Equal(t, "application/json", w.Header().Get("Content-Type"))

			if tt.expectedError != "" {
				require.Contains(t, w.Body.String(), tt.expectedError)
			}

			if tt.validateResponse != nil {
				var response models.Cupcake
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				tt.validateResponse(t, response)
			}
		})
	}
}

func TestUpdateCupcake(t *testing.T) {
	tests := []struct {
		name             string
		cupcakeID        string
		updatePayload    map[string]interface{}
		setupCupcake     map[string]interface{}
		expectedStatus   int
		expectedError    string
		validateResponse func(t *testing.T, response models.Cupcake)
	}{
		{
			name:      "valid update returns 200",
			cupcakeID: "1",
			setupCupcake: map[string]interface{}{
				"name":        "Original Name",
				"flavor":      "Original Flavor",
				"price_cents": 1000,
			},
			updatePayload: map[string]interface{}{
				"name":         "Updated Name",
				"flavor":       "Updated Flavor",
				"price_cents":  2000,
				"is_available": false,
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, response models.Cupcake) {
				require.Equal(t, "Updated Name", response.Name)
				require.Equal(t, "Updated Flavor", response.Flavor)
				require.Equal(t, 2000, response.PriceCents)
				require.False(t, response.IsAvailable)
			},
		},
		{
			name:      "partial update returns 200",
			cupcakeID: "1",
			setupCupcake: map[string]interface{}{
				"name":        "Original Name",
				"flavor":      "Original Flavor",
				"price_cents": 1000,
			},
			updatePayload: map[string]interface{}{
				"name": "Updated Name Only",
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, response models.Cupcake) {
				require.Equal(t, "Updated Name Only", response.Name)
				require.Equal(t, "Original Flavor", response.Flavor)
				require.Equal(t, 1000, response.PriceCents)
			},
		},
		{
			name:           "non-existent ID returns 400",
			cupcakeID:      "9999",
			updatePayload:  map[string]interface{}{"name": "Updated"},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "record not found",
		},
		{
			name:           "invalid ID format returns 400",
			cupcakeID:      "invalid",
			updatePayload:  map[string]interface{}{"name": "Updated"},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid ID",
		},
		{
			name:      "invalid update data returns 400",
			cupcakeID: "1",
			setupCupcake: map[string]interface{}{
				"name":        "Original Name",
				"flavor":      "Original Flavor",
				"price_cents": 1000,
			},
			updatePayload: map[string]interface{}{
				"name": "A",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "name must have at least 2 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := newTestRouter(t)

			if tt.setupCupcake != nil {
				jsonBody, _ := json.Marshal(tt.setupCupcake)
				req := httptest.NewRequest("POST", "/api/v1/cupcakes", bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				require.Equal(t, http.StatusCreated, w.Code)
			}

			jsonBody, _ := json.Marshal(tt.updatePayload)
			req := httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/cupcakes/%s", tt.cupcakeID), bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedStatus, w.Code)
			require.Equal(t, "application/json", w.Header().Get("Content-Type"))

			if tt.expectedError != "" {
				require.Contains(t, w.Body.String(), tt.expectedError)
			}

			if tt.validateResponse != nil {
				var response models.Cupcake
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				tt.validateResponse(t, response)
			}
		})
	}
}

func TestUpdateCupcake_InvalidJSON(t *testing.T) {
	tests := []struct {
		name           string
		payload        string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "malformed JSON",
			payload:        `{"name":"Test", "flavor":"Test", "price_cents":1000, "extra_field": "invalid"`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Error decoding request",
		},
		{
			name:           "invalid JSON syntax",
			payload:        `{"name": "Test", "flavor": "Test", "price_cents": 1000,}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Error decoding request",
		},
		{
			name:           "empty string",
			payload:        "",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Error decoding request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := newTestRouter(t)

			req := httptest.NewRequest("PUT", "/api/v1/cupcakes/1", bytes.NewBufferString(tt.payload))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedStatus, w.Code)
			require.Contains(t, w.Body.String(), tt.expectedError)
		})
	}
}

func TestDeleteCupcake(t *testing.T) {
	tests := []struct {
		name           string
		cupcakeID      string
		setupCupcake   map[string]interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name:      "valid ID returns 204",
			cupcakeID: "1",
			setupCupcake: map[string]interface{}{
				"name":        "To Delete",
				"flavor":      "Test Flavor",
				"price_cents": 1000,
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "non-existent ID returns 400",
			cupcakeID:      "9999",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "record not found",
		},
		{
			name:           "invalid ID format returns 400",
			cupcakeID:      "invalid",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid ID",
		},
		{
			name:           "zero ID returns 400",
			cupcakeID:      "0",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := newTestRouter(t)

			if tt.setupCupcake != nil {
				jsonBody, _ := json.Marshal(tt.setupCupcake)
				req := httptest.NewRequest("POST", "/api/v1/cupcakes", bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				require.Equal(t, http.StatusCreated, w.Code)
			}

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/cupcakes/%s", tt.cupcakeID), nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				require.Contains(t, w.Body.String(), tt.expectedError)
			}
		})
	}
}

func TestHealthCheck(t *testing.T) {
	tests := []struct {
		name             string
		expectedStatus   int
		validateResponse func(t *testing.T, response map[string]interface{})
	}{
		{
			name:           "health check returns 200",
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, response map[string]interface{}) {
				require.Equal(t, "ok", response["status"])
				require.Equal(t, "Cupcake Store API is running!", response["message"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := newHandler(t)
			r := chi.NewRouter()
			r.Get("/health", handler.HealthCheck)

			req := httptest.NewRequest("GET", "/health", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			require.Equal(t, tt.expectedStatus, w.Code)
			require.Equal(t, "application/json", w.Header().Get("Content-Type"))

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			if tt.validateResponse != nil {
				tt.validateResponse(t, response)
			}
		})
	}
}
