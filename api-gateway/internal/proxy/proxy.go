package proxy

import (
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
}

func New(authURL, booksURL string) *ServiceProxy {
	return &ServiceProxy{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
		authServiceURL:  authURL,
		booksServiceURL: booksURL,
	}
}

func (p *ServiceProxy) ProxyAuthPath(w http.ResponseWriter, r *http.Request, path string) {
	p.ProxyRequest(w, r, joinPath(p.authServiceURL, path, r.URL.RawQuery))
}

func (p *ServiceProxy) ProxyBooksPath(w http.ResponseWriter, r *http.Request, path string) {
	p.ProxyRequest(w, r, joinPath(p.booksServiceURL, path, r.URL.RawQuery))
}

func (p *ServiceProxy) ProxyRequest(w http.ResponseWriter, r *http.Request, targetURL string) {
	req, err := p.newUpstreamRequest(r, targetURL)
	if err != nil {
		p.handleProxyError(w, r, targetURL, err)

		return
	}

	response, err := p.client.Do(req)
	if err != nil {
		p.handleProxyError(w, r, targetURL, err)

		return
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
	body := map[string]string{
		"code":    code,
		"message": message,
	}
	if requestID != "" {
		body["request_id"] = requestID
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
