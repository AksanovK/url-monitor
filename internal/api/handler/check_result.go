package handler

import (
	"encoding/json"
	"github.com/AksanovK/url-monitor/internal/domain"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/AksanovK/url-monitor/internal/service"
)

type CheckResultHandler struct {
	service *service.CheckResultService
}

func NewCheckResultHandler(s *service.CheckResultService) *CheckResultHandler {
	return &CheckResultHandler{service: s}
}

type CheckResultsResponse struct {
	Data       []*domain.CheckResult `json:"data"`
	NextCursor string                `json:"next_cursor,omitempty"`
	HasMore    bool                  `json:"has_more"`
}

func (h *CheckResultHandler) List(w http.ResponseWriter, r *http.Request) {
	monitorID := r.PathValue("id")

	limitStr := r.URL.Query().Get("limit")
	limit := 20
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil {
			limit = parsed
		}
	}

	var cursor *time.Time
	cursorStr := r.URL.Query().Get("cursor")
	if cursorStr != "" {
		parsed, err := time.Parse(time.RFC3339, cursorStr)
		if err != nil {
			http.Error(w, "invalid cursor format, use RFC3339", http.StatusBadRequest)
			return
		}
		cursor = &parsed
	}

	results, err := h.service.GetByMonitor(r.Context(), monitorID, cursor, limit)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	response := CheckResultsResponse{
		Data:    results,
		HasMore: len(results) == limit,
	}

	if len(results) > 0 {
		last := results[len(results)-1]
		response.NextCursor = last.CheckedAt.UTC().Format(time.RFC3339)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
