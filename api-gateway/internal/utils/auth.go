package utils

import (
	"net/http"
	"strings"
)

const bearerPrefix = "Bearer "

// ExtractBearerToken returns the token from the Authorization header,
// or "" if the header is missing or not a Bearer token.
func ExtractBearerToken(r *http.Request) string {
	authHdr := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHdr, bearerPrefix) {
		return ""
	}

	return strings.TrimPrefix(authHdr, bearerPrefix)
}
