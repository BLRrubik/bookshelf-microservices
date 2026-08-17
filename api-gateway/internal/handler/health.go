package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"bookshelf/api-gateway/internal/cache"
	"bookshelf/api-gateway/internal/domain"
	"bookshelf/api-gateway/internal/proxy"
)

const readyCheckTimeout = 5 * time.Second

// HealthHandler answers orchestrator questions about service state
// (liveness/readiness), rather than serving business requests.
type HealthHandler struct {
	cache   *cache.Cache
	proxy   *proxy.ServiceProxy
	version string
}

func NewHealthHandler(c *cache.Cache, p *proxy.ServiceProxy, version string) *HealthHandler {
	return &HealthHandler{
		cache:   c,
		proxy:   p,
		version: version,
	}
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	writeHealthResponse(w, http.StatusOK, domain.HealthResponse{
		Status:    "ok",
		Version:   h.version,
		Timestamp: time.Now().UTC(),
		Services:  map[string]string{},
	})
}

func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), readyCheckTimeout)
	defer cancel()

	checks := map[string]string{}
	status := "ok"

	if err := h.cache.Ping(ctx); err != nil {
		checks["redis"] = "unhealthy: " + err.Error()
		status = "degraded"
	} else {
		checks["redis"] = "ok"
	}

	if err := h.proxy.CheckAuthService(ctx); err != nil {
		checks["auth_service"] = "unhealthy: " + err.Error()
		status = "degraded"
	} else {
		checks["auth_service"] = "ok"
	}

	if err := h.proxy.CheckBooksService(ctx); err != nil {
		checks["books_service"] = "unhealthy: " + err.Error()
		status = "degraded"
	} else {
		checks["books_service"] = "ok"
	}

	code := http.StatusOK
	if status != "ok" {
		code = http.StatusServiceUnavailable
	}

	writeHealthResponse(w, code, domain.HealthResponse{
		Status:    status,
		Version:   h.version,
		Timestamp: time.Now().UTC(),
		Services:  checks,
	})
}

func writeHealthResponse(w http.ResponseWriter, code int, body domain.HealthResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}
