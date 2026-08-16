package handler

import (
	"net/http"
	"time"

	"bookshelf/api-gateway/internal/cache"
	"bookshelf/api-gateway/internal/middleware"
	"bookshelf/api-gateway/internal/proxy"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	authHandler      *AuthHandler
	booksHandler     *BooksHandler
	reviewsHandler   *ReviewsHandler
	dashboardHandler *DashboardHandler

	generalCache    func(http.Handler) http.Handler
	booksListCache  func(http.Handler) http.Handler
	bookDetailCache func(http.Handler) http.Handler
}

func NewHandler(p *proxy.ServiceProxy, c *cache.Cache, ttl time.Duration) *Handler {
	return &Handler{
		authHandler:      NewAuthHandler(p),
		booksHandler:     NewBooksHandler(p, c),
		reviewsHandler:   NewReviewsHandler(p),
		dashboardHandler: NewDashboardHandler(p, c, ttl),

		generalCache: middleware.CacheMiddleware(&middleware.CacheConfig{
			Cache:     c,
			KeyPrefix: "cache",
			TTL:       ttl,
		}),
		booksListCache: middleware.CacheMiddleware(&middleware.CacheConfig{
			Cache:     c,
			KeyPrefix: "books",
			TTL:       ttl,
		}),
		bookDetailCache: middleware.CacheMiddleware(&middleware.CacheConfig{
			Cache:     c,
			KeyPrefix: "book",
			TTL:       ttl,
		}),
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Get("/api/v1/dashboard", h.dashboardHandler.GetDashboard)

	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", h.authHandler.Register)
		r.Post("/login", h.authHandler.Login)
		r.Post("/refresh", h.authHandler.Refresh)
		r.Post("/logout", h.authHandler.Logout)
	})

	r.Route("/api/v1/users", func(r chi.Router) {
		r.With(h.generalCache).Get("/me", h.authHandler.GetCurrentUser)
		r.Put("/me", h.authHandler.UpdateCurrentUser)
		r.With(h.generalCache).Get("/{userId}", h.authHandler.GetUser)
	})

	r.Route("/api/v1/books", func(r chi.Router) {
		r.With(h.booksListCache).Get("/", h.booksHandler.ListBooks)
		r.Post("/", h.booksHandler.CreateBook)
		r.With(h.bookDetailCache).Get("/{bookId}", h.booksHandler.GetBook)
		r.Put("/{bookId}", h.booksHandler.UpdateBook)
		r.Delete("/{bookId}", h.booksHandler.DeleteBook)
		r.With(h.bookDetailCache).Get("/{bookId}/cover", h.booksHandler.GetBookCover)
		r.Get("/{bookId}/cover/status", h.booksHandler.GetBookCoverStatus)
		r.Post("/{bookId}/cover", h.booksHandler.UploadBookCover)
		r.Delete("/{bookId}/cover", h.booksHandler.DeleteBookCover)
		r.With(h.generalCache).Get("/{bookId}/reviews", h.reviewsHandler.ListReviews)
		r.Post("/{bookId}/reviews", h.reviewsHandler.CreateReview)
	})

	r.Route("/api/v1/reviews", func(r chi.Router) {
		r.With(h.generalCache).Get("/{reviewId}", h.reviewsHandler.GetReview)
		r.Put("/{reviewId}", h.reviewsHandler.UpdateReview)
		r.Delete("/{reviewId}", h.reviewsHandler.DeleteReview)
	})
}
