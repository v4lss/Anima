// Package notifier provides notification services.
package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// DiscordWebhook represents a Discord webhook payload.
type DiscordWebhook struct {
	Content string  `json:"content"`
	Embeds  []Embed `json:"embeds,omitempty"`
}

// Embed represents a Discord embed object.
type Embed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Color       int    `json:"color"`
	Timestamp   string `json:"timestamp"`
}

// DiscordNotifier sends notifications to Discord webhooks.
type DiscordNotifier struct {
	client *http.Client
}

func NewDiscordNotifier() *DiscordNotifier {
	return &DiscordNotifier{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// Notify sends an alert to a Discord webhook.
func (d *DiscordNotifier) Notify(webhookURL, monitorName, target, status, message string) error {
	color := 0x00ff00 // green
	if status == "DOWN" {
		color = 0xff0000 // red
	}

	payload := DiscordWebhook{
		Content: fmt.Sprintf("🚨 Monitor Alert: **%s** is %s", monitorName, status),
		Embeds: []Embed{
			{
				Title:       monitorName,
				Description: fmt.Sprintf("**Target:** %s\n**Status:** %s\n**Message:** %s", target, status, message),
				Color:       color,
				Timestamp:   time.Now().Format(time.RFC3339),
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal webhook: %w", err)
	}

	req, err := http.NewRequest("POST", webhookURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.client.Do(req)
	if err != nil {
		return fmt.Errorf("send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	return nil
}
