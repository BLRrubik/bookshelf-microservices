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

type BookHandler struct {
	srv    *service.BookService
	logger *zap.Logger
}

func NewBookHandler(bookService *service.BookService, logger *zap.Logger) *BookHandler {
	return &BookHandler{
		srv:    bookService,
		logger: logger,
	}
}

func (h *BookHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := extractFilter(r)

	books, count, err := h.srv.List(r.Context(), filter)
	if err != nil {
		h.logger.Error("list books failed", zap.Error(err))
		writeError(w, r, http.StatusInternalServerError, err.Error())

		return
	}

	booksResp := make([]domain.BookResponse, len(books))
	for i, book := range books {
		booksResp[i] = book.ToResponse()
	}

	writeJSON(w, http.StatusOK, domain.BookListResponse{
		Data: booksResp,
		Pagination: domain.Pagination{
			Total:      count,
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalPages: (count + filter.Limit - 1) / filter.Limit,
		},
	})
}

func (h *BookHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "id")

	resp, err := h.srv.GetByID(r.Context(), bookID)
	if err != nil {
		if errors.Is(err, service.ErrBookNotFound) {
			writeError(w, r, http.StatusNotFound, "book not found")

			return
		}

		h.logger.Error("get book failed", zap.String("book_id", bookID), zap.Error(err))
		writeError(w, r, http.StatusInternalServerError, err.Error())

		return
	}

	writeJSON(w, http.StatusOK, resp.ToResponse())
}

func (h *BookHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r.Context())

	var req domain.CreateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, r, http.StatusBadRequest, "invalid request body")

		return
	}
	defer r.Body.Close()

	if err := req.Validate(); err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error())

		return
	}

	resp, err := h.srv.Create(r.Context(), userID, req)
	if err != nil {
		h.logger.Error("create book failed", zap.String("user_id", userID), zap.Error(err))
		writeError(w, r, http.StatusInternalServerError, err.Error())

		return
	}

	writeJSON(w, http.StatusCreated, resp.ToResponse())
}

func (h *BookHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r.Context())
	bookID := chi.URLParam(r, "id")

	var req domain.UpdateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, r, http.StatusBadRequest, "invalid request body")

		return
	}
	defer r.Body.Close()

	book, err := h.srv.Update(r.Context(), userID, bookID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrBookNotFound):
			writeError(w, r, http.StatusNotFound, "book not found")
		case errors.Is(err, service.ErrNotBookOwner):
			writeError(w, r, http.StatusForbidden, "not book owner")
		default:
			h.logger.Error("update book failed", zap.String("book_id", bookID), zap.Error(err))
			writeError(w, r, http.StatusInternalServerError, err.Error())
		}

		return
	}

	writeJSON(w, http.StatusOK, book.ToResponse())
}

func (h *BookHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r.Context())
	bookID := chi.URLParam(r, "id")

	if err := h.srv.Delete(r.Context(), userID, bookID); err != nil {
		switch {
		case errors.Is(err, service.ErrBookNotFound):
			writeError(w, r, http.StatusNotFound, "book not found")
		case errors.Is(err, service.ErrNotBookOwner):
			writeError(w, r, http.StatusForbidden, "not book owner")
		default:
			h.logger.Error("delete book failed", zap.String("book_id", bookID), zap.Error(err))
			writeError(w, r, http.StatusInternalServerError, err.Error())
		}

		return
	}

	writeJSON(w, http.StatusNoContent, nil)
}

func extractFilter(r *http.Request) domain.ListParams {
	page, limit := extractPageAndLimit(r)

	search := r.URL.Query().Get("search")
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	return domain.ListParams{
		Page:   page,
		Limit:  limit,
		Search: search,
		Sort:   sort,
		Order:  order,
	}
}
