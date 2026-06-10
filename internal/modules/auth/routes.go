package auth

import (
	"Project2-v7/internal/shared/middleware"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *AuthHandler, rl *middleware.RedisRateLimiter) {
	// Tight limit on auth routes to prevent brute force
	authLimiter := rl.LimitByKey(func(r *http.Request) string {
		return fmt.Sprintf("rl:ip:%s:auth", r.RemoteAddr)
	})

	r.Route("/auth", func(r chi.Router) {
		r.With(authLimiter).Post("/register", handler.Register)
		r.With(authLimiter).Post("/login", handler.Login)
	})
}
