// Package alert - config repository port.
package alert

import "context"

// ConfigRepository persists AlertConfig records.
type ConfigRepository interface {
	Save(ctx context.Context, c *Config) error
	FindByMonitorID(ctx context.Context, monitorID string) ([]*Config, error)
	Delete(ctx context.Context, id string) error
}
