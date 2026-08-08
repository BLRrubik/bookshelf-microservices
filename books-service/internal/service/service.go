package service

import (
	"bookshelf/books-service/internal/repository"
)

type Service struct {
	BookService   *BookService
	ReviewService *ReviewService
}

func New(repos *repository.Repository, jwtSecret string) *Service {
	return &Service{
		BookService:   NewBookService(repos.BookRepository),
		ReviewService: NewReviewService(repos.ReviewRepository, repos.BookRepository),
	}
}
