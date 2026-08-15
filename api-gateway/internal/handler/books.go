package handler

import (
	"net/http"

	"bookshelf/api-gateway/internal/proxy"

	"github.com/go-chi/chi/v5"
)

type BooksHandler struct {
	proxy *proxy.ServiceProxy
}

func NewBooksHandler(p *proxy.ServiceProxy) *BooksHandler {
	return &BooksHandler{proxy: p}
}

func (h *BooksHandler) ListBooks(w http.ResponseWriter, r *http.Request) {
	h.proxy.ProxyBooksPath(w, r, "/api/v1/books")
}

func (h *BooksHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	h.proxy.ProxyBooksPath(w, r, "/api/v1/books")
}

func (h *BooksHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "bookId")
	h.proxy.ProxyBooksPath(w, r, "/api/v1/books/"+bookID)
}

func (h *BooksHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "bookId")
	h.proxy.ProxyBooksPath(w, r, "/api/v1/books/"+bookID)
}

func (h *BooksHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "bookId")
	h.proxy.ProxyBooksPath(w, r, "/api/v1/books/"+bookID)
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
	h.proxy.ProxyBooksPath(w, r, "/api/v1/books/"+bookID+"/cover")
}

func (h *BooksHandler) DeleteBookCover(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "bookId")
	h.proxy.ProxyBooksPath(w, r, "/api/v1/books/"+bookID+"/cover")
}
