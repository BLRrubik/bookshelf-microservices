package handler

import (
	"bookshelf/books-service/internal/client"
	"context"
	"errors"
	"net/http"
	"strings"
)

func AuthMiddleware(authClient *client.AuthClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if len(authHeader) == 0 {
				writeError(w, r, http.StatusUnauthorized, "Authorization header required")

				return
			}

			splitToken := strings.Split(authHeader, " ")
			if strings.ToLower(splitToken[0]) != "bearer" {
				writeError(w, r, http.StatusUnauthorized, "Invalid authorization header format")

				return
			}

			resp, err := authClient.VerifyToken(r.Context(), splitToken[1])
			if err != nil {

				switch {
				case errors.Is(err, client.ErrRequestError):
					writeError(w, r, http.StatusUnauthorized, "Invalid or expired token")
				default:
					writeError(w, r, http.StatusServiceUnavailable, "Service unavailable")
				}

				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, resp.UserID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
