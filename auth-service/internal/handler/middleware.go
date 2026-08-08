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

		userID, err := h.userService.ValidateToken(splitToken[1])
		if err != nil {
			writeError(w, r, http.StatusUnauthorized, "Invalid or expired token")

			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
