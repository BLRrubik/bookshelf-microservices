package middleware

import (
	"net/http"

	"github.com/google/uuid"
)

func RequestID() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-Request-Id") == "" {
				r.Header.Set("X-Request-Id", uuid.NewString())
			}

			next.ServeHTTP(w, r)
		})
	}
}
