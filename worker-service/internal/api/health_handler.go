package api

import (
	"bookshelf/worker-service/internal/queue"
	"bookshelf/worker-service/internal/storage"
	"context"
	"encoding/json"
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

type ReadyResponse struct {
	Ready     bool             `json:"ready"`
	Service   string           `json:"service"`
	Checks    map[string]Check `json:"checks"`
	Timestamp string           `json:"timestamp"`
}

type Check struct {
	Status   string `json:"status"`   // "ok", "error"
	Duration string `json:"duration"` // "2ms"
	Error    string `json:"error,omitempty"`
}

type HealthHandler struct {
	db             *sqlx.DB
	rabbitMQClient *queue.Consumer
	minioClient    *storage.MinIOStorage
}

func NewHealthHandler(
	db *sqlx.DB,
	rabbitMQClient *queue.Consumer,
	minioClient *storage.MinIOStorage,
) *HealthHandler {
	return &HealthHandler{
		db:             db,
		rabbitMQClient: rabbitMQClient,
		minioClient:    minioClient,
	}
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := HealthResponse{
		Status:    "ok",
		Version:   "1.0.0",
		Service:   "books-service",
		Timestamp: time.Now().Format(time.RFC3339),
		Checks:    map[string]Check{},
	}

	resp.Checks["database"] = h.checkDatabase(r.Context())
	resp.Checks["rabbitMQ"] = h.checkRabbitMQ(r.Context())
	resp.Checks["minio"] = h.checkMinio(r.Context())

	for _, v := range resp.Checks {
		if v.Status == "error" {
			resp.Status = "unhealthy"

			break
		}
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := ReadyResponse{
		Ready:     true,
		Service:   "books-service",
		Timestamp: time.Now().Format(time.RFC3339),
		Checks:    map[string]Check{},
	}

	resp.Checks["database"] = h.checkDatabase(r.Context())
	resp.Checks["rabbitMQ"] = h.checkRabbitMQ(r.Context())
	resp.Checks["minio"] = h.checkMinio(r.Context())

	for _, v := range resp.Checks {
		if v.Status == "error" {
			resp.Ready = false

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

func (h *HealthHandler) checkRabbitMQ(ctx context.Context) Check {
	now := time.Now()
	err := h.rabbitMQClient.HealthCheck()
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

func (h *HealthHandler) checkMinio(ctx context.Context) Check {
	now := time.Now()
	err := h.minioClient.HealthCheck(ctx)
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

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	bytes, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, _ = w.Write(bytes)
}
