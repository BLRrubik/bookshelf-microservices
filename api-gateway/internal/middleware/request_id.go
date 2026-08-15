package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

func RequestID() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = uuid.NewString()
			}

			r.Header.Set("X-Request-Id", requestID)
			ctx := context.WithValue(r.Context(), "request_id", requestID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
