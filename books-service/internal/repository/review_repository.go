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
	createReviewQuery = `
INSERT INTO reviews 
(id, user_id, book_id, rating, title, content)
VALUES ($1, $2, $3, $4, $5, $6)
`
	getReviewByIDQuery = `
SELECT id, user_id, book_id, rating, title, content, created_at, updated_at 
FROM reviews
WHERE id = $1
`
	listReviewsByBookIDQuery = `
SELECT id, user_id, book_id, rating, title, content, created_at, updated_at
FROM reviews
WHERE book_id = $1
LIMIT $2 OFFSET $3
`
	updateReviewQuery = `
UPDATE reviews
SET rating = $1, title = $2, content = $3
WHERE id = $4
`
	deleteReviewQuery = `
DELETE FROM reviews
WHERE id = $1
`
	hasUserBookReviewQuery = `
SELECT EXISTS (
SELECT id FROM reviews WHERE user_id = $1 and book_id = $2
)
`
)

type ReviewRepository struct {
	db *sqlx.DB
}

func NewReviewRepository(db *sqlx.DB) *ReviewRepository {
	return &ReviewRepository{
		db: db,
	}
}

func (rr *ReviewRepository) Create(ctx context.Context, review *domain.Review) error {
	review.ID = uuid.NewString()

	_, err := rr.db.ExecContext(
		ctx,
		createReviewQuery,
		review.ID,
		review.UserID,
		review.BookID,
		review.Rating,
		review.Title,
		review.Content,
	)
	if err != nil {
		return err
	}

	return nil
}

func (rr *ReviewRepository) GetByID(ctx context.Context, id string) (*domain.Review, error) {
	var review domain.Review

	err := rr.db.GetContext(ctx, &review, getReviewByIDQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrReviewNotFound
		}

		return nil, err
	}

	return &review, nil
}

func (rr *ReviewRepository) ListByBookID(ctx context.Context, bookID string, page, limit int) ([]domain.Review, int, error) {
	var reviews []domain.Review

	offset := (page - 1) * limit

	err := rr.db.SelectContext(
		ctx,
		&reviews,
		listReviewsByBookIDQuery,
		bookID,
		limit,
		offset,
	)
	if err != nil {
		return nil, 0, err
	}

	var count int
	err = rr.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM reviews WHERE book_id = $1", bookID)
	if err != nil {
		return nil, 0, err
	}

	return reviews, count, nil
}

func (rr *ReviewRepository) Update(ctx context.Context, review *domain.Review) error {
	_, err := rr.db.ExecContext(
		ctx,
		updateReviewQuery,
		review.Rating,
		review.Title,
		review.Content,
		review.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (rr *ReviewRepository) Delete(ctx context.Context, id string) error {
	_, err := rr.db.ExecContext(ctx, deleteReviewQuery, id)
	if err != nil {
		return err
	}

	return nil
}

func (rr *ReviewRepository) UserHasReviewedBook(ctx context.Context, userID string, bookID string) (bool, error) {
	var exists bool
	err := rr.db.GetContext(ctx, &exists, hasUserBookReviewQuery, userID, bookID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return exists, nil
}
