// Package workers - TCPChecker consumes the "jobs:tcp" queue
// and performs TCP dial probes.
package workers

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/v4lss/animas/internal/domain/alert"
	"github.com/v4lss/animas/internal/domain/monitor"
	redisinfra "github.com/v4lss/animas/internal/infrastructure/redis"
	"github.com/v4lss/animas/pkg/notifier"
)

// TCPChecker polls the jobs:tcp Redis queue and dials TCP monitors.
type TCPChecker struct {
	monitorRepo        monitor.Repository
	checkRepo          monitor.CheckRepository
	alertConfigRepo    alert.ConfigRepository
	alertRepo          alert.Repository
	queue              *redisinfra.Queue
	discordNotifier    *notifier.DiscordNotifier
}

func NewTCPChecker(
	mr monitor.Repository,
	cr monitor.CheckRepository,
	acr alert.ConfigRepository,
	ar alert.Repository,
	q *redisinfra.Queue,
	dn *notifier.DiscordNotifier,
) *TCPChecker {
	return &TCPChecker{
		monitorRepo:     mr,
		checkRepo:       cr,
		alertConfigRepo: acr,
		alertRepo:       ar,
		queue:           q,
		discordNotifier: dn,
	}
}

// Run blocks and processes jobs from the TCP queue.
func (c *TCPChecker) Run(ctx context.Context) {
	for {
		monitorID, err := c.queue.Pop(ctx, "jobs:tcp")
		if err != nil {
			continue
		}
		c.check(ctx, monitorID)
	}
}

func (c *TCPChecker) check(ctx context.Context, monitorID string) {
	m, err := c.monitorRepo.FindByID(ctx, monitorID)
	if err != nil {
		log.Printf("[tcp_checker] monitor not found: %s", monitorID)
		return
	}

	// Get previous check to detect status change
	prevChecks, _ := c.checkRepo.FindByMonitorID(ctx, monitorID, 1)
	var prevStatus monitor.CheckStatus
	if len(prevChecks) > 0 {
		prevStatus = prevChecks[0].Status
	}

	start := time.Now()
	conn, err := net.DialTimeout("tcp", m.Target, 5*time.Second)
	elapsed := time.Since(start).Milliseconds()

	chk := &monitor.Check{
		MonitorID:    monitorID,
		ResponseTime: elapsed,
		CheckedAt:    time.Now(),
	}

	if err != nil {
		chk.Status = monitor.StatusDown
		chk.Error = err.Error()
	} else {
		conn.Close()
		chk.Status = monitor.StatusUp
	}

	if err := c.checkRepo.Save(ctx, chk); err != nil {
		log.Printf("[tcp_checker] failed to save check: %v", err)
	}

	// Check for status change and send alerts
	if prevStatus != "" && prevStatus != chk.Status {
		c.sendAlerts(ctx, m, prevStatus, chk.Status, chk.Error)
	}
}

func (c *TCPChecker) sendAlerts(ctx context.Context, m *monitor.Monitor, oldStatus, newStatus monitor.CheckStatus, errorMsg string) {
	configs, err := c.alertConfigRepo.FindByMonitorID(ctx, m.ID)
	if err != nil {
		log.Printf("[tcp_checker] failed to get alert configs: %v", err)
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
				log.Printf("[tcp_checker] discord notify failed: %v", err)
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
					log.Printf("[tcp_checker] failed to save alert record: %v", err)
				}
			}
		}
	}
}
