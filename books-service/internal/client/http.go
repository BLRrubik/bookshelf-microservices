package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"time"

	"go.uber.org/zap"
)

var (
	ErrInternalError = errors.New("internal server error")
	ErrRequestError  = errors.New("request error")
)

type HTTPClient struct {
	client     *http.Client
	baseURL    string
	maxRetries int
	retryDelay time.Duration
	logger     *zap.Logger
}

type HTTPClientConfig struct {
	MaxRetries int
	RetryDelay time.Duration
	Timeout    time.Duration
}

func (c *HTTPClientConfig) normalize() {
	if c.Timeout == 0 {
		c.Timeout = time.Second * 5
	}

	if c.MaxRetries == 0 {
		c.MaxRetries = 3
	}

	if c.RetryDelay == 0 {
		c.RetryDelay = 100 * time.Millisecond
	}
}

func NewHTTPClient(baseURL string, cfg HTTPClientConfig, logger *zap.Logger) *HTTPClient {
	cfg.normalize()

	return &HTTPClient{
		client: &http.Client{
			Timeout: cfg.Timeout,
		},
		baseURL:    baseURL,
		maxRetries: cfg.MaxRetries,
		retryDelay: cfg.RetryDelay,
		logger:     logger,
	}
}

func (c *HTTPClient) Get(ctx context.Context, path string, headers map[string]string) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return c.do(req)
}

func (c *HTTPClient) Post(ctx context.Context, path string, body interface{}, headers map[string]string) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+path, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return c.do(req)
}

func (c *HTTPClient) do(req *http.Request) (*http.Response, error) {
	var (
		attempt int
		resp    *http.Response
		delay   = c.retryDelay
		err     error
	)
	for attempt < c.maxRetries {
		resp, err = c.client.Do(req)
		if err != nil && !isTimeout(err) {
			return nil, err
		}

		if err == nil && (resp.StatusCode == http.StatusOK || resp.StatusCode < http.StatusInternalServerError) {
			return resp, nil
		}

		attempt++

		c.logger.Warn("http request retry",
			zap.String("url", req.URL.String()),
			zap.Int("attempt", attempt),
			zap.Error(err),
		)

		time.Sleep(delay)

		delay *= 2
	}

	return resp, err
}

func isTimeout(err error) bool {
	var netErr net.Error
	if errors.As(err, &netErr) && (netErr.Timeout() || netErr.Temporary()) {
		return true
	}

	return false
}
