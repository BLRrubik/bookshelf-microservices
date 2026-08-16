package middleware

import (
	"bookshelf/api-gateway/internal/cache"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type CachedResponse struct {
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
	Body       []byte            `json:"body"`
}

type CacheConfig struct {
	Cache     *cache.Cache
	TTL       time.Duration
	KeyPrefix string
}

// uncacheableHeaders set per-request by other middleware (rate limit, request id).
// Storing stale copies would overwrite fresh values on cache HIT.
var uncacheableHeaders = map[string]struct{}{
	"X-Ratelimit-Limit":     {},
	"X-Ratelimit-Remaining": {},
	"X-Ratelimit-Reset":     {},
	"Retry-After":           {},
	"X-Request-Id":          {},
	"X-Cache":               {},
	"Set-Cookie":            {},
}

type CachedResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (w *CachedResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *CachedResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)

	return w.ResponseWriter.Write(b)
}

func CacheMiddleware(cfg *CacheConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				next.ServeHTTP(w, r)

				return
			}

			key := cfg.Cache.GenerateKey(cfg.KeyPrefix, r.URL.Path, r.URL.RawQuery)

			if cached, ok := getCachedResponse(r, cfg, key); ok {
				writeCachedResponse(w, cached)

				return
			}

			w.Header().Set("X-Cache", "MISS")

			ww := &CachedResponseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
				body:           bytes.NewBuffer(nil),
			}

			next.ServeHTTP(ww, r)

			storeResponse(r, cfg, key, ww)
		})
	}
}

func getCachedResponse(r *http.Request, cfg *CacheConfig, key string) (CachedResponse, bool) {
	res, err := cfg.Cache.Get(r.Context(), key)
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			// TODO: log cache read error
		}

		return CachedResponse{}, false
	}

	var cached CachedResponse
	if err = json.Unmarshal(res, &cached); err != nil {
		return CachedResponse{}, false
	}

	return cached, true
}

func writeCachedResponse(w http.ResponseWriter, cached CachedResponse) {
	header := w.Header()
	for k, v := range cached.Headers {
		header.Set(k, v)
	}
	header.Set("X-Cache", "HIT")

	w.WriteHeader(cached.StatusCode)
	w.Write(cached.Body)
}

func storeResponse(r *http.Request, cfg *CacheConfig, key string, ww *CachedResponseWriter) {
	if ww.statusCode < 200 || ww.statusCode > 299 {
		return
	}

	headers := make(map[string]string, len(ww.Header()))
	for k, v := range ww.Header() {
		if _, skip := uncacheableHeaders[k]; skip {
			continue
		}

		headers[k] = strings.Join(v, ",")
	}

	cached := CachedResponse{
		StatusCode: ww.statusCode,
		Headers:    headers,
		Body:       ww.body.Bytes(),
	}

	jsonBody, err := json.Marshal(cached)
	if err != nil {
		return
	}

	cfg.Cache.Set(r.Context(), key, jsonBody, cfg.TTL)
}
