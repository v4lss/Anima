// Package monitor - repository interface.
// Defined in the domain so the application layer can depend on it
// without knowing the concrete database implementation.
package monitor

import "context"

// Repository is the port that infrastructure must satisfy.
type Repository interface {
	Create(ctx context.Context, m *Monitor) error
	FindByID(ctx context.Context, id string) (*Monitor, error)
	FindByUserID(ctx context.Context, userID string) ([]*Monitor, error)
	FindAllEnabled(ctx context.Context) ([]*Monitor, error)
	Update(ctx context.Context, m *Monitor) error
	Delete(ctx context.Context, id string) error
}

// CheckRepository persists the result of each probe.
type CheckRepository interface {
	Save(ctx context.Context, c *Check) error
	FindByMonitorID(ctx context.Context, monitorID string, limit int) ([]*Check, error)
}
