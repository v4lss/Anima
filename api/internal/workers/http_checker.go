// Package workers - HTTPChecker consumes the "jobs:http" queue
// and performs HTTP/HTTPS probes, saving Check results to MongoDB.
package workers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/v4lss/animas/internal/domain/alert"
	"github.com/v4lss/animas/internal/domain/monitor"
	redisinfra "github.com/v4lss/animas/internal/infrastructure/redis"
	"github.com/v4lss/animas/pkg/notifier"
)

// HTTPChecker polls the jobs:http Redis queue and checks HTTP monitors.
type HTTPChecker struct {
	monitorRepo        monitor.Repository
	checkRepo          monitor.CheckRepository
	alertConfigRepo    alert.ConfigRepository
	alertRepo          alert.Repository
	queue              *redisinfra.Queue
	client             *http.Client
	discordNotifier    *notifier.DiscordNotifier
}

func NewHTTPChecker(
	mr monitor.Repository,
	cr monitor.CheckRepository,
	acr alert.ConfigRepository,
	ar alert.Repository,
	q *redisinfra.Queue,
	dn *notifier.DiscordNotifier,
) *HTTPChecker {
	return &HTTPChecker{
		monitorRepo:     mr,
		checkRepo:       cr,
		alertConfigRepo: acr,
		alertRepo:       ar,
		queue:           q,
		client:          &http.Client{Timeout: 10 * time.Second},
		discordNotifier: dn,
	}
}

// Run blocks and processes jobs from the HTTP queue.
func (c *HTTPChecker) Run(ctx context.Context) {
	for {
		monitorID, err := c.queue.Pop(ctx, "jobs:http")
		if err != nil {
			continue
		}
		c.check(ctx, monitorID)
	}
}

func (c *HTTPChecker) check(ctx context.Context, monitorID string) {
	m, err := c.monitorRepo.FindByID(ctx, monitorID)
	if err != nil {
		log.Printf("[http_checker] monitor not found: %s", monitorID)
		return
	}

	// Get previous check to detect status change
	prevChecks, _ := c.checkRepo.FindByMonitorID(ctx, monitorID, 1)
	var prevStatus monitor.CheckStatus
	if len(prevChecks) > 0 {
		prevStatus = prevChecks[0].Status
	}

	start := time.Now()
	resp, err := c.client.Get(m.Target)
	elapsed := time.Since(start).Milliseconds()

	chk := &monitor.Check{
		MonitorID:    monitorID,
		ResponseTime: elapsed,
		CheckedAt:    time.Now(),
	}

	if err != nil || resp.StatusCode >= 500 {
		chk.Status = monitor.StatusDown
		if err != nil {
			chk.Error = err.Error()
		}
	} else {
		chk.Status = monitor.StatusUp
	}

	if err := c.checkRepo.Save(ctx, chk); err != nil {
		log.Printf("[http_checker] failed to save check: %v", err)
	}

	// Check for status change and send alerts
	if prevStatus != "" && prevStatus != chk.Status {
		c.sendAlerts(ctx, m, prevStatus, chk.Status, chk.Error)
	}
}

func (c *HTTPChecker) sendAlerts(ctx context.Context, m *monitor.Monitor, oldStatus, newStatus monitor.CheckStatus, errorMsg string) {
	configs, err := c.alertConfigRepo.FindByMonitorID(ctx, m.ID)
	if err != nil {
		log.Printf("[http_checker] failed to get alert configs: %v", err)
		return
	}

	for _, cfg := range configs {
		if !cfg.Enabled {
			continue
		}

		if cfg.Type == alert.Discord {
			message := "Monitor is UP"
			if newStatus == monitor.StatusDown {
				message = errorMsg
			}

			if err := c.discordNotifier.Notify(cfg.Webhook, m.Name, m.Target, string(newStatus), message); err != nil {
				log.Printf("[http_checker] discord notify failed: %v", err)
			} else {
				// Record alert
				alertRecord := &alert.Alert{
					MonitorID: m.ID,
					UserID:    m.UserID,
					Type:      alert.Discord,
					Message:   message,
					SentAt:    time.Now(),
				}
				if err := c.alertRepo.Save(ctx, alertRecord); err != nil {
					log.Printf("[http_checker] failed to save alert record: %v", err)
				}
			}
		}
	}
}
