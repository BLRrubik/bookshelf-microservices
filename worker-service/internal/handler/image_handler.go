package handler

import (
	"bookshelf/worker-service/internal/queue"
	"bookshelf/worker-service/internal/repository"
	"bookshelf/worker-service/internal/storage"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"log/slog"

	"github.com/disintegration/imaging"
	"go.uber.org/zap"
)

const (
	CoverWidth  = 400
	CoverHeight = 600
	ThumbWidth  = 100
	ThumbHeight = 150
	JPEGQuality = 85
)

type ImageCompressMessage struct {
	BookID       string `json:"book_id"`
	CoverID      string `json:"cover_id"`
	OriginalPath string `json:"original_path"`
}

type ImageHandler struct {
	storage *storage.MinIOStorage
	repo    *repository.CoverRepository
}

func NewImageHandler(storage *storage.MinIOStorage, repo *repository.CoverRepository) *ImageHandler {
	return &ImageHandler{
		storage: storage,
		repo:    repo,
	}
}

func (h *ImageHandler) HandleImageCompress(body []byte) error {
	ctx := context.Background()
	var msg ImageCompressMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		slog.Error("Failed to unmarshal image compress message", zap.String("body", string(body)), zap.Error(err))

		return err
	}

	fileBytes, err := h.storage.GetFile(ctx, msg.OriginalPath)
	if err != nil {
		slog.Error("Failed to get file", zap.String("original_path", msg.OriginalPath), zap.Error(err))

		return fmt.Errorf("%w, failed to get file: %w", queue.ErrTemporary, err)
	}

	img, err := imaging.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		slog.Error("Failed to decode image", zap.String("original_path", msg.OriginalPath), zap.Error(err))
		if err = h.repo.UpdateStatus(ctx, msg.CoverID, "failed", "", "", err.Error()); err != nil {
			slog.Error("Failed to update status", zap.String("original_path", msg.OriginalPath), zap.Error(err))
		}

		return nil
	}

	coverKey := fmt.Sprintf("covers/%s/cover.jpg", msg.BookID)
	if err = h.processImage(ctx, img, CoverWidth, CoverHeight, JPEGQuality, coverKey); err != nil {
		slog.Error("Failed to process image", zap.String("original_path", msg.OriginalPath), zap.Error(err))
		if err = h.repo.UpdateStatus(ctx, msg.CoverID, "failed", "", "", err.Error()); err != nil {
			slog.Error("Failed to update status", zap.String("original_path", msg.OriginalPath), zap.Error(err))
		}

		return nil
	}

	thumbKey := fmt.Sprintf("covers/%s/thumb.jpg", msg.BookID)
	if err = h.processImage(ctx, img, ThumbWidth, ThumbHeight, JPEGQuality, thumbKey); err != nil {
		slog.Error("Failed to process image", zap.String("original_path", msg.OriginalPath), zap.Error(err))
		if err = h.repo.UpdateStatus(ctx, msg.CoverID, "failed", "", "", err.Error()); err != nil {
			slog.Error("Failed to update status", zap.String("original_path", msg.OriginalPath), zap.Error(err))
		}

		return nil
	}

	if err = h.repo.UpdateStatus(
		ctx,
		msg.CoverID,
		"ready",
		coverKey,
		thumbKey,
		"",
	); err != nil {
		slog.Error("Failed to update status", zap.String("original_path", msg.OriginalPath), zap.Error(err))
	}

	if err = h.repo.UpdateBookCover(
		ctx,
		msg.BookID,
		"ready",
		h.storage.GetFileURL(coverKey),
		h.storage.GetFileURL(thumbKey),
	); err != nil {
		slog.Error("Failed to update status", zap.String("original_path", msg.OriginalPath), zap.Error(err))
	}

	return nil
}

func (h *ImageHandler) processImage(ctx context.Context, img image.Image, width, height, quality int, objectName string) error {
	data := imaging.Fill(img, width, height, imaging.Center, imaging.Lanczos)

	var buf bytes.Buffer
	if err := imaging.Encode(&buf, data, imaging.JPEG, imaging.JPEGQuality(quality)); err != nil {
		slog.Error("Failed to encode image", zap.Error(err))

		return err
	}

	if err := h.storage.UploadFile(ctx, objectName, buf.Bytes(), "image/jpg"); err != nil {
		slog.Error("Failed to upload image", zap.String("object_name", objectName), zap.Error(err))

		return err
	}

	return nil
}
