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

	"golang.org/x/sync/errgroup"
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

	authToken := r.Header.Get("Authorization")

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

	var dashboardResp domain.DashboardResponse

	group, groupCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		dashboardResp.PopularBooks = h.fetchPopularBooks(groupCtx, authToken)

		return nil
	})
	group.Go(func() error {
		dashboardResp.RecentReviews = h.fetchRecentReviews(groupCtx, authToken)

		return nil
	})
	if userID != "" {
		group.Go(func() error {
			dashboardResp.UserStats = h.fetchUserStats(groupCtx, authToken, userID)

			return nil
		})
	}

	group.Wait()

	if dashboardResp.PopularBooks == nil {
		dashboardResp.PopularBooks = []domain.BookSummary{}
	}

	if dashboardResp.RecentReviews == nil {
		dashboardResp.RecentReviews = []domain.RecentReview{}
	}

	bytes, err := json.Marshal(dashboardResp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Cache", "MISS")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
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

	var booksResp domain.ResponseWithPagination[domain.BookSummary]
	err = json.Unmarshal(body, &booksResp)
	if err != nil {
		return nil
	}

	return booksResp.Data
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

		var reviews domain.ResponseWithPagination[domain.RecentReview]
		if err = json.Unmarshal(rBody, &reviews); err != nil {
			continue
		}

		for _, rv := range reviews.Data {
			if len(result) >= 10 {
				break
			}

			rv.BookTitle = book.Title
			rv.Content = truncateContent(rv.Content, 200)

			result = append(result, rv)
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

	var resp domain.Pagination
	if err = json.Unmarshal(body, &resp); err != nil {
		return 0, err
	}

	return resp.Total, nil
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
