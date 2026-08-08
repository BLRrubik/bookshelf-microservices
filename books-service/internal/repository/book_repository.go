package repository

import (
	"bookshelf/books-service/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	createBookQuery = `
INSERT INTO books 
(id, title, author, description, isbn, published_year, created_by)
VALUES($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT DO NOTHING;
`
	getBookByIDQuery = `
SELECT 
	b.id, b.title, b.author, 
	b.description, b.isbn, b.published_year,
	b.created_by, b.created_at, b.updated_at,
	AVG(r.rating) as average_rating,
	COUNT(r.id) as reviews_count
FROM books AS b
LEFT JOIN reviews AS r ON b.id = r.book_id
WHERE b.id = $1
GROUP BY b.id
`
	listBooksQuery = `
SELECT 
	b.id, b.title, b.author, 
	b.description, b.isbn, b.published_year,
	b.created_by, b.created_at, b.updated_at,
	AVG(r.rating) as average_rating,
	COUNT(r.id) as reviews_count
FROM books AS b
LEFT JOIN reviews AS r ON b.id = r.book_id
WHERE b.title ILIKE $1
GROUP BY b.id
ORDER BY %s %s
LIMIT $2 OFFSET $3
`
	updateBookQuery = `
UPDATE books
SET title=$1, author=$2, description=$3, isbn=$4, published_year=$5
WHERE id = $6
`
	deleteBookQuery = `
DELETE FROM books
WHERE id = $1
`
)

type BookRepository struct {
	db *sqlx.DB
}

func NewBookRepository(db *sqlx.DB) *BookRepository {
	return &BookRepository{
		db: db,
	}
}

func (br *BookRepository) Create(ctx context.Context, book *domain.Book) error {
	book.ID = uuid.NewString()

	_, err := br.db.ExecContext(
		ctx,
		createBookQuery,
		book.ID,
		book.Title,
		book.Author,
		book.Description,
		book.ISBN,
		book.PublishedYear,
		book.UserID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (br *BookRepository) GetByID(ctx context.Context, id string) (*domain.Book, error) {
	var book domain.Book
	err := br.db.GetContext(ctx, &book, getBookByIDQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBookNotFound
		}

		return nil, err
	}

	return &book, nil
}

func (br *BookRepository) List(ctx context.Context, filter domain.BookFilter) ([]domain.Book, int, error) {
	var books []domain.Book
	var count int

	offset := (filter.Page - 1) * filter.Limit

	err := br.db.SelectContext(
		ctx,
		&books,
		fmt.Sprintf(listBooksQuery, filter.Sort, filter.Order),
		"%"+filter.Search+"%",
		filter.Limit,
		offset,
	)
	if err != nil {
		return nil, 0, err
	}

	err = br.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM books")
	if err != nil {
		return nil, 0, err
	}

	return books, count, nil
}

func (br *BookRepository) Update(ctx context.Context, book *domain.Book) error {
	_, err := br.db.ExecContext(
		ctx,
		updateBookQuery,
		book.Title,
		book.Author,
		book.Description,
		book.ISBN,
		book.PublishedYear,
		book.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (br *BookRepository) Delete(ctx context.Context, id string) error {
	_, err := br.db.ExecContext(ctx, deleteBookQuery, id)
	if err != nil {
		return err
	}

	return nil
}
