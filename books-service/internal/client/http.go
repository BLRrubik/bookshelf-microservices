package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"syscall"
	"time"
)

type HTTPClient struct {
	client  *http.Client
	baseURL string
}

func NewHTTPClient(baseURL string, timeout time.Duration) *HTTPClient {
	if timeout == 0 {
		timeout = time.Second * 5
	}

	return &HTTPClient{
		client: &http.Client{
			Timeout: timeout,
		},
		baseURL: baseURL,
	}
}

func (c *HTTPClient) Get(ctx context.Context, path string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, errors.New("request timeout")
		}

		if errors.Is(err, syscall.ECONNREFUSED) {
			return nil, errors.New("connection refused")
		}

		return nil, err
	}

	if !successRequest(resp.StatusCode) {
		return nil, fmt.Errorf("request failed with status code %d", resp.StatusCode)
	}

	return resp, nil
}

func (c *HTTPClient) Post(ctx context.Context, path string, body interface{}, headers map[string]string) (*http.Response, error) {
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

	resp, err := c.client.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, errors.New("request timeout")
		}

		if errors.Is(err, syscall.ECONNREFUSED) {
			return nil, errors.New("connection refused")
		}

		return nil, err
	}

	if !successRequest(resp.StatusCode) {
		return nil, fmt.Errorf("request failed with status code %d", resp.StatusCode)
	}

	return resp, nil
}

func successRequest(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}
