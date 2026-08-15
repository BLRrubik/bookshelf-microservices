package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const requestIDKey = "requestIDKey"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-Id")
		if requestID == "" {
			requestID = uuid.NewString()
		}

		r.Header.Set("X-Request-Id", requestID)
		w.Header().Set("X-Request-Id", requestID)

		ctx := context.WithValue(r.Context(), requestIDKey, requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetRequestID(ctx context.Context) string {
	requestID, _ := ctx.Value(requestIDKey).(string)

	return requestID
}
