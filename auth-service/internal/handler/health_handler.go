package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
)

type HealthResponse struct {
	Status    string           `json:"status"` // "ok", "unhealthy"
	Service   string           `json:"service"`
	Version   string           `json:"version"`
	Checks    map[string]Check `json:"checks"`
	Timestamp string           `json:"timestamp"`
}

type Check struct {
	Status   string `json:"status"`   // "ok", "error"
	Duration string `json:"duration"` // "2ms"
	Error    string `json:"error,omitempty"`
}

type HealthHandler struct {
	db *sqlx.DB
}

func NewHealthHandler(db *sqlx.DB) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := HealthResponse{
		Status:    "ok",
		Version:   "1.0.0",
		Service:   "auth-service",
		Timestamp: time.Now().Format(time.RFC3339),
		Checks:    map[string]Check{},
	}

	resp.Checks["database"] = h.checkDatabase(r.Context())

	for _, v := range resp.Checks {
		if v.Status == "error" {
			resp.Status = "unhealthy"

			break
		}
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *HealthHandler) checkDatabase(ctx context.Context) Check {
	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	now := time.Now()
	err := h.db.PingContext(pingCtx)
	after := time.Since(now)

	check := Check{
		Status:   "ok",
		Duration: after.String(),
	}

	if err != nil {
		check.Error = err.Error()
		check.Status = "error"
	}

	return check
}
