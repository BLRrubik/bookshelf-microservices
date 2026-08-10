package service

import (
	"bookshelf/books-service/internal/client"
	"bookshelf/books-service/internal/repository"
)

type Service struct {
	BookService   *BookService
	ReviewService *ReviewService
}

func New(repos *repository.Repository, authClient *client.AuthClient) *Service {
	return &Service{
		BookService:   NewBookService(repos.BookRepository),
		ReviewService: NewReviewService(repos.ReviewRepository, repos.BookRepository, authClient),
	}
}
