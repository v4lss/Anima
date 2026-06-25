// Package handler - AlertHandler manages alert configuration endpoints.
package handler

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/v4lss/animas/internal/domain/alert"
	"github.com/v4lss/animas/internal/domain/monitor"
	"github.com/v4lss/animas/internal/infrastructure/http/middleware"
	"github.com/v4lss/animas/pkg/response"
)

// AlertHandler handles alert configuration requests.
type AlertHandler struct {
	alertConfigRepo alert.ConfigRepository
	monitorRepo     monitor.Repository
}

func NewAlertHandler(acr alert.ConfigRepository, mr monitor.Repository) *AlertHandler {
	return &AlertHandler{alertConfigRepo: acr, monitorRepo: mr}
}

// Create handles POST /api/monitors/:monitorId/alerts
func (h *AlertHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	monitorID := chi.URLParam(r, "monitorId")

	// Verify user owns this monitor
	m, err := h.monitorRepo.FindByID(r.Context(), monitorID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "monitor not found")
		return
	}
	if m.UserID != userID {
		response.Error(w, http.StatusForbidden, "forbidden")
		return
	}

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

	// Validate webhook URL format (basic check)
	if !strings.HasPrefix(input.Webhook, "https://") && !strings.HasPrefix(input.Webhook, "http://") {
		response.Error(w, http.StatusBadRequest, "webhook must be a valid URL")
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
	userID := middleware.UserIDFromContext(r.Context())
	monitorID := chi.URLParam(r, "monitorId")

	// Verify user owns this monitor
	m, err := h.monitorRepo.FindByID(r.Context(), monitorID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "monitor not found")
		return
	}
	if m.UserID != userID {
		response.Error(w, http.StatusForbidden, "forbidden")
		return
	}

	configs, err := h.alertConfigRepo.FindByMonitorID(r.Context(), monitorID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusOK, configs)
}

// Delete handles DELETE /api/monitors/:monitorId/alerts/:id
func (h *AlertHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	id := chi.URLParam(r, "id")

	// Get the alert config to verify ownership
	configs, err := h.alertConfigRepo.FindByMonitorID(r.Context(), chi.URLParam(r, "monitorId"))
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	var targetConfig *alert.Config
	for _, cfg := range configs {
		if cfg.ID == id {
			targetConfig = cfg
			break
		}
	}

	if targetConfig == nil {
		response.Error(w, http.StatusNotFound, "alert config not found")
		return
	}

	if targetConfig.UserID != userID {
		response.Error(w, http.StatusForbidden, "forbidden")
		return
	}

	if err := h.alertConfigRepo.Delete(r.Context(), id); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusNoContent, nil)
}
