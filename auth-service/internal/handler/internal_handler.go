package handler

import (
	"bookshelf/auth-service/internal/domain"
	"bookshelf/auth-service/internal/service"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type InternalHandler struct {
	userService *service.UserService
}

func NewInternalHandler(svc *service.UserService) *InternalHandler {
	return &InternalHandler{
		userService: svc,
	}
}

func (h *InternalHandler) RegisterRoutes(r chi.Router) {
	r.Route("/internal/v1", func(r chi.Router) {
		r.Post("/auth/verify", h.VerifyToken)
	})
}

func (h *InternalHandler) VerifyToken(w http.ResponseWriter, r *http.Request) {
	var req domain.VerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, r, http.StatusBadRequest, "invalid request body")

		return
	}
	defer r.Body.Close()

	res := h.userService.ValidateToken(req.Token)
	if res.Error != "" {
		writeError(w, r, http.StatusBadRequest, res.Error)

		return
	}

	writeJSON(w, http.StatusOK, res)
}
