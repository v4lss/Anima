// Package middleware - Auth validates the JWT in the Authorization header.
// It injects the userID into the request context on success.
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/v4lss/animas/pkg/jwt"
	"github.com/v4lss/animas/pkg/response"
)

type contextKey string

const UserIDKey contextKey = "userID"

// Auth returns an HTTP middleware that guards routes with JWT verification.
func Auth(jwtSvc *jwt.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				response.Error(w, http.StatusUnauthorized, "missing token")
				return
			}

			token := strings.TrimPrefix(header, "Bearer ")
			userID, err := jwtSvc.Verify(token)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "invalid token")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UserIDFromContext extracts the userID injected by the Auth middleware.
func UserIDFromContext(ctx context.Context) string {
	id, _ := ctx.Value(UserIDKey).(string)
	return id
}
