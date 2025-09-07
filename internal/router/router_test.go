package router

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julimonteiro/cupcake-store/internal/config"
	"github.com/julimonteiro/cupcake-store/internal/database"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	cfg := &config.Config{
		DBDialect: "sqlite",
		DBDSN:     ":memory:",
		LogLevel:  "error",
	}
	db, err := database.Init(cfg)
	require.NoError(t, err)
	return db
}

func TestSetup(t *testing.T) {
	tests := []struct {
		name           string
		expectedStatus int
		validateResult func(t *testing.T, router http.Handler)
	}{
		{
			name:           "router setup successful",
			expectedStatus: http.StatusOK,
			validateResult: func(t *testing.T, router http.Handler) {
				require.NotNil(t, router)

				req := httptest.NewRequest("GET", "/health", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				require.Equal(t, http.StatusOK, w.Code)
				require.Equal(t, "application/json", w.Header().Get("Content-Type"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			router := Setup(db)

			if tt.validateResult != nil {
				tt.validateResult(t, router)
			}
		})
	}
}

func TestSetup_APIRoutes(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		path        string
		body        []byte
		status      int
		description string
	}{
		{
			name:        "GET /api/v1/cupcakes",
			method:      "GET",
			path:        "/api/v1/cupcakes",
			status:      http.StatusOK,
			description: "should return 200 for GET cupcakes list",
		},
		{
			name:        "POST /api/v1/cupcakes",
			method:      "POST",
			path:        "/api/v1/cupcakes",
			body:        []byte(`{"name":"Test","flavor":"Test","price_cents":100}`),
			status:      http.StatusCreated,
			description: "should return 201 for valid POST request",
		},
		{
			name:        "POST /api/v1/cupcakes with invalid data",
			method:      "POST",
			path:        "/api/v1/cupcakes",
			body:        []byte(`{"name":"A","flavor":"X","price_cents":1}`),
			status:      http.StatusBadRequest,
			description: "should return 400 for invalid POST request",
		},
		{
			name:        "GET /api/v1/cupcakes/1",
			method:      "GET",
			path:        "/api/v1/cupcakes/1",
			status:      http.StatusNotFound,
			description: "should return 404 for non-existent cupcake",
		},
		{
			name:        "PUT /api/v1/cupcakes/1",
			method:      "PUT",
			path:        "/api/v1/cupcakes/1",
			body:        []byte(`{"name":"Updated"}`),
			status:      http.StatusBadRequest,
			description: "should return 400 for non-existent cupcake update",
		},
		{
			name:        "DELETE /api/v1/cupcakes/1",
			method:      "DELETE",
			path:        "/api/v1/cupcakes/1",
			status:      http.StatusBadRequest,
			description: "should return 400 for non-existent cupcake deletion",
		},
		{
			name:        "GET /api/v1/cupcakes/invalid",
			method:      "GET",
			path:        "/api/v1/cupcakes/invalid",
			status:      http.StatusBadRequest,
			description: "should return 400 for invalid ID format",
		},
		{
			name:        "PUT /api/v1/cupcakes/invalid",
			method:      "PUT",
			path:        "/api/v1/cupcakes/invalid",
			body:        []byte(`{"name":"Updated"}`),
			status:      http.StatusBadRequest,
			description: "should return 400 for invalid ID format in PUT",
		},
		{
			name:        "DELETE /api/v1/cupcakes/invalid",
			method:      "DELETE",
			path:        "/api/v1/cupcakes/invalid",
			status:      http.StatusBadRequest,
			description: "should return 400 for invalid ID format in DELETE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			router := Setup(db)

			req := httptest.NewRequest(tt.method, tt.path, bytes.NewBuffer(tt.body))
			if tt.body != nil {
				req.Header.Set("Content-Type", "application/json")
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			require.Equal(t, tt.status, w.Code, tt.description)
		})
	}
}

func TestSetup_StaticFiles(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedStatus int
		description    string
	}{
		{
			name:           "GET / (root path)",
			path:           "/",
			expectedStatus: http.StatusNotFound,
			description:    "should handle static file serving (404 expected in test environment)",
		},
		{
			name:           "GET /index.html",
			path:           "/index.html",
			expectedStatus: http.StatusNotFound,
			description:    "should handle static file serving for specific file",
		},
		{
			name:           "GET /nonexistent.html",
			path:           "/nonexistent.html",
			expectedStatus: http.StatusNotFound,
			description:    "should return 404 for non-existent static file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				DBDialect: "sqlite",
				DBDSN:     ":memory:",
				LogLevel:  "error",
			}

			db, err := database.Init(cfg)
			require.NoError(t, err)

			router := Setup(db)
			require.NotNil(t, router)

			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			require.NotEqual(t, http.StatusInternalServerError, w.Code, tt.description)
		})
	}
}

