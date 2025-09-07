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
			r.Get("/", handler.GetAllCupcakes)
			r.Post("/", handler.CreateCupcake)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", handler.GetCupcake)
				r.Put("/", handler.UpdateCupcake)
				r.Delete("/", handler.DeleteCupcake)
			})
		})
	})

	return r
}

func TestCreateCupcake_Returns201AndBody(t *testing.T) {
	router := newTestRouter(t)

	reqBody := map[string]interface{}{
		"name":        "Chocolate Special",
		"flavor":      "Belgian Chocolate",
		"price_cents": 1500,
	}

	jsonBody, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/api/v1/cupcakes", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response models.Cupcake
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	require.Greater(t, response.ID, uint(0))
	require.Equal(t, "Chocolate Special", response.Name)
	require.Equal(t, "Belgian Chocolate", response.Flavor)
	require.Equal(t, 1500, response.PriceCents)
	require.True(t, response.IsAvailable)
}

func TestCreateCupcake_InvalidPayload_Returns400(t *testing.T) {
	router := newTestRouter(t)

	reqBody := map[string]interface{}{
		"name":        "A",
		"flavor":      "X",
		"price_cents": 1,
	}

	jsonBody, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/api/v1/cupcakes", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), "name must have at least 2 characters")
}

func TestCreateCupcake_InvalidJSON_Returns400(t *testing.T) {
	router := newTestRouter(t)

	req := httptest.NewRequest("POST", "/api/v1/cupcakes", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), "Error decoding request")
}

