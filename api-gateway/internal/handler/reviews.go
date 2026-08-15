package handler

import (
	"net/http"

	"bookshelf/api-gateway/internal/proxy"

	"github.com/go-chi/chi/v5"
)

type ReviewsHandler struct {
	proxy *proxy.ServiceProxy
}

func NewReviewsHandler(p *proxy.ServiceProxy) *ReviewsHandler {
	return &ReviewsHandler{proxy: p}
}

func (h *ReviewsHandler) ListReviews(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "bookId")
	h.proxy.ProxyBooksPath(w, r, "/api/v1/books/"+bookID+"/reviews")
}

func (h *ReviewsHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "bookId")
	h.proxy.ProxyBooksPath(w, r, "/api/v1/books/"+bookID+"/reviews")
}

func (h *ReviewsHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	reviewID := chi.URLParam(r, "reviewId")
	h.proxy.ProxyBooksPath(w, r, "/api/v1/reviews/"+reviewID)
}

func (h *ReviewsHandler) UpdateReview(w http.ResponseWriter, r *http.Request) {
	reviewID := chi.URLParam(r, "reviewId")
	h.proxy.ProxyBooksPath(w, r, "/api/v1/reviews/"+reviewID)
}

func (h *ReviewsHandler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	reviewID := chi.URLParam(r, "reviewId")
	h.proxy.ProxyBooksPath(w, r, "/api/v1/reviews/"+reviewID)
}
