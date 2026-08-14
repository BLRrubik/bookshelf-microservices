package handler

import (
	"context"
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
