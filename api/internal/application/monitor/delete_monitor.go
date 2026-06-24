// Package monitor (application layer) - use-case: DeleteMonitor.
package monitor

import (
	"context"
	"errors"

	"github.com/v4lss/animas/internal/domain/monitor"
)

// DeleteMonitor removes a monitor if it belongs to the requesting user.
func DeleteMonitor(ctx context.Context, repo monitor.Repository, id, userID string) error {
	m, err := repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if m.UserID != userID {
		return errors.New("forbidden")
	}
	return repo.Delete(ctx, id)
}
