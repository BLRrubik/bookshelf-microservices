package domain

import (
	"bookshelf/books-service/internal/utils"
	"database/sql"
	"errors"
	"time"
)

type Book struct {
	ID            string          `json:"id" db:"id"`
	Title         string          `json:"title" db:"title"`
	Author        string          `json:"author" db:"author"`
	UserID        string          `json:"created_by" db:"created_by"`
	CreatedAt     time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at" db:"updated_at"`
	Description   sql.NullString  `json:"description" db:"description"`
	ISBN          sql.NullString  `json:"isbn" db:"isbn"`
	PublishedYear sql.NullInt32   `json:"published_year" db:"published_year"`
	AverageRating sql.NullFloat64 `json:"-" db:"average_rating"`
	ReviewsCount  int             `json:"reviews_count" db:"reviews_count"`
}

func (b *Book) ToResponse() BookResponse {
	book := BookResponse{
		ID:            b.ID,
		Title:         b.Title,
		Author:        b.Author,
		CreatedBy:     b.UserID,
		CreatedAt:     b.CreatedAt,
		UpdatedAt:     b.UpdatedAt,
		Description:   utils.NullToString(b.Description),
		ISBN:          utils.NullToString(b.ISBN),
		PublishedYear: utils.NullToInt32(b.PublishedYear),
		AverageRating: utils.NullToFloat64(b.AverageRating),
		ReviewsCount:  b.ReviewsCount,
	}

	return book
}

type BookResponse struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	CreatedBy     string    `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Description   *string   `json:"description"`
	ISBN          *string   `json:"isbn"`
	PublishedYear *int32    `json:"published_year"`
	AverageRating *float64  `json:"-"`
	ReviewsCount  int       `json:"reviews_count"`
}

type CreateBookRequest struct {
	Title         string  `json:"title"`
	Author        string  `json:"author"`
	Description   *string `json:"description,omitempty"`
	ISBN          *string `json:"isbn,omitempty"`
	PublishedYear *int32  `json:"published_year,omitempty"`
}

func (r CreateBookRequest) Validate() error {
	if r.Title == "" {
		return errors.New("title is required")
	}

	if r.Author == "" {
		return errors.New("author is required")
	}

	return nil
}

type UpdateBookRequest struct {
	Title         *string `json:"title,omitempty"`
	Author        *string `json:"author,omitempty"`
	Description   *string `json:"description,omitempty"`
	ISBN          *string `json:"isbn,omitempty"`
	PublishedYear *int32  `json:"published_year,omitempty"`
}

type ListParams struct {
	Search string `json:"search"`
	Sort   string `json:"sort"`
	Order  string `json:"order"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
}

func (b *ListParams) Normalize() {
	if b.Page <= 0 {
		b.Page = 1
	}

	if b.Limit <= 0 {
		b.Limit = 10
	}

	if b.Sort != "" {
		b.Sort = "title"
	}

	if len(b.Search) == 0 {
		b.Search = ""
	}

	if b.Order != "asc" && b.Order != "desc" {
		b.Order = "desc"
	}
}

type BookListResponse struct {
	Data       []BookResponse `json:"data"`
	Pagination Pagination     `json:"pagination"`
}
