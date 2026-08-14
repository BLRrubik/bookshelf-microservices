package service

import (
	"bookshelf/books-service/internal/client"
	"bookshelf/books-service/internal/domain"
	"bookshelf/books-service/internal/repository"
	"bookshelf/books-service/internal/utils"
	"context"
	"errors"

	"go.uber.org/zap"
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
	authClient *client.AuthClient
	logger     *zap.Logger
}

func NewReviewService(
	reviewRepo *repository.ReviewRepository,
	bookRepo *repository.BookRepository,
	authClient *client.AuthClient,
	logger *zap.Logger,
) *ReviewService {
	return &ReviewService{
		reviewRepo: reviewRepo,
		bookRepo:   bookRepo,
		authClient: authClient,
		logger:     logger,
	}
}

func (s *ReviewService) Create(
	ctx context.Context,
	userID string,
	bookID string,
	req domain.CreateReviewRequest,
) (*domain.ReviewResponse, error) {
	if req.Rating < 1 || req.Rating > 5 {
		return nil, errors.New("rating must be between 1 and 5")
	}

	if len(req.Content) < 10 {
		return nil, errors.New("content must be at least 10 characters")
	}

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
		UserID:  userID,
		Rating:  req.Rating,
		Title:   utils.StringToNull(req.Title),
		Content: req.Content,
	}

	if err = s.reviewRepo.Create(ctx, review); err != nil {
		return nil, err
	}

	s.logger.Info("review created", zap.String("review_id", review.ID), zap.String("book_id", bookID), zap.String("user_id", userID))

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

func (s *ReviewService) ListByBook(ctx context.Context, bookID string, page, limit int) (*domain.ReviewListResponse, error) {
	book, err := s.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		if errors.Is(err, repository.ErrBookNotFound) {
			return nil, ErrBookNotFound
		}

		return nil, err
	}

	reviews, count, err := s.reviewRepo.ListByBookID(ctx, book.ID, page, limit)
	if err != nil {
		return nil, err
	}

	users := make([]string, len(reviews))
	for i, review := range reviews {
		users[i] = review.UserID
	}

	usersMap := make(map[string]client.UserPublic, len(users))
	usersSummary, err := s.authClient.GetUsersByIDs(ctx, users)
	if err != nil {
		s.logger.Warn("failed to fetch review authors", zap.String("book_id", bookID), zap.Error(err))
	}

	for _, user := range usersSummary {
		usersMap[user.ID] = user
	}

	reviewsResp := make([]domain.ReviewResponse, len(reviews))
	for i, review := range reviews {
		reviewsResp[i] = review.ToResponse()
		userInfo := usersMap[review.UserID]
		reviewsResp[i].User = domain.UserSummary{
			ID:       userInfo.ID,
			Username: userInfo.Username,
		}
	}

	return &domain.ReviewListResponse{
		Data: reviewsResp,
		Pagination: domain.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      count,
			TotalPages: (count + limit - 1) / limit,
		},
	}, nil
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

	s.logger.Info("review updated", zap.String("review_id", review.ID), zap.String("user_id", userID))

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

	if err = s.reviewRepo.Delete(ctx, reviewID); err != nil {
		return err
	}

	s.logger.Info("review deleted", zap.String("review_id", reviewID), zap.String("user_id", userID))

	return nil
}
