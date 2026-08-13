package repository

import (
	"bookshelf/books-service/internal/domain"
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	createCoverQuery = `
INSERT INTO covers 
(id, book_id, status, original_path, cover_path, thumb_path, error)
VALUES ($1, $2, $3, $4, $5, $6, $7)
`
	getCoverByBookIdQuery = `
SELECT id, book_id, status, original_path, cover_path, thumb_path, error
FROM covers
WHERE book_id = $1
ORDER BY created_at DESC LIMIT 1
`

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
	db *sqlx.DB
}

func NewCoverRepository(db *sqlx.DB) *CoverRepository {
	return &CoverRepository{db: db}
}

func (r *CoverRepository) Create(ctx context.Context, cover *domain.Cover) error {
	cover.ID = uuid.NewString()

	_, err := r.db.ExecContext(
		ctx,
		createCoverQuery,
		cover.ID,
		cover.BookID,
		cover.Status,
		cover.OriginalPath,
		cover.CoverPath,
		cover.CoverPath,
		cover.Error,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *CoverRepository) GetByBookID(ctx context.Context, bookID string) (*domain.Cover, error) {
	var cover domain.Cover

	if err := r.db.GetContext(ctx, &cover, getCoverByBookIdQuery, bookID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCoverNotFound
		}
		return nil, err
	}

	return &cover, nil
}

func (r *CoverRepository) UpdateStatus(
	ctx context.Context,
	id string,
	status domain.CoverStatus,
	coverPath string,
	thumbPath string,
	errorMsg string,
) error {
	query := updateCoverStateQuery
	if status == domain.CoverStatusReady || status == domain.CoverStatusFailed {
		query = updateCoverStateCompletedQuery
	}

	_, err := r.db.ExecContext(ctx, query, status, coverPath, thumbPath, errorMsg, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *CoverRepository) UpdateBookCover(
	ctx context.Context,
	bookID string,
	status domain.CoverStatus,
	coverURL string,
	thumbURL string,
) error {
	_, err := r.db.ExecContext(ctx, updateBookCoverQuery, bookID, status, coverURL, thumbURL)
	if err != nil {
		return err
	}

	return nil
}

func (r *CoverRepository) DeleteByBookID(ctx context.Context, bookID string) error {
	_, err := r.db.ExecContext(ctx, deleteCoverByBookIdQuery, bookID)
	if err != nil {
		return err
	}

	return nil
}
