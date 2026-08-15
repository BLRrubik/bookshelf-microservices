package proxy

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"
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

func (p *ServiceProxy) ProxyToAuth(w http.ResponseWriter, r *http.Request) {
	p.ProxyRequest(w, r, joinPath(p.authServiceURL, r.URL.Path, r.URL.RawQuery))
}

func (p *ServiceProxy) ProxyToBooks(w http.ResponseWriter, r *http.Request) {
	p.ProxyRequest(w, r, joinPath(p.booksServiceURL, r.URL.Path, r.URL.RawQuery))
}

func (p *ServiceProxy) ProxyRequest(w http.ResponseWriter, r *http.Request, targetURL string) {
	req, err := p.newUpstreamRequest(r, targetURL)
	if err != nil {
		writeProxyError(w, err)

		return
	}

	response, err := p.client.Do(req)
	if err != nil {
		writeProxyError(w, err)

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
	header.Del("X-Forwarded-For")
	header.Del("X-Real-IP")
	header.Del("X-Forwarded-Proto")
	header.Del("X-Forwarded-Host")

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		header.Set("X-Forwarded-For", host)
		header.Set("X-Real-IP", host)
	}

	proto := "http"
	if r.TLS != nil {
		proto = "https"
	}
	header.Set("X-Forwarded-Proto", proto)
	header.Set("X-Forwarded-Host", r.Host)
}

func writeProxyError(w http.ResponseWriter, err error) {
	var netErr net.Error
	if errors.Is(err, context.DeadlineExceeded) ||
		(errors.As(err, &netErr) && netErr.Timeout()) {
		http.Error(w, "upstream timeout", http.StatusGatewayTimeout)

		return
	}

	http.Error(w, "upstream unavailable", http.StatusBadGateway)
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
