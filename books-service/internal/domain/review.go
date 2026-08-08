package domain

import (
	"database/sql"
	"errors"
	"time"
)

type Review struct {
	ID        string         `json:"id" db:"id"`
	BookID    string         `json:"book_id" db:"book_id"`
	UserID    string         `json:"user_id" db:"user_id"`
	Rating    int            `json:"rating" db:"rating"`
	Title     sql.NullString `json:"title" db:"title"`
	Content   string         `json:"content" db:"content"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`
}

func (r *Review) ToResponse() ReviewResponse {
	review := ReviewResponse{
		ID:        r.ID,
		BookID:    r.BookID,
		UserID:    r.UserID,
		Rating:    r.Rating,
		Content:   r.Content,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}

	_ = r.Title.Scan(&review.Title)

	return review
}

type ReviewResponse struct {
	ID        string    `json:"id"`
	BookID    string    `json:"book_id"`
	UserID    string    `json:"user_id"`
	Rating    int       `json:"rating"`
	Title     *string   `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateReviewRequest struct {
	Rating  int     `json:"rating"`
	Content string  `json:"content"`
	Title   *string `json:"title,omitempty"`
}

func (r *CreateReviewRequest) Validate() error {
	if r.Rating < 1 || r.Rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}

	if len(r.Content) < 10 {
		return errors.New("content must be at least 10 characters")
	}

	return nil
}

type UpdateReviewRequest struct {
	Rating  *int    `json:"rating,omitempty"`
	Content *string `json:"content,omitempty"`
	Title   *string `json:"title,omitempty"`
}

type ReviewListResponse struct {
	Data       []ReviewResponse `json:"data"`
	Pagination Pagination       `json:"pagination"`
}
