package middleware

import (
	"bookshelf/api-gateway/internal/ratelimit"
	"net/http"
	"strconv"
	"time"
)

type RateLimitConfig struct {
	Limiter *ratelimit.RateLimiter
	Limit   int
}

func RateLimitMiddleware(cfg *RateLimitConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			res, err := cfg.Limiter.Allow(r.Context(), r.RemoteAddr, cfg.Limit)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			switch {
			case res.Allowed:
				setRateLimitHeader(w, res.Limit, res.Remaining, res.ResetAt)

				next.ServeHTTP(w, r)

			default:
				w.Header().Set("Retry-After", strconv.Itoa(int(res.ResetAt.Unix())))
				w.WriteHeader(http.StatusTooManyRequests)

				return
			}
		})
	}
}

func setRateLimitHeader(w http.ResponseWriter, limit, remaining int, resetAt time.Time) {
	w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
	w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
	w.Header().Set("X-RateLimit-Reset", strconv.Itoa(int(resetAt.Unix())))
}
