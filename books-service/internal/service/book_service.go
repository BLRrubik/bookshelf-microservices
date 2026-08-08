package service

import (
	"bookshelf/books-service/internal/domain"
	"bookshelf/books-service/internal/repository"
	"bookshelf/books-service/internal/utils"
	"context"
	"errors"
)

var (
	ErrBookNotFound    = errors.New("book not found")
	ErrNotBookOwner    = errors.New("user is not the owner of the book")
	ErrBookTitleEmpty  = errors.New("book title cannot be empty")
	ErrBookAuthorEmpty = errors.New("book author cannot be empty")
)

type BookService struct {
	bookRepo *repository.BookRepository
}

func NewBookService(bookRepo *repository.BookRepository) *BookService {
	return &BookService{
		bookRepo: bookRepo,
	}
}

func (s *BookService) Create(ctx context.Context, userID string, req domain.CreateBookRequest) (*domain.Book, error) {
	//creator, err := s.userRepo.GetByID(ctx, userID)
	//if err != nil {
	//	if errors.Is(err, ErrUserNotFound) {
	//		return nil, ErrUserNotFound
	//	}
	//
	//	return nil, err
	//}

	book := domain.Book{
		Title:         req.Title,
		Author:        req.Author,
		UserID:        userID,
		Description:   utils.StringToNull(req.Description),
		ISBN:          utils.StringToNull(req.ISBN),
		PublishedYear: utils.Int32ToNull(req.PublishedYear),
	}

	if err := s.bookRepo.Create(ctx, &book); err != nil {
		return nil, err
	}

	return &book, nil
}

func (s *BookService) GetByID(ctx context.Context, id string) (*domain.Book, error) {
	return s.bookRepo.GetByID(ctx, id)
}

func (s *BookService) List(ctx context.Context, params domain.ListParams) ([]domain.Book, int, error) {
	params.Normalize()

	return s.bookRepo.List(ctx, params)
}

func (s *BookService) ListByUser(ctx context.Context, userID string, params domain.ListParams) ([]domain.Book, int, error) {
	params.Normalize()

	books, count, err := s.bookRepo.ListByUserID(ctx, userID, params)
	if err != nil {
		return nil, 0, err
	}

	return books, count, nil
}

func (s *BookService) Update(
	ctx context.Context,
	userID string,
	bookID string,
	req domain.UpdateBookRequest,
) (*domain.Book, error) {
	book, err := s.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		return nil, ErrBookNotFound
	}

	if book.UserID != userID {
		return nil, ErrNotBookOwner
	}

	if req.Title != nil {
		book.Title = *req.Title
	}

	if req.Author != nil {
		book.Author = *req.Author
	}

	if req.Description != nil {
		book.Description = utils.StringToNull(req.Description)
	}

	if req.ISBN != nil {
		book.ISBN = utils.StringToNull(req.ISBN)
	}

	if req.PublishedYear != nil {
		book.PublishedYear = utils.Int32ToNull(req.PublishedYear)
	}

	if err = s.bookRepo.Update(ctx, book); err != nil {
		return nil, err
	}

	//creator, err := s.userRepo.GetByID(ctx, book.UserID)
	//if err != nil {
	//	return nil, ErrBookNotFound
	//}

	return book, nil
}

func (s *BookService) Delete(ctx context.Context, userID string, bookID string) error {
	book, err := s.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		if errors.Is(err, ErrBookNotFound) {
			return ErrBookNotFound
		}
	}

	if book.UserID != userID {
		return ErrNotBookOwner
	}

	return s.bookRepo.Delete(ctx, bookID)
}
