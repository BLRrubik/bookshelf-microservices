package handler

import (
	"bookshelf/books-service/internal/domain"
	"bookshelf/books-service/internal/service"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ReviewHandler struct {
	svc *service.ReviewService
}

func NewReviewHandler(svc *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{
		svc: svc,
	}
}

func (h *ReviewHandler) List(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "book_id")
	page, limit := extractPageAndLimit(r)

	reviews, count, err := h.svc.ListByBook(r.Context(), bookID, page, limit)
	if err != nil {
		if errors.Is(err, service.ErrBookNotFound) {
			writeError(w, r, http.StatusNotFound, "book not found")

			return
		}

		writeError(w, r, http.StatusInternalServerError, err.Error())

		return
	}

	reviewsResp := make([]domain.ReviewResponse, len(reviews))
	for i, review := range reviews {
		reviewsResp[i] = review.ToResponse()
	}

	writeJSON(w, http.StatusOK, domain.ReviewListResponse{
		Data: reviewsResp,
		Pagination: domain.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      count,
			TotalPages: (count + limit - 1) / limit,
		},
	})
}

func (h *ReviewHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	reviewID := chi.URLParam(r, "id")

	review, err := h.svc.GetByID(r.Context(), reviewID)
	if err != nil {
		if errors.Is(err, service.ErrReviewNotFound) {
			writeError(w, r, http.StatusNotFound, "review not found")

			return
		}

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
			writeError(w, r, http.StatusInternalServerError, err.Error())
		}

		return
	}

	writeJSON(w, http.StatusNoContent, nil)
}
