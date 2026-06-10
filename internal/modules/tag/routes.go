package tag

import (
	"Project2-v7/internal/shared/middleware"
	"Project2-v7/internal/shared/roles"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *TagHandler, rl *middleware.RedisRateLimiter) {
	userLimiter := rl.LimitByKey(func(r *http.Request) string {
		claims, ok := middleware.GetClaimsFromContext(r.Context())
		if !ok {
			return fmt.Sprintf("rl:ip:%s", r.RemoteAddr)
		}
		return fmt.Sprintf("rl:user:%d", claims.UserId)
	})

	r.Route("/tags", func(r chi.Router) {
		r.With(rl.LimitByIP()).Get("/", handler.GetAllTags)
		r.With(rl.LimitByIP()).Get("/{id}", handler.GetTagById)

		r.With(middleware.RequireAuth, middleware.RequireRole(roles.RoleAdmin, roles.RoleDBA), userLimiter).Post("/", handler.CreateTag)
		r.With(middleware.RequireAuth, middleware.RequireRole(roles.RoleAdmin, roles.RoleDBA), userLimiter).Put("/{id}", handler.UpdateTag)
		r.With(middleware.RequireAuth, middleware.RequireRole(roles.RoleAdmin, roles.RoleDBA), userLimiter).Delete("/{id}", handler.DeleteTag)
	})
}
