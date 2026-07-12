package auth

import (
	"context"
	"net/http"
)

// contextKey namespaces values stored on the request context by this package.
type contextKey int

const userContextKey contextKey = iota

// Middleware extracts and validates the Bearer token, attaching the authenticated user to
// the request context, or responding 401 if missing/invalid.
func Middleware(sessions *SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: implement
			next.ServeHTTP(w, r)
		})
	}
}

// UserFromContext retrieves the authenticated user attached by Middleware.
func UserFromContext(ctx context.Context) (userID string, ok bool) {
	// TODO: implement
	return "", false
}
