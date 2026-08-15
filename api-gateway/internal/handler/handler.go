package handler

import (
	"net/http"

	"bookshelf/api-gateway/internal/proxy"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	authHandler    *AuthHandler
	booksHandler   *BooksHandler
	reviewsHandler *ReviewsHandler
}

func NewHandler(p *proxy.ServiceProxy) *Handler {
	return &Handler{
		authHandler:    NewAuthHandler(p),
		booksHandler:   NewBooksHandler(p),
		reviewsHandler: NewReviewsHandler(p),
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", h.authHandler.Register)
		r.Post("/login", h.authHandler.Login)
		r.Post("/refresh", h.authHandler.Refresh)
		r.Post("/logout", h.authHandler.Logout)
	})

	r.Route("/api/v1/users", func(r chi.Router) {
		r.Get("/me", h.authHandler.GetCurrentUser)
		r.Put("/me", h.authHandler.UpdateCurrentUser)
		r.Get("/{userId}", h.authHandler.GetUser)
	})

	r.Route("/api/v1/books", func(r chi.Router) {
		r.Get("/", h.booksHandler.ListBooks)
		r.Post("/", h.booksHandler.CreateBook)
		r.Get("/{bookId}", h.booksHandler.GetBook)
		r.Put("/{bookId}", h.booksHandler.UpdateBook)
		r.Delete("/{bookId}", h.booksHandler.DeleteBook)
		r.Get("/{bookId}/cover", h.booksHandler.GetBookCover)
		r.Get("/{bookId}/cover/status", h.booksHandler.GetBookCoverStatus)
		r.Post("/{bookId}/cover", h.booksHandler.UploadBookCover)
		r.Delete("/{bookId}/cover", h.booksHandler.DeleteBookCover)
		r.Get("/{bookId}/reviews", h.reviewsHandler.ListReviews)
		r.Post("/{bookId}/reviews", h.reviewsHandler.CreateReview)
	})

	r.Route("/api/v1/reviews", func(r chi.Router) {
		r.Get("/{reviewId}", h.reviewsHandler.GetReview)
		r.Put("/{reviewId}", h.reviewsHandler.UpdateReview)
		r.Delete("/{reviewId}", h.reviewsHandler.DeleteReview)
	})
}
