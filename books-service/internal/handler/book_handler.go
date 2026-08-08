package handler

import (
	"bookshelf/books-service/internal/domain"
	"bookshelf/books-service/internal/service"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) ListBooks(w http.ResponseWriter, r *http.Request) {
	filter := extractFilter(r)

	resp, err := h.services.BookService.List(r.Context(), filter)
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, err.Error())

		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetBook(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "bookId")

	resp, err := h.services.BookService.GetByID(r.Context(), bookID)
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

func (h *Handler) CreateBook(w http.ResponseWriter, r *http.Request) {
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

	resp, err := h.services.BookService.Create(r.Context(), userID, req)
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, err.Error())

		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *Handler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r.Context())
	bookID := chi.URLParam(r, "bookId")

	var req domain.UpdateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, r, http.StatusBadRequest, "invalid request body")

		return
	}
	defer r.Body.Close()

	book, err := h.services.BookService.Update(r.Context(), userID, bookID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrBookNotFound):
			writeError(w, r, http.StatusNotFound, "book not found")
		case errors.Is(err, service.ErrNotBookOwner):
			writeError(w, r, http.StatusForbidden, "not book owner")
		default:
			writeError(w, r, http.StatusInternalServerError, err.Error())
		}
	}

	writeJSON(w, http.StatusOK, book)
}

func (h *Handler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r.Context())
	bookID := chi.URLParam(r, "bookId")

	if err := h.services.BookService.Delete(r.Context(), userID, bookID); err != nil {
		switch {
		case errors.Is(err, service.ErrBookNotFound):
			writeError(w, r, http.StatusNotFound, "book not found")
		case errors.Is(err, service.ErrNotBookOwner):
			writeError(w, r, http.StatusForbidden, "not book owner")
		default:
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
