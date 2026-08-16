package middleware

import (
	"bookshelf/api-gateway/internal/ratelimit"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type RateLimitConfig struct {
	Limiter         *ratelimit.RateLimiter
	AnonymousLimit  int
	AuthorizedLimit int
	UploadLimit     int
}

func RateLimitMiddleware(cfg *RateLimitConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			key, limit := getRateLimitKeyAndLimit(r, cfg)
			res, err := cfg.Limiter.Allow(r.Context(), key, limit)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			setRateLimitHeader(w, res.Limit, res.Remaining, res.ResetAt)

			if !res.Allowed {
				w.Header().Set("Retry-After", strconv.Itoa(int(res.ResetAt.Sub(time.Now()).Seconds())))
				w.WriteHeader(http.StatusTooManyRequests)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func setRateLimitHeader(w http.ResponseWriter, limit, remaining int, resetAt time.Time) {
	w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
	w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
	w.Header().Set("X-RateLimit-Reset", strconv.Itoa(int(resetAt.Unix())))
}

func getRateLimitKeyAndLimit(r *http.Request, cfg *RateLimitConfig) (string, int) {
	ip := getClientIP(r)
	if strings.Contains(r.RequestURI, "/cover") {
		return fmt.Sprintf("cover:ip:%s", ip), cfg.UploadLimit
	}

	authHdr := r.Header.Get("Authorization")
	if authHdr == "" || !strings.HasPrefix(authHdr, "Bearer ") {
		return fmt.Sprintf("anon:ip:%s", ip), cfg.AnonymousLimit
	}

	token := strings.TrimPrefix(authHdr, "Bearer ")

	return fmt.Sprintf("auth:ip:%s", token[:7]), cfg.UploadLimit
}

func getClientIP(r *http.Request) string {
	splitted := strings.Split(r.Header.Get("X-Forwarded-For"), ",")
	if len(splitted[0]) > 0 {
		return splitted[0]
	}

	ip := r.Header.Get("X-Real-Ip")
	if len(ip) > 0 {
		return ip
	}

	return r.RemoteAddr
}