func TestListCupcakes_Returns200AndArray(t *testing.T) {
	router := newTestRouter(t)

	reqBody1 := map[string]interface{}{
		"name":        "Chocolate",
		"flavor":      "Belgian",
		"price_cents": 1500,
	}
	jsonBody1, _ := json.Marshal(reqBody1)
	req1 := httptest.NewRequest("POST", "/api/v1/cupcakes", bytes.NewBuffer(jsonBody1))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)
	require.Equal(t, http.StatusCreated, w1.Code)

	reqBody2 := map[string]interface{}{
		"name":        "Vanilla",
		"flavor":      "Madagascar",
		"price_cents": 1200,
	}
	jsonBody2, _ := json.Marshal(reqBody2)
	req2 := httptest.NewRequest("POST", "/api/v1/cupcakes", bytes.NewBuffer(jsonBody2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	require.Equal(t, http.StatusCreated, w2.Code)

	req := httptest.NewRequest("GET", "/api/v1/cupcakes", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var cupcakes []models.Cupcake
	err := json.Unmarshal(w.Body.Bytes(), &cupcakes)
	require.NoError(t, err)
	require.Len(t, cupcakes, 2)

	require.Equal(t, "Chocolate", cupcakes[0].Name)
	require.Equal(t, "Vanilla", cupcakes[1].Name)
}

func TestGetCupcake_NotFound_Returns404(t *testing.T) {
	router := newTestRouter(t)

	req := httptest.NewRequest("GET", "/api/v1/cupcakes/9999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
	require.Contains(t, w.Body.String(), "Cupcake not found")
}

func TestGetCupcake_InvalidID_Returns400(t *testing.T) {
	router := newTestRouter(t)

	req := httptest.NewRequest("GET", "/api/v1/cupcakes/invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), "Invalid ID")
}

func TestGetCupcake_Success_Returns200(t *testing.T) {
	router := newTestRouter(t)

	reqBody := map[string]interface{}{
		"name":        "Test Cupcake",
		"flavor":      "Test Flavor",
		"price_cents": 1000,
	}
	jsonBody, _ := json.Marshal(reqBody)
	createReq := httptest.NewRequest("POST", "/api/v1/cupcakes", bytes.NewBuffer(jsonBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)
	require.Equal(t, http.StatusCreated, createW.Code)

	var created models.Cupcake
	err := json.Unmarshal(createW.Body.Bytes(), &created)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/api/v1/cupcakes/"+fmt.Sprintf("%d", created.ID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response models.Cupcake
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	require.Equal(t, created.ID, response.ID)
	require.Equal(t, "Test Cupcake", response.Name)
	require.Equal(t, "Test Flavor", response.Flavor)
	require.Equal(t, 1000, response.PriceCents)
}

func TestUpdateCupcake_Returns200AndUpdatedBody(t *testing.T) {
	router := newTestRouter(t)

	reqBody := map[string]interface{}{
		"name":        "Original",
		"flavor":      "Original",
		"price_cents": 1000,
	}
	jsonBody, _ := json.Marshal(reqBody)
	createReq := httptest.NewRequest("POST", "/api/v1/cupcakes", bytes.NewBuffer(jsonBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)
	require.Equal(t, http.StatusCreated, createW.Code)

	var created models.Cupcake
	err := json.Unmarshal(createW.Body.Bytes(), &created)
	require.NoError(t, err)

	updateBody := map[string]interface{}{
		"name":        "Updated Name",
		"flavor":      "Updated Flavor",
		"price_cents": 2000,
	}
	updateJsonBody, _ := json.Marshal(updateBody)
	updateReq := httptest.NewRequest("PUT", "/api/v1/cupcakes/"+fmt.Sprintf("%d", created.ID), bytes.NewBuffer(updateJsonBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateW := httptest.NewRecorder()
	router.ServeHTTP(updateW, updateReq)

	require.Equal(t, http.StatusOK, updateW.Code)
	require.Equal(t, "application/json", updateW.Header().Get("Content-Type"))

	var response models.Cupcake
	err = json.Unmarshal(updateW.Body.Bytes(), &response)
	require.NoError(t, err)

	require.Equal(t, created.ID, response.ID)
	require.Equal(t, "Updated Name", response.Name)
	require.Equal(t, "Updated Flavor", response.Flavor)
	require.Equal(t, 2000, response.PriceCents)
	require.True(t, response.IsAvailable)
}

func TestUpdateCupcake_NotFound_Returns400(t *testing.T) {
	router := newTestRouter(t)

	updateBody := map[string]interface{}{
		"name": "New Name",
	}
	updateJsonBody, _ := json.Marshal(updateBody)
	updateReq := httptest.NewRequest("PUT", "/api/v1/cupcakes/999", bytes.NewBuffer(updateJsonBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateW := httptest.NewRecorder()
	router.ServeHTTP(updateW, updateReq)

	require.Equal(t, http.StatusBadRequest, updateW.Code)
	require.Contains(t, updateW.Body.String(), "cupcake not found")
}

func TestUpdateCupcake_InvalidID_Returns400(t *testing.T) {
	router := newTestRouter(t)

	updateBody := map[string]interface{}{
		"name": "New Name",
	}
	updateJsonBody, _ := json.Marshal(updateBody)
	updateReq := httptest.NewRequest("PUT", "/api/v1/cupcakes/invalid", bytes.NewBuffer(updateJsonBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateW := httptest.NewRecorder()
	router.ServeHTTP(updateW, updateReq)

	require.Equal(t, http.StatusBadRequest, updateW.Code)
	require.Contains(t, updateW.Body.String(), "Invalid ID")
}

func TestDeleteCupcake_Returns204(t *testing.T) {
	router := newTestRouter(t)

	reqBody := map[string]interface{}{
		"name":        "To Delete",
		"flavor":      "Delete Flavor",
		"price_cents": 1000,
	}
	jsonBody, _ := json.Marshal(reqBody)
	createReq := httptest.NewRequest("POST", "/api/v1/cupcakes", bytes.NewBuffer(jsonBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)
	require.Equal(t, http.StatusCreated, createW.Code)

	var created models.Cupcake
	err := json.Unmarshal(createW.Body.Bytes(), &created)
	require.NoError(t, err)

	deleteReq := httptest.NewRequest("DELETE", "/api/v1/cupcakes/"+fmt.Sprintf("%d", created.ID), nil)
	deleteW := httptest.NewRecorder()
	router.ServeHTTP(deleteW, deleteReq)

	require.Equal(t, http.StatusNoContent, deleteW.Code)
	require.Empty(t, deleteW.Body.String())

	getReq := httptest.NewRequest("GET", "/api/v1/cupcakes/"+fmt.Sprintf("%d", created.ID), nil)
	getW := httptest.NewRecorder()
	router.ServeHTTP(getW, getReq)

	require.Equal(t, http.StatusNotFound, getW.Code)
}

func TestDeleteCupcake_NotFound_Returns400(t *testing.T) {
	router := newTestRouter(t)

	deleteReq := httptest.NewRequest("DELETE", "/api/v1/cupcakes/999", nil)
	deleteW := httptest.NewRecorder()
	router.ServeHTTP(deleteW, deleteReq)

	require.Equal(t, http.StatusBadRequest, deleteW.Code)
	require.Contains(t, deleteW.Body.String(), "cupcake not found")
}

func TestDeleteCupcake_InvalidID_Returns400(t *testing.T) {
	router := newTestRouter(t)

	deleteReq := httptest.NewRequest("DELETE", "/api/v1/cupcakes/invalid", nil)
	deleteW := httptest.NewRecorder()
	router.ServeHTTP(deleteW, deleteReq)

	require.Equal(t, http.StatusBadRequest, deleteW.Code)
	require.Contains(t, deleteW.Body.String(), "Invalid ID")
}
