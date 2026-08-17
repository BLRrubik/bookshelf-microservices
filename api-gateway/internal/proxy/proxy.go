package proxy

import (
	"bookshelf/api-gateway/internal/domain"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"bookshelf/api-gateway/internal/middleware"
)

type ServiceProxy struct {
	client          *http.Client
	authServiceURL  string
	booksServiceURL string
	serviceKey      string
}

func New(authURL, booksURL, serviceKey string) *ServiceProxy {
	return &ServiceProxy{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
		authServiceURL:  authURL,
		booksServiceURL: booksURL,
		serviceKey:      serviceKey,
	}
}

type verifyTokenResponse struct {
	Valid  bool   `json:"valid"`
	UserID string `json:"user_id"`
}

// VerifyToken asks auth-service's internal API whether token is valid,
// returning the associated user ID.
func (p *ServiceProxy) VerifyToken(ctx context.Context, token string) (string, error) {
	reqBody, err := json.Marshal(map[string]string{"token": token})
	if err != nil {
		return "", fmt.Errorf("marshal verify request: %w", err)
	}

	url := strings.TrimRight(p.authServiceURL, "/") + "/internal/v1/auth/verify"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return "", fmt.Errorf("create verify request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Service-Key", p.serviceKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("verify token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("verify token: unexpected status %d", resp.StatusCode)
	}

	var result verifyTokenResponse
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode verify response: %w", err)
	}

	if !result.Valid {
		return "", errors.New("invalid token")
	}

	return result.UserID, nil
}

func (p *ServiceProxy) ProxyAuthPath(w http.ResponseWriter, r *http.Request, path string) int {
	return p.ProxyRequest(w, r, joinPath(p.authServiceURL, path, r.URL.RawQuery))
}

func (p *ServiceProxy) ProxyBooksPath(w http.ResponseWriter, r *http.Request, path string) int {
	return p.ProxyRequest(w, r, joinPath(p.booksServiceURL, path, r.URL.RawQuery))
}

func (p *ServiceProxy) ProxyRequest(w http.ResponseWriter, r *http.Request, targetURL string) int {
	req, err := p.newUpstreamRequest(r, targetURL)
	if err != nil {
		p.handleProxyError(w, r, targetURL, err)

		return 0
	}

	response, err := p.client.Do(req)
	if err != nil {
		p.handleProxyError(w, r, targetURL, err)

		return 0
	}
	defer response.Body.Close()

	removeHopByHop(response.Header)
	copyHeaders(w.Header(), response.Header)
	w.WriteHeader(response.StatusCode)

	if _, err = io.Copy(w, response.Body); err != nil {
		slog.WarnContext(
			r.Context(),
			"copy upstream response",
			"error", err,
		)
	}

	return response.StatusCode
}

func (p *ServiceProxy) CheckAuthService(ctx context.Context) error {
	return p.HealthCheck(ctx, p.authServiceURL)
}

func (p *ServiceProxy) CheckBooksService(ctx context.Context) error {
	return p.HealthCheck(ctx, p.booksServiceURL)
}

// HealthCheck queries {baseURL}/health and reports an error unless the
// upstream service responds with 200 OK before ctx's deadline.
func (p *ServiceProxy) HealthCheck(ctx context.Context, baseURL string) error {
	url := strings.TrimRight(baseURL, "/") + "/health"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("create health check request: %w", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("health check request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check: unexpected status %d", resp.StatusCode)
	}

	return nil
}

func (p *ServiceProxy) GetWithCacheBooksPath(
	ctx context.Context,
	url string,
	headers map[string]string,
) (statusCode int, body []byte, err error) {
	return p.GetWithCache(ctx, p.booksServiceURL+url, headers)
}

func (p *ServiceProxy) GetWithCacheAuthPath(
	ctx context.Context,
	url string,
	headers map[string]string,
) (statusCode int, body []byte, err error) {
	return p.GetWithCache(ctx, p.authServiceURL+url, headers)
}

func (p *ServiceProxy) GetWithCache(ctx context.Context, url string, headers map[string]string) (statusCode int, body []byte, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, nil, fmt.Errorf("create request: %w", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("get request: %w", err)
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, fmt.Errorf("read response body: %w", err)
	}

	return resp.StatusCode, buf, nil
}

func (p *ServiceProxy) newUpstreamRequest(r *http.Request, targetURL string) (*http.Request, error) {
	out, err := http.NewRequestWithContext(
		r.Context(),
		r.Method,
		targetURL,
		r.Body,
	)
	if err != nil {
		return nil, fmt.Errorf("create upstream request: %w", err)
	}
	out.ContentLength = r.ContentLength

	copyHeaders(out.Header, r.Header)
	removeHopByHop(out.Header)
	setForwardedHeaders(out.Header, r)

	return out, nil
}

func joinPath(base, incoming, query string) string {
	return strings.TrimRight(base, "/") + "/" +
		strings.TrimLeft(incoming, "/") + "?" + query
}

func copyHeaders(dst, src http.Header) {
	for key, values := range src {
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}

func setForwardedHeaders(header http.Header, r *http.Request) {
	ip := getClientIP(r)
	header.Set("X-Forwarded-For", ip)
	header.Set("X-Real-Ip", ip)

	proto := "http"
	if r.TLS != nil {
		proto = "https"
	}
	header.Set("X-Forwarded-Proto", proto)
	header.Set("X-Forwarded-Host", r.Host)
}

func (p *ServiceProxy) handleProxyError(w http.ResponseWriter, r *http.Request, targetURL string, err error) {
	requestID := middleware.GetRequestID(r.Context())

	status, code, message := classifyProxyError(err)

	slog.ErrorContext(
		r.Context(),
		"proxy request failed",
		"target", targetURL,
		"error", err,
		"code", code,
		"request_id", requestID,
	)

	writeError(w, status, code, message, requestID)
}

func classifyProxyError(err error) (status int, code, message string) {
	var netErr net.Error
	if errors.Is(err, context.DeadlineExceeded) ||
		(errors.As(err, &netErr) && netErr.Timeout()) {
		return http.StatusGatewayTimeout, "GATEWAY_TIMEOUT", "Backend service timeout"
	}

	return http.StatusBadGateway, "BAD_GATEWAY", "Backend service unavailable"
}

func writeError(w http.ResponseWriter, status int, code, message, requestID string) {
	body := domain.ErrorResponse{
		Code:    code,
		Message: message,
	}
	if requestID != "" {
		body.RequestID = requestID
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		slog.Error("write proxy error response", "error", err)
	}
}

var hopByHopHeaders = []string{
	"Connection",
	"Proxy-Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",
	"Trailer",
	"Transfer-Encoding",
	"Upgrade",
}

func removeHopByHop(header http.Header) {
	for _, name := range strings.Split(header.Get("Connection"), ",") {
		if name = strings.TrimSpace(name); name != "" {
			header.Del(name)
		}
	}
	for _, name := range hopByHopHeaders {
		header.Del(name)
	}
}

func getClientIP(r *http.Request) string {
	splitted := strings.Split(r.Header.Get("X-Forwarded-For"), ",")
	if len(splitted[0]) > 0 {
		return splitted[0]
	}

	ip := r.Header.Get("X-Real-Ip")
	if len(ip) > 0 {
		return ip
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return host
}
