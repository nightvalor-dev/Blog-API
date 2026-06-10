package media

import (
	"Project2-v7/internal/shared/middleware"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *MediaHandler, rl *middleware.RedisRateLimiter) {
	userLimiter := rl.LimitByKey(func(r *http.Request) string {
		claims, ok := middleware.GetClaimsFromContext(r.Context())
		if !ok {
			return fmt.Sprintf("rl:ip:%s", r.RemoteAddr)
		}
		return fmt.Sprintf("rl:user:%d", claims.UserId)
	})

	r.Route("/blogs/{id}/media", func(r chi.Router) {
		r.Use(middleware.RequireAuth)
		r.With(userLimiter).Post("/", handler.Upload)
		r.With(rl.LimitByIP()).Get("/", handler.GetByBlogId)
		r.With(userLimiter).Delete("/{mid}", handler.Delete)
	})
}
