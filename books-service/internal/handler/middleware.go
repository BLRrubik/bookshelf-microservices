package handler

import (
	"bookshelf/books-service/internal/client"
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func RequestLogger(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			fields := []zap.Field{
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", ww.Status()),
				zap.Duration("duration", time.Since(start)),
				zap.String("request_id", middleware.GetReqID(r.Context())),
				zap.String("remote_addr", r.RemoteAddr),
			}

			switch {
			case ww.Status() >= 500:
				logger.Error("http request", fields...)
			case ww.Status() >= 400:
				logger.Warn("http request", fields...)
			default:
				logger.Info("http request", fields...)
			}
		})
	}
}

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
