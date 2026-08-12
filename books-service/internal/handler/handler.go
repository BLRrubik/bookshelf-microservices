package handler

import (
	"bookshelf/books-service/internal/client"
	"bookshelf/books-service/internal/domain"
	"bookshelf/books-service/internal/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type Handler struct {
	bookHandler   *BookHandler
	reviewHandler *ReviewHandler
	healthHandler *HealthHandler
	coverHandler  *CoverHandler
	authClient    *client.AuthClient
}

func NewHandler(
	service *service.Service,
	db *sqlx.DB,
	authClient *client.AuthClient,
	rabbitMQClient *client.RabbitMQClient,
) *Handler {
	return &Handler{
		bookHandler:   NewBookHandler(service.BookService),
		reviewHandler: NewReviewHandler(service.ReviewService),
		healthHandler: NewHealthHandler(db, authClient, rabbitMQClient),
		coverHandler:  NewCoverHandler(service.CoverService),
		authClient:    authClient,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/health", h.healthHandler.Health)
	r.Get("/ready", h.healthHandler.Ready)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/books/{book_id}/reviews", h.reviewHandler.List)
		r.Get("/reviews/{id} ", h.reviewHandler.GetReview)

		r.Get("/books", h.bookHandler.List)
		r.Get("/books/{id}", h.bookHandler.GetByID)

		r.Group(func(r chi.Router) {
			r.Use(AuthMiddleware(h.authClient))

			r.Post("/books/{book_id}/reviews", h.reviewHandler.Create)
			r.Put("/reviews/{id} ", h.reviewHandler.Update)
			r.Delete("/reviews/{id} ", h.reviewHandler.Delete)

			r.Post("/books", h.bookHandler.Create)
			r.Put("/books/{id}", h.bookHandler.Update)
			r.Delete("/books/{id}", h.bookHandler.Delete)

			r.Post("/books/{id}/cover", h.coverHandler.UploadBookCover)
		})

	})
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
		Code:    status,
		Message: message,
	}

	if reqID, ok := r.Context().Value(requestIDKey).(string); ok {
		resp.RequestID = reqID
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, _ = w.Write(bytes)
}

func extractPageAndLimit(r *http.Request) (int, int) {
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	return page, limit
}
