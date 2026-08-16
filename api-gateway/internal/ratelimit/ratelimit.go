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

func (r *RateLimiter) Allow(ctx context.Context, key string, limit int) (*Result, error) {
	fullKey := rateLimitKey + key

	rate, err := r.client.Incr(ctx, fullKey).Result()
	if err != nil {
		return nil, err
	}

	if rate == 1 {
		r.client.Expire(ctx, fullKey, r.window)
	}

	remainingTTL, err := r.client.TTL(ctx, fullKey).Result()
	if err != nil {
		return nil, err
	}

	return &Result{
		Allowed:   rate <= int64(limit),
		Limit:     limit,
		Remaining: limit - int(rate),
		ResetAt:   time.Now().Add(remainingTTL),
	}, nil
}
