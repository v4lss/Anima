// Package monitor defines the core Monitor entity and its value objects.
// This is pure domain - no framework, no database, no HTTP dependencies.
package monitor

import "time"

// MonitorType classifies what protocol the monitor checks.
type MonitorType string

const (
	HTTP  MonitorType = "HTTP"
	HTTPS MonitorType = "HTTPS"
	TCP   MonitorType = "TCP"
)

// Monitor is the central aggregate of Animas.
// It describes a target that should be periodically checked.
type Monitor struct {
	ID        string
	UserID    string
	Name      string
	Target    string      // URL for HTTP/HTTPS, "host:port" for TCP
	Type      MonitorType
	Interval  int  // seconds between checks
	Enabled   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Check represents a single probe result for a Monitor.
type Check struct {
	ID           string
	MonitorID    string
	Status       CheckStatus
	ResponseTime int64 // milliseconds
	Error        string
	CheckedAt    time.Time
}

// CheckStatus indicates whether the probe succeeded or failed.
type CheckStatus string

const (
	StatusUp   CheckStatus = "UP"
	StatusDown CheckStatus = "DOWN"
)
