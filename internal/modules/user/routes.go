package user

import (
	"Project2-v7/internal/shared/middleware"
	"Project2-v7/internal/shared/roles"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, h *UserHandler, rl *middleware.RedisRateLimiter) {
	userLimiter := rl.LimitByKey(func(r *http.Request) string {
		claims, ok := middleware.GetClaimsFromContext(r.Context())
		if !ok {
			return fmt.Sprintf("rl:ip:%s", r.RemoteAddr)
		}
		return fmt.Sprintf("rl:user:%d", claims.UserId)
	})

	r.Route("/users", func(r chi.Router) {
		r.With(middleware.RequireAuth, middleware.RequireRole(roles.RoleAdmin, roles.RoleDBA), userLimiter).Get("/", h.GetAllUsers)
		r.With(middleware.RequireAuth, middleware.RequireRole(roles.RoleAdmin), userLimiter).Post("/", h.CreateUser)

		r.Route("/{id}", func(r chi.Router) {
			r.With(rl.LimitByIP()).Get("/", h.GetUserById)
			r.With(middleware.RequireAuth, userLimiter).Put("/", h.UpdateUser)
			r.With(middleware.RequireAuth, middleware.RequireRole(roles.RoleAdmin, roles.RoleDBA), userLimiter).Delete("/", h.DeleteUser)
			r.With(middleware.RequireAuth, middleware.RequireRole(roles.RoleAdmin), userLimiter).Patch("/role", h.ChangeUserRole)
		})
	})
}
