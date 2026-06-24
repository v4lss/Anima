// Package redis - Queue pushes and pops monitor job IDs via Redis lists.
// jobs:http and jobs:tcp are the two queues consumed by workers.
package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// Queue wraps Redis list operations as a simple FIFO job queue.
type Queue struct {
	client *redis.Client
}

func NewQueue(addr, password string) *Queue {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})
	return &Queue{client: rdb}
}

// Client returns the underlying Redis client.
func (q *Queue) Client() *redis.Client {
	return q.client
}

// Push adds a monitorID to the given queue name (e.g. "jobs:http").
func (q *Queue) Push(ctx context.Context, queue, monitorID string) error {
	return q.client.RPush(ctx, queue, monitorID).Err()
}

// Pop blocks until a job is available and returns the monitorID.
func (q *Queue) Pop(ctx context.Context, queue string) (string, error) {
	result, err := q.client.BLPop(ctx, 0, queue).Result()
	if err != nil {
		return "", err
	}
	return result[1], nil // result[0] is the queue name
}
