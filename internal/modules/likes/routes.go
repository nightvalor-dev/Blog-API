package likes

import (
	"Project2-v7/internal/shared/middleware"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *LikeHandler, rl *middleware.RedisRateLimiter) {
	userLimiter := rl.LimitByKey(func(r *http.Request) string {
		claims, ok := middleware.GetClaimsFromContext(r.Context())
		if !ok {
			return fmt.Sprintf("rl:ip:%s", r.RemoteAddr)
		}
		return fmt.Sprintf("rl:user:%d", claims.UserId)
	})

	r.With(middleware.RequireAuth, userLimiter).Post("/blogs/{id}/like", handler.Like)
	r.With(middleware.RequireAuth, userLimiter).Delete("/blogs/{id}/like", handler.Unlike)
	r.With(rl.LimitByIP()).Get("/blogs/{id}/likes", handler.GetCount)
}
