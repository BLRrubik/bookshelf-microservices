package repository

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")

	ErrBookNotFound = errors.New("book not found")

	ErrReviewNotFound = errors.New("review not found")
	ErrCoverNotFound  = errors.New("cover not found")
)

type Repository struct {
	BookRepository   *BookRepository
	ReviewRepository *ReviewRepository
	CoverRepository  *CoverRepository
}

func New(db *sqlx.DB, logger *zap.Logger) *Repository {
	return &Repository{
		BookRepository:   NewBookRepository(db, logger.Named("book_repository")),
		ReviewRepository: NewReviewRepository(db, logger.Named("review_repository")),
		CoverRepository:  NewCoverRepository(db, logger.Named("cover_repository")),
	}
}
