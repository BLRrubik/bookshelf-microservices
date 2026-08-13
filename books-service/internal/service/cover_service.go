package service

import (
	"bookshelf/books-service/internal/client"
	"bookshelf/books-service/internal/domain"
	"bookshelf/books-service/internal/repository"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
)

var ErrForbidden = errors.New("forbidden")

type CoverService struct {
	coverRepo    *repository.CoverRepository
	bookRepo     *repository.BookRepository
	minioClient  *client.MinIOClient
	rabbitClient *client.RabbitMQClient
}

func NewCoverService(
	coverRepo *repository.CoverRepository,
	bookRepo *repository.BookRepository,
	minioClient *client.MinIOClient,
	rabbitClient *client.RabbitMQClient,
) *CoverService {
	return &CoverService{
		coverRepo:    coverRepo,
		bookRepo:     bookRepo,
		minioClient:  minioClient,
		rabbitClient: rabbitClient,
	}
}

func (s *CoverService) UploadCover(
	ctx context.Context,
	userID string,
	bookID string,
	file io.Reader,
	fileSize int64,
	filename string,
) (*domain.CoverUploadResponse, error) {
	book, err := s.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		return nil, ErrBookNotFound
	}

	if book.UserID != userID {
		return nil, fmt.Errorf("%w: invalid user id", ErrForbidden)
	}

	if err = validateImageFile(filename, fileSize); err != nil {
		return nil, err
	}

	splitted := strings.Split(filename, ".")
	ext := splitted[len(splitted)-1]

	switch ext {
	case "jpg", "jpeg", "png", "webp":
	default:
		return nil, errors.New("invalid file type")
	}

	filepath := fmt.Sprintf("covers/%s/original.%s", book.ID, ext)

	if err = s.minioClient.UploadFile(ctx, filepath, file, fileSize, client.GetContentType(filename)); err != nil {
		return nil, err
	}

	cover := &domain.Cover{
		BookID:       book.ID,
		Status:       domain.CoverStatusProcessing,
		OriginalPath: filepath,
	}

	if err = s.coverRepo.Create(ctx, cover); err != nil {
		return nil, err
	}

	if err = s.coverRepo.UpdateBookCover(
		ctx,
		bookID,
		domain.CoverStatusProcessing,
		"",
		"",
	); err != nil {
		return nil, err
	}

	if err = s.rabbitClient.PublishImageCompress(ctx, client.ImageCompressMessage{
		BookID:       bookID,
		CoverID:      cover.ID,
		OriginalPath: filepath,
	}); err != nil {
		return nil, err
	}

	return &domain.CoverUploadResponse{
		CoverID: cover.ID,
		Status:  domain.CoverStatusProcessing,
	}, err
}

func validateImageFile(filename string, size int64) error {
	if filename == "" {
		return errors.New("empty filename")
	}

	if size <= 0 || size > 5*(1<<20) {
		return errors.New("file size too big")
	}

	return nil
}

func (s *CoverService) GetCover(ctx context.Context, bookID string) (*domain.CoverResponse, error) {
	cover, err := s.coverRepo.GetByBookID(ctx, bookID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrCoverNotFound):
			return &domain.CoverResponse{
				Status:  domain.CoverStatusNone,
				Message: "No cover",
			}, nil
		}

		return nil, err
	}

	resp := &domain.CoverResponse{
		Status:   cover.Status,
		CoverURL: s.minioClient.GetFileURL(cover.CoverPath),
		ThumbURL: s.minioClient.GetFileURL(cover.ThumbPath),
		Message:  cover.Error,
	}

	if cover.Status == domain.CoverStatusProcessing {
		resp.Message = "Cover is being processed. Please wait."
	}

	return resp, nil
}

func (s *CoverService) GetCoverStatus(ctx context.Context, bookID string) (*domain.CoverStatusResponse, error) {
	cover, err := s.coverRepo.GetByBookID(ctx, bookID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrCoverNotFound):
			return &domain.CoverStatusResponse{
				Status: domain.CoverStatusNone,
			}, nil
		}

		return nil, err
	}

	return &domain.CoverStatusResponse{
		CoverID:     cover.ID,
		Status:      cover.Status,
		CoverURL:    s.minioClient.GetFileURL(cover.CoverPath),
		ThumbURL:    s.minioClient.GetFileURL(cover.ThumbPath),
		Error:       cover.Error,
		CreatedAt:   cover.CreatedAt,
		CompletedAt: cover.CompletedAt,
	}, nil
}

func (s *CoverService) DeleteCover(ctx context.Context, userID, bookID string) error {
	book, err := s.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		return err
	}

	if book.UserID != userID {
		return fmt.Errorf("%w: invalid user id", ErrForbidden)
	}

	cover, err := s.coverRepo.GetByBookID(ctx, bookID)
	if err != nil {
		return err
	}

	files := []string{cover.OriginalPath, cover.ThumbPath, cover.ThumbPath}
	for _, file := range files {
		if err = s.minioClient.DeleteFile(ctx, file); err != nil {
			return err
		}
	}

	if err = s.coverRepo.DeleteByBookID(ctx, bookID); err != nil {
		return err
	}

	return s.coverRepo.UpdateBookCover(ctx, bookID, domain.CoverStatusNone, "", "")
}
