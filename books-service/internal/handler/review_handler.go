package handler

import (
	"bookshelf/books-service/internal/domain"
	"bookshelf/books-service/internal/service"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type ReviewHandler struct {
	svc    *service.ReviewService
	logger *zap.Logger
}

func NewReviewHandler(svc *service.ReviewService, logger *zap.Logger) *ReviewHandler {
	return &ReviewHandler{
		svc:    svc,
		logger: logger,
	}
}

func (h *ReviewHandler) List(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "book_id")
	page, limit := extractPageAndLimit(r)

	reviews, err := h.svc.ListByBook(r.Context(), bookID, page, limit)
	if err != nil {
		if errors.Is(err, service.ErrBookNotFound) {
			writeError(w, r, http.StatusNotFound, "book not found")

			return
		}

		h.logger.Error("list reviews failed", zap.String("book_id", bookID), zap.Error(err))
		writeError(w, r, http.StatusInternalServerError, err.Error())

		return
	}

	writeJSON(w, http.StatusOK, reviews)
}

func (h *ReviewHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	reviewID := chi.URLParam(r, "id")

	review, err := h.svc.GetByID(r.Context(), reviewID)
	if err != nil {
		if errors.Is(err, service.ErrReviewNotFound) {
			writeError(w, r, http.StatusNotFound, "review not found")

			return
		}

		h.logger.Error("get review failed", zap.String("review_id", reviewID), zap.Error(err))
		writeError(w, r, http.StatusInternalServerError, err.Error())

		return
	}

	writeJSON(w, http.StatusOK, review.ToResponse())
}

func (h *ReviewHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r.Context())
	bookID := chi.URLParam(r, "book_id")

	var req domain.CreateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error())

		return
	}
	defer r.Body.Close()

	if err := req.Validate(); err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error())

		return
	}

	review, err := h.svc.Create(r.Context(), userID, bookID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrBookNotFound),
			errors.Is(err, service.ErrAlreadyReviewed):
			writeError(w, r, http.StatusConflict, "cannot create review")
		default:
			h.logger.Error("create review failed", zap.String("book_id", bookID), zap.String("user_id", userID), zap.Error(err))
			writeError(w, r, http.StatusInternalServerError, err.Error())
		}

		return
	}

	writeJSON(w, http.StatusCreated, review)
}

func (h *ReviewHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r.Context())
	reviewID := chi.URLParam(r, "id")

	var req domain.UpdateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error())

		return
	}
	defer r.Body.Close()

	review, err := h.svc.Update(r.Context(), reviewID, userID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrReviewNotFound):
			writeError(w, r, http.StatusNotFound, "review not found")
		case errors.Is(err, service.ErrNotReviewOwner):
			writeError(w, r, http.StatusForbidden, "not review owner")
		default:
			h.logger.Error("update review failed", zap.String("review_id", reviewID), zap.Error(err))
			writeError(w, r, http.StatusInternalServerError, err.Error())
		}

		return
	}

	writeJSON(w, http.StatusOK, review.ToResponse())
}

func (h *ReviewHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r.Context())
	reviewID := chi.URLParam(r, "id")

	if err := h.svc.Delete(r.Context(), reviewID, userID); err != nil {
		switch {
		case errors.Is(err, service.ErrReviewNotFound):
			writeError(w, r, http.StatusNotFound, "review not found")
		case errors.Is(err, service.ErrNotReviewOwner):
			writeError(w, r, http.StatusForbidden, "not review owner")
		default:
			h.logger.Error("delete review failed", zap.String("review_id", reviewID), zap.Error(err))
			writeError(w, r, http.StatusInternalServerError, err.Error())
		}

		return
	}

	writeJSON(w, http.StatusNoContent, nil)
}
