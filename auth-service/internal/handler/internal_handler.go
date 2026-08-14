package handler

import (
	"bookshelf/auth-service/internal/domain"
	"bookshelf/auth-service/internal/service"
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type InternalHandler struct {
	userService *service.UserService
	logger      *zap.Logger
}

func NewInternalHandler(svc *service.UserService, logger *zap.Logger) *InternalHandler {
	return &InternalHandler{
		userService: svc,
		logger:      logger,
	}
}

func (h *InternalHandler) VerifyToken(w http.ResponseWriter, r *http.Request) {
	var req domain.VerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, r, http.StatusBadRequest, "invalid request body")

		return
	}
	defer r.Body.Close()

	res, err := h.userService.ValidateToken(req.Token)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, &domain.VerifyResponse{Error: err.Error()})

		return
	}

	writeJSON(w, http.StatusOK, &domain.VerifyResponse{
		UserID:    res.UserID,
		Valid:     true,
		ExpiresAt: res.ExpiresAt.Format(time.RFC3339),
	})
}

func (h *InternalHandler) GetUsersByIDs(w http.ResponseWriter, r *http.Request) {
	var req domain.GetUsersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, r, http.StatusBadRequest, "invalid request body")

		return
	}
	defer r.Body.Close()

	users, err := h.userService.GetUsersByIDs(r.Context(), req.IDs)
	if err != nil {
		h.logger.Error("get users by ids failed", zap.Error(err))
		writeError(w, r, http.StatusInternalServerError, err.Error())

		return
	}

	writeJSON(w, http.StatusOK, users)
}
