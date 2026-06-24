// Package alert - repository port.
package alert

import "context"

// Repository persists Alert records.
type Repository interface {
	Save(ctx context.Context, a *Alert) error
	FindByMonitorID(ctx context.Context, monitorID string) ([]*Alert, error)
}
