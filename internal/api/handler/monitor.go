package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AksanovK/url-monitor/internal/service"
)

type MonitorHandler struct {
	service *service.MonitorService
}

func NewMonitorHandler(s *service.MonitorService) *MonitorHandler {
	return &MonitorHandler{service: s}
}

func (h *MonitorHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL            string `json:"url"`
		IntervalSec    int    `json:"interval_sec"`
		ExpectedStatus int    `json:"expected_status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	m, err := h.service.Create(r.Context(), req.URL, req.IntervalSec, req.ExpectedStatus)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(m); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func (h *MonitorHandler) List(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.List(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func (h *MonitorHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	m, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(m); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
