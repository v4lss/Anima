// Package workers - Scheduler periodically enqueues monitor IDs for checking.
// Every 10 seconds it fetches all enabled monitors and pushes those whose
// interval has elapsed since their last check into the appropriate Redis queue.
package workers

import (
	"context"
	"sync"
	"time"

	"github.com/v4lss/animas/internal/domain/monitor"
	redisinfra "github.com/v4lss/animas/internal/infrastructure/redis"
	"github.com/v4lss/animas/pkg/logger"
)

// Scheduler enqueues monitors that are due for a check.
type Scheduler struct {
	repo      monitor.Repository
	queue     *redisinfra.Queue
	lastCheck sync.Map // monitorID, time.Time of last enqueue
}

func NewScheduler(repo monitor.Repository, queue *redisinfra.Queue) *Scheduler {
	return &Scheduler{repo: repo, queue: queue}
}

// Run blocks and ticks every 10 seconds.
func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	// Enqueue immediately on startup so we don't wait 10s for the first check.
	s.enqueue(ctx)

	for {
		select {
		case <-ticker.C:
			s.enqueue(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (s *Scheduler) enqueue(ctx context.Context) {
	monitors, err := s.repo.FindAllEnabled(ctx)
	if err != nil {
		logger.Error("[scheduler] failed to fetch monitors: %v\n", err)
		return
	}

	now := time.Now()
	queued := 0

	for _, m := range monitors {
		if !s.isDue(m, now) {
			continue
		}

		queue := queueFor(m.Type)
		if err := s.queue.Push(ctx, queue, m.ID); err != nil {
			logger.Error("[scheduler] failed to enqueue monitor %s: %v\n", m.ID, err)
			continue
		}

		s.lastCheck.Store(m.ID, now)
		queued++
	}

	if queued > 0 {
		logger.Info("[scheduler] enqueued %d monitor(s)\n", queued)
	}
}

// isDue returns true if the monitor hasn't been enqueued within its interval.
func (s *Scheduler) isDue(m *monitor.Monitor, now time.Time) bool {
	last, ok := s.lastCheck.Load(m.ID)
	if !ok {
		return true // never checked, enqueue immediately
	}
	return now.Sub(last.(time.Time)) >= time.Duration(m.Interval)*time.Second
}

// queueFor maps a MonitorType to its Redis queue name.
func queueFor(t monitor.MonitorType) string {
	if t == monitor.TCP {
		return "jobs:tcp"
	}
	return "jobs:http"
}
