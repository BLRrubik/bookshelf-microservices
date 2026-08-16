package handler

import (
	"bookshelf/api-gateway/internal/cache"
	"bookshelf/api-gateway/internal/proxy"
	"bookshelf/api-gateway/internal/utils"
	"context"
	"net/http"
	"time"
)

type DashboardHandler struct {
	proxy *proxy.ServiceProxy
	cache *cache.Cache
	ttl   time.Duration
}

func NewDashboardHandler(proxy *proxy.ServiceProxy, cache *cache.Cache, ttl time.Duration) *DashboardHandler {
	return &DashboardHandler{proxy: proxy, cache: cache, ttl: ttl}
}

func (h *DashboardHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var userID string
	if token := utils.ExtractBearerToken(r); token != "" {
		userID = h.verifyToken(ctx, token)
	}

	cacheKey := dashboardCacheKey(userID)

	if cached, err := h.cache.Get(ctx, cacheKey); err == nil {
		w.Header().Set("X-Cache", "HIT")
		w.Header().Set("Content-Type", "application/json")
		w.Write(cached)

		return
	}

	w.Header().Set("X-Cache", "MISS")
	w.WriteHeader(http.StatusOK)
}

func dashboardCacheKey(userID string) string {
	if userID == "" {
		return "dashboard:anon"
	}

	if len(userID) > 8 {
		userID = userID[:8]
	}

	return "dashboard:user:" + userID
}

func (h *DashboardHandler) verifyToken(ctx context.Context, token string) string {
	userID, err := h.proxy.VerifyToken(ctx, token)
	if err != nil {
		return ""
	}

	return userID
}
