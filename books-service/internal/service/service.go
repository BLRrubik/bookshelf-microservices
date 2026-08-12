package service

import (
	"bookshelf/books-service/internal/client"
	"bookshelf/books-service/internal/repository"
)

type Service struct {
	BookService   *BookService
	ReviewService *ReviewService
	CoverService  *CoverService
}

func New(
	repos *repository.Repository,
	authClient *client.AuthClient,
	minioClient *client.MinIOClient,
	rabbitMQClient *client.RabbitMQClient,
) *Service {
	return &Service{
		BookService:   NewBookService(repos.BookRepository),
		ReviewService: NewReviewService(repos.ReviewRepository, repos.BookRepository, authClient),
		CoverService:  NewCoverService(repos.CoverRepository, repos.BookRepository, minioClient, rabbitMQClient),
	}
}
