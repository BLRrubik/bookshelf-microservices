package handler

import (
	"context"
	"net/http"
	"strings"
)

func (h *AuthHandler) AuthMiddleware(next http.Handler) http.Handler {
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

		res, err := h.userService.ValidateToken(splitToken[1])
		if err != nil {
			writeError(w, r, http.StatusUnauthorized, "Invalid or expired token")

			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, res.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

const serviceKeyHeader = "X-Service-Key"

func ServiceKeyMiddleware(expectedKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			serviceKey := r.Header.Get(serviceKeyHeader)

			if serviceKey == "" {
				writeError(w, r, http.StatusUnauthorized, "missing service key")

				return
			}

			if serviceKey != expectedKey {
				writeError(w, r, http.StatusForbidden, "invalid service key")

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
