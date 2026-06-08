package middleware

import (
	"Project2-v7/internal/shared/roles"
	"Project2-v7/pkg/jwt"
	"Project2-v7/pkg/utils"
	"context"
	"errors"
	"net/http"
	"strings"
)

type contextKey string

const ClaimsKey contextKey = "claims"

func extractClaims(r *http.Request) (*jwt.Claims, error) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, errors.New("missing or invalid authorization header")
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	return jwt.ValidateToken(tokenStr)
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := extractClaims(r)
		if err != nil {
			utils.WriteJSON(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		ctx := context.WithValue(r.Context(), ClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireRole(allowedRoles ...roles.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(ClaimsKey).(*jwt.Claims)
			if !ok || claims == nil {
				utils.WriteJSON(w, http.StatusUnauthorized, "unauthorized")
				return
			}
			userRole := roles.Role(claims.Role)
			for _, role := range allowedRoles {
				if userRole == role {
					next.ServeHTTP(w, r)
					return
				}
			}
			utils.WriteJSON(w, http.StatusForbidden, "forbidden")
		})
	}
}

func GetClaimsFromContext(ctx context.Context) (*jwt.Claims, bool) {
	claims, ok := ctx.Value(ClaimsKey).(*jwt.Claims)
	return claims, ok
}
