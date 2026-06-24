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
	ID        string    `bson:"id"`
	UserID    string    `bson:"userid"`
	Name      string    `bson:"name"`
	Target    string    `bson:"target"` // URL for HTTP/HTTPS, "host:port" for TCP
	Type      MonitorType `bson:"type"`
	Interval  int       `bson:"interval"`
	Enabled   bool      `bson:"enabled"`
	CreatedAt time.Time `bson:"createdat"`
	UpdatedAt time.Time `bson:"updatedat"`
}

// Check represents a single probe result for a Monitor.
type Check struct {
	ID           string      `bson:"id"`
	MonitorID    string      `bson:"monitorid"`
	Status       CheckStatus `bson:"status"`
	ResponseTime int64       `bson:"responsetime"` // milliseconds
	Error        string      `bson:"error"`
	CheckedAt    time.Time   `bson:"checkedat"`
}

// CheckStatus indicates whether the probe succeeded or failed.
type CheckStatus string

const (
	StatusUp   CheckStatus = "UP"
	StatusDown CheckStatus = "DOWN"
)
