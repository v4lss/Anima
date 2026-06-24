// Package workers - HTTPChecker consumes the "jobs:http" queue
// and performs HTTP/HTTPS probes, saving Check results to MongoDB.
package workers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/v4lss/animas/internal/domain/monitor"
	redisinfra "github.com/v4lss/animas/internal/infrastructure/redis"
)

// HTTPChecker polls the jobs:http Redis queue and checks HTTP monitors.
type HTTPChecker struct {
	monitorRepo monitor.Repository
	checkRepo   monitor.CheckRepository
	queue       *redisinfra.Queue
	client      *http.Client
}

func NewHTTPChecker(mr monitor.Repository, cr monitor.CheckRepository, q *redisinfra.Queue) *HTTPChecker {
	return &HTTPChecker{
		monitorRepo: mr,
		checkRepo:   cr,
		queue:       q,
		client:      &http.Client{Timeout: 10 * time.Second},
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
}
