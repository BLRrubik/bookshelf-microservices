package handler

import (
	"log/slog"
	"net/http"

	"bookshelf/api-gateway/internal/cache"
	"bookshelf/api-gateway/internal/proxy"

	"github.com/go-chi/chi/v5"
)

const (
	booksListPattern  = "books:*"
	bookDetailPattern = "book:*"
)

type BooksHandler struct {
	proxy *proxy.ServiceProxy
	cache *cache.Cache
}

func NewBooksHandler(p *proxy.ServiceProxy, c *cache.Cache) *BooksHandler {
	return &BooksHandler{proxy: p, cache: c}
}

func (h *BooksHandler) ListBooks(w http.ResponseWriter, r *http.Request) {
	h.proxy.ProxyBooksPath(w, r, "/api/v1/books")
}

func (h *BooksHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	status := h.proxy.ProxyBooksPath(w, r, "/api/v1/books")
	h.invalidateOnSuccess(r, status, booksListPattern)
}

func (h *BooksHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "bookId")
	h.proxy.ProxyBooksPath(w, r, "/api/v1/books/"+bookID)
}

func (h *BooksHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "bookId")
	status := h.proxy.ProxyBooksPath(w, r, "/api/v1/books/"+bookID)
	h.invalidateOnSuccess(r, status, bookDetailPattern, booksListPattern)
}

func (h *BooksHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "bookId")
	status := h.proxy.ProxyBooksPath(w, r, "/api/v1/books/"+bookID)
	h.invalidateOnSuccess(r, status, bookDetailPattern, booksListPattern)
}

func (h *BooksHandler) GetBookCover(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "bookId")
	h.proxy.ProxyBooksPath(w, r, "/api/v1/books/"+bookID+"/cover")
}

func (h *BooksHandler) GetBookCoverStatus(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "bookId")
	h.proxy.ProxyBooksPath(w, r, "/api/v1/books/"+bookID+"/cover/status")
}

func (h *BooksHandler) UploadBookCover(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "bookId")
	status := h.proxy.ProxyBooksPath(w, r, "/api/v1/books/"+bookID+"/cover")
	h.invalidateOnSuccess(r, status, bookDetailPattern)
}

func (h *BooksHandler) DeleteBookCover(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "bookId")
	status := h.proxy.ProxyBooksPath(w, r, "/api/v1/books/"+bookID+"/cover")
	h.invalidateOnSuccess(r, status, bookDetailPattern)
}

func (h *BooksHandler) invalidateOnSuccess(r *http.Request, status int, patterns ...string) {
	if status < 200 || status > 299 {
		return
	}

	for _, pattern := range patterns {
		if err := h.cache.DeletePattern(r.Context(), pattern); err != nil {
			slog.ErrorContext(
				r.Context(),
				"invalidate cache",
				"pattern", pattern,
				"error", err,
			)
		}
	}
}
