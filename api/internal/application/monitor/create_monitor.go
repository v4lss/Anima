// Package monitor (application layer) - use-case: CreateMonitor.
// Orchestrates domain validation + persistence via the repository port.
package monitor

import (
	"context"
	"time"

	"github.com/v4lss/animas/internal/domain/monitor"
)

// CreateMonitorInput carries the data from the HTTP handler.
type CreateMonitorInput struct {
	UserID   string
	Name     string
	Target   string
	Type     monitor.MonitorType
	Interval int
}

// CreateMonitor validates and persists a new monitor.
func CreateMonitor(ctx context.Context, repo monitor.Repository, in CreateMonitorInput) (*monitor.Monitor, error) {
	m := &monitor.Monitor{
		UserID:    in.UserID,
		Name:      in.Name,
		Target:    in.Target,
		Type:      in.Type,
		Interval:  in.Interval,
		Enabled:   true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := monitor.Validate(m); err != nil {
		return nil, err
	}

	if err := repo.Create(ctx, m); err != nil {
		return nil, err
	}

	return m, nil
}
