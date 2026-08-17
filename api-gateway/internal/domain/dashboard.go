package domain

type BookSummary struct {
	ID       string
	Title    string
	Author   string
	Rating   float64
	CoverURL string
}

type RecentReview struct {
}

type UserStats struct {
	BooksAdded     int
	ReviewsWritten int
	TotalBooks     int
	TotalReviews   int
}

type DashboardResponse struct {
	PopularBooks  []BookSummary  `json:"popular_books"`
	RecentReviews []RecentReview `json:"recent_reviews"`
	UserStats     *UserStats     `json:"user_stats"`
}
