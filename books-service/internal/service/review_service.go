package service

import (
	"bookshelf/books-service/internal/domain"
	"bookshelf/books-service/internal/repository"
	"bookshelf/books-service/internal/utils"
	"context"
	"errors"
)

var (
	ErrReviewNotFound        = errors.New("review not found")
	ErrNotReviewOwner        = errors.New("not review owner")
	ErrAlreadyReviewed       = errors.New("user has already reviewed this book")
	ErrInvalidRating         = errors.New("invalid rating")
	ErrReviewContentTooShort = errors.New("review content is too short")
)

type ReviewService struct {
	reviewRepo *repository.ReviewRepository
	bookRepo   *repository.BookRepository
}

func NewReviewService(
	reviewRepo *repository.ReviewRepository,
	bookRepo *repository.BookRepository,
) *ReviewService {
	return &ReviewService{
		reviewRepo: reviewRepo,
		bookRepo:   bookRepo,
	}
}

func (s *ReviewService) Create(
	ctx context.Context,
	userID string,
	bookID string,
	req domain.CreateReviewRequest,
) (*domain.ReviewResponse, error) {
	book, err := s.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		if errors.Is(err, ErrBookNotFound) {
			return nil, ErrBookNotFound
		}

		return nil, err
	}

	reviewExist, err := s.reviewRepo.UserHasReviewedBook(ctx, userID, bookID)
	if err != nil {
		return nil, err
	}

	if reviewExist {
		return nil, ErrAlreadyReviewed
	}

	review := &domain.Review{
		BookID:  book.ID,
		Rating:  req.Rating,
		Title:   utils.StringToNull(req.Title),
		Content: req.Content,
	}

	if err = s.reviewRepo.Create(ctx, review); err != nil {
		return nil, err
	}

	reviewResponse := review.ToResponse()

	return &reviewResponse, nil
}

func (s *ReviewService) GetByID(ctx context.Context, id string) (*domain.Review, error) {
	review, err := s.reviewRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrReviewNotFound) {
			return nil, ErrReviewNotFound
		}

		return nil, err
	}

	return review, nil
}

func (s *ReviewService) ListByBook(ctx context.Context, bookID string) ([]domain.Review, error) {
	book, err := s.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		if errors.Is(err, repository.ErrBookNotFound) {
			return nil, ErrBookNotFound
		}

		return nil, err
	}

	reviews, err := s.reviewRepo.ListByBookID(ctx, book.ID)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

func (s *ReviewService) Update(
	ctx context.Context,
	userID string,
	reviewID string,
	req domain.UpdateReviewRequest,
) (*domain.Review, error) {
	review, err := s.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		if errors.Is(err, repository.ErrReviewNotFound) {
			return nil, ErrReviewNotFound
		}

		return nil, err
	}

	if review.UserID != userID {
		return nil, ErrNotReviewOwner
	}

	if err = s.validateUpdate(req); err != nil {
		return nil, err
	}

	if req.Rating != nil {
		review.Rating = *req.Rating
	}

	if req.Content != nil {
		review.Content = *req.Content
	}

	if req.Title != nil {
		review.Title = utils.StringToNull(req.Title)
	}

	if err = s.reviewRepo.Update(ctx, review); err != nil {
		return nil, err
	}

	return review, nil
}

func (s *ReviewService) validateUpdate(req domain.UpdateReviewRequest) error {
	if req.Rating != nil {
		if *req.Rating < 1 || *req.Rating > 5 {
			return ErrInvalidRating
		}
	}

	if req.Content != nil {
		if len(*req.Content) < 10 {
			return ErrReviewContentTooShort
		}
	}

	return nil
}

func (s *ReviewService) Delete(ctx context.Context, userID string, reviewID string) error {
	review, err := s.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		if errors.Is(err, repository.ErrReviewNotFound) {
			return ErrReviewNotFound
		}

		return err
	}

	if review.UserID != userID {
		return ErrNotReviewOwner
	}

	return s.reviewRepo.Delete(ctx, reviewID)
}
