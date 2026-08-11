package repository

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")

	ErrBookNotFound = errors.New("book not found")

	ErrReviewNotFound = errors.New("review not found")
)

type Repository struct {
	BookRepository   *BookRepository
	ReviewRepository *ReviewRepository
	CoverRepository  *CoverRepository
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		BookRepository:   NewBookRepository(db),
		ReviewRepository: NewReviewRepository(db),
		CoverRepository:  NewCoverRepository(db),
	}
}
