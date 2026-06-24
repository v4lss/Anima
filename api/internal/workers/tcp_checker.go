// Package workers - TCPChecker consumes the "jobs:tcp" queue
// and performs TCP dial probes.
package workers

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/v4lss/animas/internal/domain/monitor"
	redisinfra "github.com/v4lss/animas/internal/infrastructure/redis"
)

// TCPChecker polls the jobs:tcp Redis queue and dials TCP monitors.
type TCPChecker struct {
	monitorRepo monitor.Repository
	checkRepo   monitor.CheckRepository
	queue       *redisinfra.Queue
}

func NewTCPChecker(mr monitor.Repository, cr monitor.CheckRepository, q *redisinfra.Queue) *TCPChecker {
	return &TCPChecker{monitorRepo: mr, checkRepo: cr, queue: q}
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
}
