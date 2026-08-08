package handler

import (
	"bookshelf/books-service/internal/client"
	"bookshelf/books-service/internal/domain"
	"bookshelf/books-service/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	bookHandler   *BookHandler
	reviewHandler *ReviewHandler
	authClient    *client.AuthClient
}

func NewHandler(service *service.Service, authClient *client.AuthClient) *Handler {
	return &Handler{
		bookHandler:   NewBookHandler(service.BookService),
		reviewHandler: NewReviewHandler(service.ReviewService),
		authClient:    authClient,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/health", h.Health)
	r.Get("/ready", h.Health)

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
		})

	})
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"status":"ok", "service":"books-service", "timestamp":%d}`, time.Now().Unix())))
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
