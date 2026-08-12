package handler

import (
	"bookshelf/books-service/internal/service"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type CoverHandler struct {
	svc *service.CoverService
}

func NewCoverHandler(svc *service.CoverService) *CoverHandler {
	return &CoverHandler{svc: svc}
}

func (h *CoverHandler) UploadBookCover(w http.ResponseWriter, r *http.Request) {
	body := http.MaxBytesReader(w, r.Body, 5<<20)
	defer body.Close()

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		slog.Error(fmt.Sprintf("failed to parse multipart form: %s", err))
		writeError(w, r, http.StatusBadRequest, "failed to parse multipart form")

		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		slog.Error(fmt.Sprintf("failed to parse form file: %s", err))
		writeError(w, r, http.StatusBadRequest, "failed to parse form file")

		return
	}
	defer file.Close()

	bookID := chi.URLParam(r, "id")
	userID := getUserID(r.Context())

	coverResp, err := h.svc.UploadCover(r.Context(), userID, bookID, file, header.Size, header.Filename)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrBookNotFound):
			writeError(w, r, http.StatusNotFound, "book not found")
		case errors.Is(err, service.ErrForbidden):
			writeError(w, r, http.StatusForbidden, "forbidden")
		default:
			writeError(w, r, http.StatusInternalServerError, "internal server error")
		}

		return
	}

	writeJSON(w, http.StatusAccepted, coverResp)
}
