// Package alert defines the Alert entity.
// Alerts are fired when a monitor transitions from UP to DOWN (or vice-versa).
package alert

import "time"

// AlertType identifies the notification channel.
type AlertType string

const (
	Email   AlertType = "EMAIL"
	Discord AlertType = "DISCORD"
)

// Alert records a single notification that was sent.
type Alert struct {
	ID        string     `bson:"id"`
	MonitorID string     `bson:"monitorid"`
	UserID    string     `bson:"userid"`
	Type      AlertType  `bson:"type"`
	Message   string     `bson:"message"`
	SentAt    time.Time  `bson:"sentat"`
}
