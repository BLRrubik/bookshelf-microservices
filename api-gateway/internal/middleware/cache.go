package middleware

import (
	"bookshelf/api-gateway/internal/cache"
	"net/http"
	"time"
)

type CacheConfig struct {
	Cache     *cache.Cache
	TTL       time.Duration
	KeyPrefix string
}

func CacheMiddleware(cfg *CacheConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				next.ServeHTTP(w, r)

				return
			}

			key := cfg.Cache.GenerateKey(cfg.KeyPrefix, r.URL.Path, r.URL.RawQuery)

			res, err := cfg.Cache.Get(r.Context(), key)
			if err != nil {
				w.Header().Set("X-Cache", "MISS")
				next.ServeHTTP(w, r)

				return
			}

			switch {
			case len(res) > 0:
				w.WriteHeader(http.StatusOK)
				w.Header().Set("X-Cache", "HIT")
				w.Write(res)

				return
			default:
				w.Header().Set("X-Cache", "MISS")

				next.ServeHTTP(w, r)
			}
		})
	}
}
