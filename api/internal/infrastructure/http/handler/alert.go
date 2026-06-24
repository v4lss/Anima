// Package handler - AlertHandler manages alert configuration endpoints.
package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/v4lss/animas/internal/domain/alert"
	"github.com/v4lss/animas/internal/infrastructure/http/middleware"
	"github.com/v4lss/animas/pkg/response"
)

// AlertHandler handles alert configuration requests.
type AlertHandler struct {
	alertConfigRepo alert.ConfigRepository
}

func NewAlertHandler(acr alert.ConfigRepository) *AlertHandler {
	return &AlertHandler{alertConfigRepo: acr}
}

// Create handles POST /api/monitors/:monitorId/alerts
func (h *AlertHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	monitorID := chi.URLParam(r, "monitorId")

	var input struct {
		Type    alert.AlertType `json:"type"`
		Webhook string          `json:"webhook"`
	}

	if err := response.DecodeJSON(r, &input); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if input.Type == "" || input.Webhook == "" {
		response.Error(w, http.StatusBadRequest, "type and webhook are required")
		return
	}

	config := &alert.Config{
		MonitorID: monitorID,
		UserID:    userID,
		Type:      input.Type,
		Webhook:   input.Webhook,
		Enabled:   true,
	}

	if err := h.alertConfigRepo.Save(r.Context(), config); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusCreated, config)
}

// List handles GET /api/monitors/:monitorId/alerts
func (h *AlertHandler) List(w http.ResponseWriter, r *http.Request) {
	monitorID := chi.URLParam(r, "monitorId")

	configs, err := h.alertConfigRepo.FindByMonitorID(r.Context(), monitorID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusOK, configs)
}

// Delete handles DELETE /api/monitors/:monitorId/alerts/:id
func (h *AlertHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.alertConfigRepo.Delete(r.Context(), id); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusNoContent, nil)
}
