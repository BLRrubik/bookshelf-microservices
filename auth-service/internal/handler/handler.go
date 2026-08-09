package handler

import (
	"bookshelf/auth-service/internal/domain"
	"bookshelf/auth-service/internal/service"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type Handler struct {
	authHandler     *AuthHandler
	internalHandler *InternalHandler
	healthHandler   *HealthHandler
}

func NewHandler(userService *service.UserService, db *sqlx.DB, jwtSecret string) *Handler {
	return &Handler{
		authHandler:     NewAuthHandler(userService, jwtSecret),
		internalHandler: NewInternalHandler(userService),
		healthHandler:   NewHealthHandler(db),
	}
}

func (h *Handler) RegisterRoutes(r chi.Router, serviceKey string) {
	r.Get("/ready", h.Ready)
	r.Get("/health", h.healthHandler.Health)

	r.Route("/api/v1", func(r chi.Router) {
		// Публичные роуты — без авторизации
		r.Post("/auth/register", h.authHandler.Register)
		r.Post("/auth/login", h.authHandler.Login)

		// Защищённые роуты — с AuthMiddleware
		r.Group(func(r chi.Router) {
			r.Use(h.authHandler.AuthMiddleware)

			r.Get("/users/me", h.authHandler.GetMe)
			r.Put("/users/me", h.authHandler.UpdateMe)
		})
	})

	r.Route("/internal/v1", func(r chi.Router) {
		r.Use(ServiceKeyMiddleware(serviceKey))
		r.Post("/auth/verify", h.internalHandler.VerifyToken)
		r.Post("/users/batch", h.internalHandler.GetUsersByIDs)
	})
}

func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
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