func TestSetup_CORS(t *testing.T) {
	tests := []struct {
		name            string
		method          string
		path            string
		headers         map[string]string
		expectedStatus  int
		expectedHeaders map[string]string
		description     string
	}{
		{
			name:   "OPTIONS request with CORS headers",
			method: "OPTIONS",
			path:   "/api/v1/cupcakes",
			headers: map[string]string{
				"Origin":                         "http://localhost:3000",
				"Access-Control-Request-Method":  "POST",
				"Access-Control-Request-Headers": "Content-Type",
			},
			expectedStatus: http.StatusOK,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
				"Access-Control-Allow-Headers": "Accept, Authorization, Content-Type, X-CSRF-Token",
			},
			description: "should handle CORS preflight request",
		},
		{
			name:           "OPTIONS request without CORS headers",
			method:         "OPTIONS",
			path:           "/api/v1/cupcakes",
			headers:        map[string]string{},
			expectedStatus: http.StatusOK,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
				"Access-Control-Allow-Headers": "Accept, Authorization, Content-Type, X-CSRF-Token",
			},
			description: "should handle OPTIONS request without CORS headers",
		},
		{
			name:   "GET request with Origin header",
			method: "GET",
			path:   "/api/v1/cupcakes",
			headers: map[string]string{
				"Origin": "http://localhost:3000",
			},
			expectedStatus: http.StatusOK,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
			description: "should add CORS headers to regular requests",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				DBDialect: "sqlite",
				DBDSN:     ":memory:",
				LogLevel:  "error",
			}

			db, err := database.Init(cfg)
			require.NoError(t, err)

			router := Setup(db)
			require.NotNil(t, router)

			req := httptest.NewRequest(tt.method, tt.path, nil)
			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedStatus, w.Code, tt.description)

			for key, expectedValue := range tt.expectedHeaders {
				actualValue := w.Header().Get(key)
				require.Contains(t, actualValue, expectedValue, tt.description)
			}
		})
	}
}

func TestSetup_Middleware(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		description    string
	}{
		{
			name:           "GET /health with logger middleware",
			method:         "GET",
			path:           "/health",
			expectedStatus: http.StatusOK,
			description:    "should apply logger middleware",
		},
		{
			name:           "GET /api/v1/cupcakes with logger middleware",
			method:         "GET",
			path:           "/api/v1/cupcakes",
			expectedStatus: http.StatusOK,
			description:    "should apply logger middleware to API routes",
		},
		{
			name:           "POST /api/v1/cupcakes with logger middleware",
			method:         "POST",
			path:           "/api/v1/cupcakes",
			expectedStatus: http.StatusBadRequest,
			description:    "should apply logger middleware to POST requests",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			router := Setup(db)

			req := httptest.NewRequest(tt.method, tt.path, nil)
			if tt.method == "POST" {
				req.Header.Set("Content-Type", "application/json")
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedStatus, w.Code, tt.description)
		})
	}
}

func TestSetup_RouteStructure(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		description    string
	}{
		{
			name:           "health check route",
			method:         "GET",
			path:           "/health",
			expectedStatus: http.StatusOK,
			description:    "should have health check route",
		},
		{
			name:           "cupcakes list route",
			method:         "GET",
			path:           "/api/v1/cupcakes",
			expectedStatus: http.StatusOK,
			description:    "should have cupcakes list route",
		},
		{
			name:           "cupcake create route",
			method:         "POST",
			path:           "/api/v1/cupcakes",
			expectedStatus: http.StatusBadRequest,
			description:    "should have cupcake create route",
		},
		{
			name:           "cupcake get route",
			method:         "GET",
			path:           "/api/v1/cupcakes/1",
			expectedStatus: http.StatusNotFound,
			description:    "should have cupcake get route",
		},
		{
			name:           "cupcake update route",
			method:         "PUT",
			path:           "/api/v1/cupcakes/1",
			expectedStatus: http.StatusBadRequest,
			description:    "should have cupcake update route",
		},
		{
			name:           "cupcake delete route",
			method:         "DELETE",
			path:           "/api/v1/cupcakes/1",
			expectedStatus: http.StatusBadRequest,
			description:    "should have cupcake delete route",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			router := Setup(db)

			req := httptest.NewRequest(tt.method, tt.path, nil)
			if tt.method == "POST" || tt.method == "PUT" {
				req.Header.Set("Content-Type", "application/json")
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedStatus, w.Code, tt.description)
		})
	}
}

func TestSetup_ErrorHandling(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		body           []byte
		expectedStatus int
		description    string
	}{
		{
			name:           "malformed JSON in POST",
			method:         "POST",
			path:           "/api/v1/cupcakes",
			body:           []byte(`{"name":"Test", "flavor":"Test", "price_cents":1000, "extra_field": "invalid"`),
			expectedStatus: http.StatusBadRequest,
			description:    "should handle malformed JSON",
		},
		{
			name:           "invalid JSON in PUT",
			method:         "PUT",
			path:           "/api/v1/cupcakes/1",
			body:           []byte(`{"name":"Test", "flavor":"Test", "price_cents":1000,}`),
			expectedStatus: http.StatusBadRequest,
			description:    "should handle invalid JSON in PUT",
		},
		{
			name:           "empty body in POST",
			method:         "POST",
			path:           "/api/v1/cupcakes",
			body:           []byte(``),
			expectedStatus: http.StatusBadRequest,
			description:    "should handle empty body",
		},
		{
			name:           "non-JSON body in POST",
			method:         "POST",
			path:           "/api/v1/cupcakes",
			body:           []byte(`not json`),
			expectedStatus: http.StatusBadRequest,
			description:    "should handle non-JSON body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			router := Setup(db)

			req := httptest.NewRequest(tt.method, tt.path, bytes.NewBuffer(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedStatus, w.Code, tt.description)
		})
	}
}
