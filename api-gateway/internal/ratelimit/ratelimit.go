package ratelimit

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const rateLimitKey = "ratelimit:"

type Result struct {
	Allowed   bool
	Limit     int
	Remaining int
	ResetAt   time.Time
}

type RateLimiter struct {
	client *redis.Client
	window time.Duration
}

func New(client *redis.Client, window time.Duration) *RateLimiter {
	return &RateLimiter{client: client, window: window}
}

func (r *RateLimiter) Allow(ctx context.Context, key string, limit int) bool {
	rate, err := r.client.Incr(ctx, rateLimitKey+key).Result()
	if err != nil {
		return false
	}

	if rate == 1 {
		r.client.Expire(ctx, rateLimitKey+key, r.window)
	}

	return rate <= int64(limit)
}
