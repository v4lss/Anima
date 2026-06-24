// Package redis - Cache wraps go-redis for simple key/value operations.
package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache provides a thin wrapper over the Redis client.
type Cache struct {
	client *redis.Client
}

func NewCache(addr, password string) *Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})
	return &Cache{client: rdb}
}

func (c *Cache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
