// Package workers - Scheduler periodically enqueues monitor IDs for checking.
// It reads all enabled monitors from the DB and pushes them to the Redis queue
// according to their configured interval.
package workers

import (
	"context"
	"log"
	"time"

	"github.com/v4lss/animas/internal/domain/monitor"
	redisinfra "github.com/v4lss/animas/internal/infrastructure/redis"
)

// Scheduler ticks every minute and enqueues due monitors.
type Scheduler struct {
	repo  monitor.Repository
	queue *redisinfra.Queue
}

func NewScheduler(repo monitor.Repository, queue *redisinfra.Queue) *Scheduler {
	return &Scheduler{repo: repo, queue: queue}
}

// Run blocks and ticks every 30 seconds.
func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

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
	// TODO: iterate all enabled monitors and push to the right queue.
	// Stub: implementation comes in v1 milestone.
	log.Println("[scheduler] tick - enqueue pass")
}
