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
	logger  *zap.Logger
}

func NewImageHandler(storage *storage.MinIOStorage, repo *repository.CoverRepository, logger *zap.Logger) *ImageHandler {
	return &ImageHandler{
		storage: storage,
		repo:    repo,
		logger:  logger,
	}
}

func (h *ImageHandler) HandleImageCompress(body []byte) error {
	ctx := context.Background()
	var msg ImageCompressMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		h.logger.Error("failed to unmarshal image compress message", zap.ByteString("body", body), zap.Error(err))

		return err
	}

	h.logger.Info("processing image compress message", zap.String("book_id", msg.BookID), zap.String("cover_id", msg.CoverID))

	fileBytes, err := h.storage.GetFile(ctx, msg.OriginalPath)
	if err != nil {
		h.logger.Error("failed to get file", zap.String("original_path", msg.OriginalPath), zap.Error(err))

		return fmt.Errorf("%w, failed to get file: %w", queue.ErrTemporary, err)
	}

	img, err := imaging.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		h.logger.Error("failed to decode image", zap.String("original_path", msg.OriginalPath), zap.Error(err))
		if err = h.repo.UpdateStatus(ctx, msg.CoverID, "failed", "", "", err.Error()); err != nil {
			h.logger.Error("failed to update status", zap.String("cover_id", msg.CoverID), zap.Error(err))
		}

		return nil
	}

	coverKey := fmt.Sprintf("covers/%s/cover.jpg", msg.BookID)
	if err = h.processImage(ctx, img, CoverWidth, CoverHeight, JPEGQuality, coverKey); err != nil {
		h.logger.Error("failed to process cover image", zap.String("original_path", msg.OriginalPath), zap.Error(err))
		if err = h.repo.UpdateStatus(ctx, msg.CoverID, "failed", "", "", err.Error()); err != nil {
			h.logger.Error("failed to update status", zap.String("cover_id", msg.CoverID), zap.Error(err))
		}

		return nil
	}

	thumbKey := fmt.Sprintf("covers/%s/thumb.jpg", msg.BookID)
	if err = h.processImage(ctx, img, ThumbWidth, ThumbHeight, JPEGQuality, thumbKey); err != nil {
		h.logger.Error("failed to process thumbnail image", zap.String("original_path", msg.OriginalPath), zap.Error(err))
		if err = h.repo.UpdateStatus(ctx, msg.CoverID, "failed", "", "", err.Error()); err != nil {
			h.logger.Error("failed to update status", zap.String("cover_id", msg.CoverID), zap.Error(err))
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
		h.logger.Error("failed to update status", zap.String("cover_id", msg.CoverID), zap.Error(err))
	}

	if err = h.repo.UpdateBookCover(
		ctx,
		msg.BookID,
		"ready",
		h.storage.GetFileURL(coverKey),
		h.storage.GetFileURL(thumbKey),
	); err != nil {
		h.logger.Error("failed to update book cover", zap.String("book_id", msg.BookID), zap.Error(err))
	}

	h.logger.Info("image compress finished", zap.String("book_id", msg.BookID), zap.String("cover_id", msg.CoverID))

	return nil
}

func (h *ImageHandler) processImage(ctx context.Context, img image.Image, width, height, quality int, objectName string) error {
	data := imaging.Fill(img, width, height, imaging.Center, imaging.Lanczos)

	var buf bytes.Buffer
	if err := imaging.Encode(&buf, data, imaging.JPEG, imaging.JPEGQuality(quality)); err != nil {
		h.logger.Error("failed to encode image", zap.Error(err))

		return err
	}

	if err := h.storage.UploadFile(ctx, objectName, buf.Bytes(), "image/jpg"); err != nil {
		h.logger.Error("failed to upload image", zap.String("object_name", objectName), zap.Error(err))

		return err
	}

	return nil
}
