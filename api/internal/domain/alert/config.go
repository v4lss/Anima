// Package alert - AlertConfig defines notification settings for a monitor.
package alert

import "time"

// Config stores notification settings for a monitor.
type Config struct {
	ID        string     `bson:"id"`
	MonitorID string     `bson:"monitorid"`
	UserID    string     `bson:"userid"`
	Type      AlertType  `bson:"type"`
	Webhook   string     `bson:"webhook"`   // Discord webhook URL or email address
	Enabled   bool       `bson:"enabled"`
	CreatedAt time.Time  `bson:"createdat"`
	UpdatedAt time.Time  `bson:"updatedat"`
}
