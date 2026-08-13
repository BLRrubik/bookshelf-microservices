package handler

import (
	"bookshelf/worker-service/internal/queue"
	"bookshelf/worker-service/internal/storage"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"log/slog"

	"github.com/disintegration/imaging"
	"github.com/jmoiron/sqlx"
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
	db      *sqlx.DB
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

		return err
	}

	if err = h.processImage(ctx, img, CoverWidth, CoverHeight, JPEGQuality, fmt.Sprintf("/covers/%s/cover.jpg", msg.BookID)); err != nil {
		slog.Error("Failed to process image", zap.String("original_path", msg.OriginalPath), zap.Error(err))

		return err
	}

	if err = h.processImage(ctx, img, ThumbWidth, ThumbHeight, JPEGQuality, fmt.Sprintf("/covers/%s/thumb.jpg", msg.BookID)); err != nil {
		slog.Error("Failed to process image", zap.String("original_path", msg.OriginalPath), zap.Error(err))

		return err
	}

	return nil
}

func (h *ImageHandler) processImage(ctx context.Context, img image.Image, width, height, quality int, objectName string) error {
	data := imaging.Fill(img, ThumbWidth, ThumbHeight, imaging.Center, imaging.Lanczos)

	var buf bytes.Buffer
	if err := imaging.Encode(&buf, data, imaging.JPEG, imaging.JPEGQuality(JPEGQuality)); err != nil {
		slog.Error("Failed to encode image", zap.Error(err))

		return err
	}

	if err := h.storage.UploadFile(ctx, objectName, buf.Bytes(), "image/jpg"); err != nil {
		slog.Error("Failed to upload image", zap.String("object_name", objectName), zap.Error(err))

		return err
	}

	return nil
}
