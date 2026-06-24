// Package monitor - domain service.
// Contains business rules that don't naturally belong to a single entity.
package monitor

import (
	"errors"
	"strings"
)

// Validate checks that a Monitor has the minimum required fields.
func Validate(m *Monitor) error {
	if strings.TrimSpace(m.Name) == "" {
		return errors.New("monitor name is required")
	}
	if strings.TrimSpace(m.Target) == "" {
		return errors.New("monitor target is required")
	}
	if m.Interval < 30 {
		return errors.New("interval must be at least 30 seconds")
	}
	if m.Type != HTTP && m.Type != HTTPS && m.Type != TCP {
		return errors.New("invalid monitor type")
	}
	return nil
}
