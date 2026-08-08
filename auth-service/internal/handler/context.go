package handler

import "context"

type contextKey string

const (
	userIDKey    contextKey = "userID"
	requestIDKey contextKey = "requestID"
)

func getUserID(ctx context.Context) string {
	return ctx.Value(userIDKey).(string)
}
