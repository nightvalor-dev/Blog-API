package blog

import (
	"Project2-v7/internal/shared/middleware"
	"Project2-v7/internal/shared/roles"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *BlogHandler, rl *middleware.RedisRateLimiter) {
	userLimiter := rl.LimitByKey(func(r *http.Request) string {
		claims, ok := middleware.GetClaimsFromContext(r.Context())
		if !ok {
			return fmt.Sprintf("rl:ip:%s", r.RemoteAddr)
		}
		return fmt.Sprintf("rl:user:%d", claims.UserId)
	})

	r.Route("/blogs", func(r chi.Router) {
		// public routes — rate limit by IP
		r.With(rl.LimitByIP()).Get("/", handler.GetAllBlogs)
		r.With(rl.LimitByIP()).Get("/{id}", handler.GetBlogById)

		// write routes — rate limit by user ID after auth resolves
		r.With(middleware.RequireAuth, userLimiter).Post("/", handler.CreateBlog)
		r.With(middleware.RequireAuth, middleware.RequireRole(roles.RoleAdmin), userLimiter).Put("/{id}", handler.UpdateBlog)
		r.With(middleware.RequireAuth, middleware.RequireRole(roles.RoleAdmin), userLimiter).Delete("/{id}", handler.DeleteBlog)
	})
}
