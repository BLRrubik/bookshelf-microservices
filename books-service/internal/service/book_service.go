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

func (s *BookService) Create(ctx context.Context, userID string, req domain.CreateBookRequest) (*domain.BookResponse, error) {
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
		Description:   utils.StringToNull(req.Description),
		ISBN:          utils.StringToNull(req.ISBN),
		PublishedYear: utils.Int32ToNull(req.PublishedYear),
	}

	if err := s.bookRepo.Create(ctx, &book); err != nil {
		return nil, err
	}

	bookResponse := book.ToResponse()

	return &bookResponse, nil
}

func (s *BookService) GetByID(ctx context.Context, id string) (*domain.BookResponse, error) {
	book, err := s.bookRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrBookNotFound
	}

	bookResponse := book.ToResponse()

	return &bookResponse, nil
}

func (s *BookService) List(ctx context.Context, filter domain.ListParams) (*domain.BookListResponse, error) {
	filter.Normalize()

	books, count, err := s.bookRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	respBooks := make([]domain.BookResponse, len(books))
	for i, book := range books {
		respBooks[i] = book.ToResponse()
	}

	return &domain.BookListResponse{
		Data: respBooks,
		Pagination: domain.Pagination{
			Limit:      filter.Limit,
			Page:       filter.Page,
			Total:      count,
			TotalPages: count / filter.Limit,
		},
	}, nil
}

func (s *BookService) Update(
	ctx context.Context,
	userID string,
	bookID string,
	req domain.UpdateBookRequest,
) (*domain.BookResponse, error) {
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

	//creator, err := s.userRepo.GetByID(ctx, book.CreatedBy)
	//if err != nil {
	//	return nil, ErrBookNotFound
	//}

	bookResponse := book.ToResponse()

	return &bookResponse, nil
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
