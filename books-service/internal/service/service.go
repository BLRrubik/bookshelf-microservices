package service

import (
	"bookshelf/books-service/internal/client"
	"bookshelf/books-service/internal/repository"

	"go.uber.org/zap"
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
	logger *zap.Logger,
) *Service {
	return &Service{
		BookService:   NewBookService(repos.BookRepository, logger.Named("book_service")),
		ReviewService: NewReviewService(repos.ReviewRepository, repos.BookRepository, authClient, logger.Named("review_service")),
		CoverService:  NewCoverService(repos.CoverRepository, repos.BookRepository, minioClient, rabbitMQClient, logger.Named("cover_service")),
	}
}
