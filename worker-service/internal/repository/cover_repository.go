package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const (
	updateCoverStateQuery = `
UPDATE covers SET status = $1, cover_path = $2, thumb_path = $3, error = $4 WHERE id = $5;
`
	updateCoverStateCompletedQuery = `
UPDATE covers SET status = $1, cover_path = $2, thumb_path = $3, error = $4, completed_at = NOW() WHERE id = $5;
`
	updateBookCoverQuery = `
UPDATE books SET cover_status = $2, cover_url = $3, thumbnail_url = $4 WHERE id = $1;
`
	deleteCoverByBookIdQuery = `
DELETE FROM covers WHERE book_id = $1;
`
)

type CoverRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewCoverRepository(db *sqlx.DB, logger *zap.Logger) *CoverRepository {
	return &CoverRepository{db: db, logger: logger}
}

func (r *CoverRepository) UpdateStatus(
	ctx context.Context,
	id string,
	status string,
	coverPath string,
	thumbPath string,
	errorMsg string,
) error {
	query := updateCoverStateQuery
	if status == "ready" || status == "failed" {
		query = updateCoverStateCompletedQuery
	}

	_, err := r.db.ExecContext(ctx, query, status, coverPath, thumbPath, errorMsg, id)
	if err != nil {
		r.logger.Error("failed to update cover status", zap.String("cover_id", id), zap.Error(err))

		return err
	}

	return nil
}

func (r *CoverRepository) UpdateBookCover(
	ctx context.Context,
	bookID string,
	status string,
	coverURL string,
	thumbURL string,
) error {
	_, err := r.db.ExecContext(ctx, updateBookCoverQuery, bookID, status, coverURL, thumbURL)
	if err != nil {
		r.logger.Error("failed to update book cover", zap.String("book_id", bookID), zap.Error(err))

		return err
	}

	return nil
}
