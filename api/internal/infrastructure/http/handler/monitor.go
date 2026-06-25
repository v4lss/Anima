// Package handler - HTTP handlers for monitor endpoints.
package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	appmonitor "github.com/v4lss/animas/internal/application/monitor"
	"github.com/v4lss/animas/internal/domain/monitor"
	"github.com/v4lss/animas/internal/infrastructure/http/middleware"
	"github.com/v4lss/animas/pkg/response"
)

// MonitorHandler holds dependencies for monitor routes.
type MonitorHandler struct {
	monitorRepo monitor.Repository
	checkRepo   monitor.CheckRepository
}

func NewMonitorHandler(monitorRepo monitor.Repository, checkRepo monitor.CheckRepository) *MonitorHandler {
	return &MonitorHandler{monitorRepo: monitorRepo, checkRepo: checkRepo}
}

// Create handles POST /api/monitors
func (h *MonitorHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	var body struct {
		Name     string             `json:"name"`
		Target   string             `json:"target"`
		Type     monitor.MonitorType `json:"type"`
		Interval int                `json:"interval"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	m, err := appmonitor.CreateMonitor(r.Context(), h.monitorRepo, appmonitor.CreateMonitorInput{
		UserID:   userID,
		Name:     body.Name,
		Target:   body.Target,
		Type:     body.Type,
		Interval: body.Interval,
	})
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, http.StatusCreated, m)
}

// List handles GET /api/monitors
func (h *MonitorHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	monitors, err := appmonitor.ListMonitors(r.Context(), h.monitorRepo, userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Attach last check status to each monitor
	type monitorWithStatus struct {
		*monitor.Monitor
		LastStatus monitor.CheckStatus `json:"lastStatus"`
	}

	result := make([]monitorWithStatus, len(monitors))
	for i, m := range monitors {
		checks, err := h.checkRepo.FindByMonitorID(r.Context(), m.ID, 1)
		lastStatus := monitor.StatusDown
		if err == nil && len(checks) > 0 {
			lastStatus = checks[0].Status
		}
		result[i] = monitorWithStatus{Monitor: m, LastStatus: lastStatus}
	}

	response.Success(w, http.StatusOK, result)
}

// Get handles GET /api/monitors/:id
func (h *MonitorHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	m, err := h.monitorRepo.FindByID(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "monitor not found")
		return
	}

	response.Success(w, http.StatusOK, m)
}

// Delete handles DELETE /api/monitors/:id
func (h *MonitorHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	id := chi.URLParam(r, "id")

	if err := appmonitor.DeleteMonitor(r.Context(), h.monitorRepo, id, userID); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		}
		response.Error(w, status, err.Error())
		return
	}

	response.Success(w, http.StatusOK, map[string]string{"deleted": id})
}

// History handles GET /api/monitors/:id/history
func (h *MonitorHandler) History(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	status := r.URL.Query().Get("status")
	fromDate := r.URL.Query().Get("from")
	toDate := r.URL.Query().Get("to")

	checks, total, err := h.checkRepo.FindByMonitorIDPaginated(r.Context(), id, page, limit, status, fromDate, toDate)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusOK, map[string]interface{}{
		"data":  checks,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}
