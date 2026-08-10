package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type VerifyRequest struct {
	Token string `json:"token"`
}

type VerifyResponse struct {
	Valid     bool   `json:"valid"`
	UserID    string `json:"user_id"`
	ExpiresAt string `json:"expires_at"`
	Error     string `json:"error,omitempty"`
}

type GetUsersRequest struct {
	IDs []string `json:"ids"`
}

type HealthResponse struct {
	Status string `json:"status"`
}

type UserPublic struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AuthClient struct {
	httpClient *HTTPClient
	serviceKey string
}

func NewAuthClient(baseURL string, timeout time.Duration, maxRetries int, retryDelay time.Duration, serviceKey string) *AuthClient {
	return &AuthClient{
		httpClient: NewHTTPClient(baseURL, HTTPClientConfig{
			Timeout:    timeout,
			MaxRetries: maxRetries,
			RetryDelay: retryDelay,
		}),
		serviceKey: serviceKey,
	}
}

func (c *AuthClient) VerifyToken(ctx context.Context, token string) (*VerifyResponse, error) {
	headers := map[string]string{
		"Content-Type":  "application/json",
		"X-Service-Key": c.serviceKey,
	}

	req := VerifyRequest{
		Token: token,
	}

	resp, err := c.httpClient.Post(ctx, "/internal/v1/auth/verify", req, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch {
	case resp.StatusCode >= 500:
		return nil, ErrInternalError
	case resp.StatusCode >= 400:
		return nil, ErrRequestError
	}

	var response VerifyResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *AuthClient) GetUsersByIDs(ctx context.Context, ids []string) ([]UserPublic, error) {
	headers := map[string]string{
		"Content-Type":  "application/json",
		"X-Service-Key": c.serviceKey,
	}

	req := GetUsersRequest{
		IDs: ids,
	}

	resp, err := c.httpClient.Post(ctx, "/internal/v1/users/batch", req, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch {
	case resp.StatusCode >= 500:
		return nil, ErrInternalError
	case resp.StatusCode >= 400:
		return nil, ErrRequestError
	}

	var response []UserPublic
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response, nil
}

func (c *AuthClient) Health(ctx context.Context) (*HealthResponse, error) {
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	resp, err := c.httpClient.Get(ctx, "/health", headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected status: " + resp.Status)
	}

	var response HealthResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}
