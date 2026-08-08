package handler

import (
	"bookshelf/auth-service/internal/domain"
	"bookshelf/auth-service/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	AuthHandler     *AuthHandler
	InternalHandler *InternalHandler
}

func NewHandler(userService *service.UserService, jwtSecret string) *Handler {
	return &Handler{
		AuthHandler:     NewAuthHandler(userService, jwtSecret),
		InternalHandler: NewInternalHandler(userService),
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/health", h.Health)
	r.Get("/ready", h.Health)

	h.AuthHandler.RegisterRoutes(r)
	h.InternalHandler.RegisterRoutes(r)
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"status":"ok", "service":"auth-service", "timestamp":%d}`, time.Now().Unix())))
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

func writeError(w http.ResponseWriter, r *http.Request, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := domain.ErrorResponse{
		Error: domain.ErrorData{
			Code:    status,
			Message: message,
		},
	}

	if reqID, ok := r.Context().Value(requestIDKey).(string); ok {
		resp.Error.RequestID = reqID
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, _ = w.Write(bytes)
}
