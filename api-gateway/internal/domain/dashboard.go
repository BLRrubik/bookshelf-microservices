package domain

import "time"

type BookSummary struct {
	ID       string  `json:"id"`
	Title    string  `json:"title"`
	Author   string  `json:"author"`
	Rating   float64 `json:"rating"`
	CoverURL string  `json:"cover_url"`
}

type RecentReview struct {
	ID        string    `json:"id"`
	BookID    string    `json:"book_id"`
	BookTitle string    `json:"book_title"`
	Rating    int       `json:"rating"`
	Content   string    `json:"content"`
	User      UserInfo  `json:"user"`
	CreatedAt time.Time `json:"created_at"`
}

type UserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type UserStats struct {
	BooksAdded     int `json:"books_added"`
	ReviewsWritten int `json:"reviews_written"`
	TotalBooks     int `json:"total_books"`
	TotalReviews   int `json:"total_reviews"`
}

type DashboardResponse struct {
	PopularBooks  []BookSummary  `json:"popular_books"`
	RecentReviews []RecentReview `json:"recent_reviews"`
	UserStats     *UserStats     `json:"user_stats"`
}
