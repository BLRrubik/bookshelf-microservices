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
	b.cover_status, COALESCE(b.cover_url, '') AS cover_url, COALESCE(b.thumbnail_url, '') AS thumb_url,
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
	b.cover_status, COALESCE(b.cover_url, '') AS cover_url, COALESCE(b.thumbnail_url, '') AS thumb_url,
	AVG(r.rating) as average_rating,
	COUNT(r.id) as reviews_count
FROM books AS b
LEFT JOIN reviews AS r ON b.id = r.book_id
WHERE b.title ILIKE $1
GROUP BY b.id
ORDER BY %s %s
LIMIT $2 OFFSET $3
`
	listBooksByUserIDQuery = `
SELECT
	b.id, b.title, b.author,
	b.description, b.isbn, b.published_year,
	b.created_by, b.created_at, b.updated_at,
	b.cover_status, COALESCE(b.cover_url, '') AS cover_url, COALESCE(b.thumbnail_url, '') AS thumb_url,
	AVG(r.rating) as average_rating,
	COUNT(r.id) as reviews_count
FROM books AS b
LEFT JOIN reviews AS r ON b.id = r.book_id
WHERE b.created_by = $1 AND b.title ILIKE $2
GROUP BY b.id
ORDER BY %s %s
LIMIT $3 OFFSET $4
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

var allowedBookSortColumns = map[string]string{
	"title":          "b.title",
	"author":         "b.author",
	"published_year": "b.published_year",
	"created_at":     "b.created_at",
	"rating":         "average_rating",
}

func sortColumn(sort string) string {
	if col, ok := allowedBookSortColumns[sort]; ok {
		return col
	}

	return "b.title"
}

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
		return fmt.Errorf("create book: %w", err)
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

		return nil, fmt.Errorf("get book by id: %w", err)
	}

	return &book, nil
}

func (br *BookRepository) List(ctx context.Context, params domain.ListParams) ([]domain.Book, int, error) {
	var books []domain.Book
	var count int

	offset := (params.Page - 1) * params.Limit

	query := fmt.Sprintf(listBooksQuery, sortColumn(params.Sort), params.Order)
	err := br.db.SelectContext(
		ctx,
		&books,
		query,
		"%"+params.Search+"%",
		params.Limit,
		offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list books: %w", err)
	}

	err = br.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM books WHERE title ILIKE $1", "%"+params.Search+"%")
	if err != nil {
		return nil, 0, fmt.Errorf("count books: %w", err)
	}

	return books, count, nil
}

func (br *BookRepository) ListByUserID(ctx context.Context, userID string, params domain.ListParams) ([]domain.Book, int, error) {
	var books []domain.Book
	var count int

	offset := (params.Page - 1) * params.Limit

	err := br.db.SelectContext(
		ctx,
		&books,
		fmt.Sprintf(listBooksByUserIDQuery, sortColumn(params.Sort), params.Order),
		userID,
		"%"+params.Search+"%",
		params.Limit,
		offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list books by user: %w", err)
	}

	err = br.db.GetContext(
		ctx, &count, "SELECT COUNT(*) FROM books WHERE created_by = $1 AND title ILIKE $2",
		userID, "%"+params.Search+"%",
	)
	if err != nil {
		return nil, 0, fmt.Errorf("count books: %w", err)
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
		return fmt.Errorf("update book: %w", err)
	}

	return nil
}

func (br *BookRepository) Delete(ctx context.Context, id string) error {
	_, err := br.db.ExecContext(ctx, deleteBookQuery, id)
	if err != nil {
		return fmt.Errorf("delete book: %w", err)
	}

	return nil
}
