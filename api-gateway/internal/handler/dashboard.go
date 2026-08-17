package handler

import (
	"bookshelf/api-gateway/internal/cache"
	"bookshelf/api-gateway/internal/domain"
	"bookshelf/api-gateway/internal/proxy"
	"bookshelf/api-gateway/internal/utils"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
	"unicode/utf8"
)

type DashboardHandler struct {
	proxy *proxy.ServiceProxy
	cache *cache.Cache
	ttl   time.Duration
}

func NewDashboardHandler(proxy *proxy.ServiceProxy, cache *cache.Cache, ttl time.Duration) *DashboardHandler {
	return &DashboardHandler{proxy: proxy, cache: cache, ttl: ttl}
}

func (h *DashboardHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var userID string
	if token := utils.ExtractBearerToken(r); token != "" {
		userID = h.verifyToken(ctx, token)
	}

	cacheKey := dashboardCacheKey(userID)

	if cached, err := h.cache.Get(ctx, cacheKey); err == nil {
		w.Header().Set("X-Cache", "HIT")
		w.Header().Set("Content-Type", "application/json")
		w.Write(cached)

		return
	}

	w.Header().Set("X-Cache", "MISS")
	w.WriteHeader(http.StatusOK)
}

func dashboardCacheKey(userID string) string {
	if userID == "" {
		return "dashboard:anon"
	}

	if len(userID) > 8 {
		userID = userID[:8]
	}

	return "dashboard:user:" + userID
}

func (h *DashboardHandler) verifyToken(ctx context.Context, token string) string {
	userID, err := h.proxy.VerifyToken(ctx, token)
	if err != nil {
		return ""
	}

	return userID
}

func (h *DashboardHandler) fetchPopularBooks(ctx context.Context, token string) []domain.BookSummary {
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": token,
	}

	code, body, err := h.proxy.GetWithCacheBooksPath(ctx, "/api/v1/books?sort=rating&order=desc&limit=10", headers)
	if err != nil {
		return nil
	}

	if code != http.StatusOK {
		return nil
	}

	var books []domain.BookSummary
	err = json.Unmarshal(body, &books)
	if err != nil {
		return nil
	}

	return books
}

type reviewsResponse struct {
	Data []struct {
		ID        string    `json:"id"`
		BookID    string    `json:"book_id"`
		Rating    int       `json:"rating"`
		Content   string    `json:"content"`
		CreatedAt time.Time `json:"created_at"`
		User      struct {
			ID       string `json:"id"`
			Username string `json:"username"`
		} `json:"user"`
	} `json:"data"`
}

func (h *DashboardHandler) fetchRecentReviews(ctx context.Context, token string) []domain.RecentReview {
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": token,
	}

	code, body, err := h.proxy.GetWithCacheBooksPath(ctx, "/api/v1/books?limit=10", headers)
	if err != nil {
		return nil
	}

	if code != http.StatusOK {
		return nil
	}

	var books []domain.BookSummary
	if err = json.Unmarshal(body, &books); err != nil {
		return nil
	}

	result := []domain.RecentReview{}

	for _, book := range books {
		if len(result) >= 10 {
			break
		}

		rCode, rBody, rErr := h.proxy.GetWithCacheBooksPath(
			ctx, "/api/v1/books/"+book.ID+"/reviews?limit=3", headers,
		)
		if rErr != nil || rCode != http.StatusOK {
			continue
		}

		var reviews reviewsResponse
		if err = json.Unmarshal(rBody, &reviews); err != nil {
			continue
		}

		for _, rv := range reviews.Data {
			if len(result) >= 10 {
				break
			}

			result = append(result, domain.RecentReview{
				ID:        rv.ID,
				BookID:    rv.BookID,
				BookTitle: book.Title,
				Rating:    rv.Rating,
				Content:   truncateContent(rv.Content, 200),
				User: domain.UserInfo{
					ID:       rv.User.ID,
					Username: rv.User.Username,
				},
				CreatedAt: rv.CreatedAt,
			})
		}
	}

	return result
}

func truncateContent(content string, maxRunes int) string {
	if utf8.RuneCountInString(content) <= maxRunes {
		return content
	}

	runes := []rune(content)

	return string(runes[:maxRunes]) + "..."
}

type paginatedResponse struct {
	Pagination struct {
		Total int `json:"total"`
	} `json:"pagination"`
}

func (h *DashboardHandler) fetchUserStats(ctx context.Context, token, userID string) *domain.UserStats {
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": token,
	}

	booksAdded, err := h.fetchTotal(ctx, "/api/v1/books?created_by="+userID+"&limit=1", headers)
	if err != nil {
		return nil
	}

	reviewsWritten, err := h.fetchReviewCount(ctx, "/api/v1/users/"+userID+"/reviews?limit=1", headers)
	if err != nil {
		return nil
	}

	totalBooks, err := h.fetchTotal(ctx, "/api/v1/books?limit=1", headers)
	if err != nil {
		return nil
	}

	return &domain.UserStats{
		TotalBooks:     totalBooks,
		ReviewsWritten: reviewsWritten,
		BooksAdded:     booksAdded,
	}
}

func (h *DashboardHandler) fetchTotal(ctx context.Context, url string, headers map[string]string) (int, error) {
	code, body, err := h.proxy.GetWithCacheBooksPath(ctx, url, headers)
	if err != nil {
		return 0, err
	}

	if code != http.StatusOK {
		return 0, errors.New("invalid status code: " + strconv.Itoa(code))
	}

	var resp paginatedResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return 0, err
	}

	return resp.Pagination.Total, nil
}

func (h *DashboardHandler) fetchReviewCount(ctx context.Context, url string, headers map[string]string) (int, error) {
	code, body, err := h.proxy.GetWithCacheBooksPath(ctx, url, headers)
	if err != nil {
		return 0, err
	}

	if code != http.StatusOK {
		return 0, errors.New("invalid status code: " + strconv.Itoa(code))
	}

	var reviews []json.RawMessage
	if err = json.Unmarshal(body, &reviews); err != nil {
		return 0, err
	}

	return len(reviews), nil
}
