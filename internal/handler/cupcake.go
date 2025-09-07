package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/julimonteiro/cupcake-store/internal/models"
	"github.com/julimonteiro/cupcake-store/internal/service"
)

func sendJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

type CupcakeHandler struct {
	service *service.CupcakeService
}

func NewCupcakeHandler(service *service.CupcakeService) *CupcakeHandler {
	return &CupcakeHandler{service: service}
}

func (h *CupcakeHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"message": "Cupcake Store API is running!",
	})
}

func (h *CupcakeHandler) CreateCupcake(w http.ResponseWriter, r *http.Request) {
	var req models.CreateCupcakeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONError(w, "Error decoding request", http.StatusBadRequest)
		return
	}

	cupcake, err := h.service.CreateCupcake(&req)
	if err != nil {
		sendJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cupcake)
}

func (h *CupcakeHandler) GetCupcake(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		sendJSONError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	cupcake, err := h.service.GetCupcake(uint(id))
	if err != nil {
		sendJSONError(w, "cupcake not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cupcake)
}

func (h *CupcakeHandler) GetAllCupcakes(w http.ResponseWriter, r *http.Request) {
	cupcakes, err := h.service.GetAllCupcakes()
	if err != nil {
		sendJSONError(w, "Error fetching cupcakes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cupcakes)
}

func (h *CupcakeHandler) UpdateCupcake(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		sendJSONError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateCupcakeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONError(w, "Error decoding request", http.StatusBadRequest)
		return
	}

	cupcake, err := h.service.UpdateCupcake(uint(id), &req)
	if err != nil {
		sendJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cupcake)
}

func (h *CupcakeHandler) DeleteCupcake(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		sendJSONError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteCupcake(uint(id)); err != nil {
		sendJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
