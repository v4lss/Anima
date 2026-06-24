// Package monitor (application layer) - use-case: ListMonitors.
package monitor

import (
	"context"

	"github.com/v4lss/animas/internal/domain/monitor"
)

// ListMonitors returns all monitors owned by a user.
func ListMonitors(ctx context.Context, repo monitor.Repository, userID string) ([]*monitor.Monitor, error) {
	return repo.FindByUserID(ctx, userID)
}
