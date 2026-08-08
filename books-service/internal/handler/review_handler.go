package handler

import (
	"bookshelf/book-service/internal/domain"
	"bookshelf/book-service/internal/service"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) ListBookReviews(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "bookId")

	page, limit := extractPageAndLimit(r)

	resp, err := h.services.ReviewService.ListByBookID(r.Context(), bookID, page, limit)
	if err != nil {
		if errors.Is(err, service.ErrBookNotFound) {
			writeError(w, r, http.StatusNotFound, "book not found")

			return
		}

		writeError(w, r, http.StatusInternalServerError, err.Error())

		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetReview(w http.ResponseWriter, r *http.Request) {
	reviewID := chi.URLParam(r, "reviewId")

	resp, err := h.services.ReviewService.GetByID(r.Context(), reviewID)
	if err != nil {
		if errors.Is(err, service.ErrReviewNotFound) {
			writeError(w, r, http.StatusNotFound, "review not found")

			return
		}

		writeError(w, r, http.StatusInternalServerError, err.Error())

		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) CreateReview(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r.Context())
	bookID := chi.URLParam(r, "bookId")

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

	resp, err := h.services.ReviewService.Create(r.Context(), userID, bookID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrBookNotFound),
			errors.Is(err, service.ErrAlreadyReviewed):
			writeError(w, r, http.StatusConflict, "cannot create review")
		default:
			writeError(w, r, http.StatusInternalServerError, err.Error())
		}

		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *Handler) UpdateReview(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r.Context())
	reviewID := chi.URLParam(r, "reviewId")

	var req domain.UpdateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error())

		return
	}
	defer r.Body.Close()

	resp, err := h.services.ReviewService.Update(r.Context(), reviewID, userID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrReviewNotFound):
			writeError(w, r, http.StatusNotFound, "review not found")
		case errors.Is(err, service.ErrNotReviewOwner):
			writeError(w, r, http.StatusForbidden, "not review owner")
		default:
			writeError(w, r, http.StatusInternalServerError, err.Error())
		}

		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r.Context())
	reviewID := chi.URLParam(r, "reviewId")

	if err := h.services.ReviewService.Delete(r.Context(), reviewID, userID); err != nil {
		switch {
		case errors.Is(err, service.ErrReviewNotFound):
			writeError(w, r, http.StatusNotFound, "review not found")
		case errors.Is(err, service.ErrNotReviewOwner):
			writeError(w, r, http.StatusForbidden, "not review owner")
		default:
			writeError(w, r, http.StatusInternalServerError, err.Error())
		}

		return
	}

	writeJSON(w, http.StatusNoContent, nil)
}
