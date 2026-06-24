// Package middleware - RateLimit provides IP-based rate limiting using Redis.
package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/v4lss/animas/pkg/response"
)

// RateLimit returns middleware that limits requests per IP.
// Uses Redis sliding window counter with 1 minute windows.
func RateLimit(redisClient *redis.Client, requestsPerMinute int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ip := getClientIP(r)
			key := fmt.Sprintf("ratelimit:%s", ip)

			// Increment counter with 1 minute expiry
			pipe := redisClient.Pipeline()
			incr := pipe.Incr(ctx, key)
			pipe.Expire(ctx, key, time.Minute)
			_, err := pipe.Exec(ctx)

			if err != nil {
				// Redis error: allow request but log
				next.ServeHTTP(w, r)
				return
			}

			count := incr.Val()
			if count > int64(requestsPerMinute) {
				response.Error(w, http.StatusTooManyRequests, "rate limit exceeded")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxies)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	return r.RemoteAddr
}
