package comment

import (
	"Project2-v7/internal/shared/middleware"
	"Project2-v7/internal/shared/roles"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *CommentHandler, rl *middleware.RedisRateLimiter) {
	userLimiter := rl.LimitByKey(func(r *http.Request) string {
		claims, ok := middleware.GetClaimsFromContext(r.Context())
		if !ok {
			return fmt.Sprintf("rl:ip:%s", r.RemoteAddr)
		}
		return fmt.Sprintf("rl:user:%d", claims.UserId)
	})

	r.Route("/comments", func(r chi.Router) {
		r.With(rl.LimitByIP()).Get("/", handler.GetAllComments)
		r.With(middleware.RequireAuth, userLimiter).Post("/", handler.CreateComment)

		r.Route("/{id}", func(r chi.Router) {
			r.With(rl.LimitByIP()).Get("/", handler.GetCommentById)
			r.With(middleware.RequireAuth, userLimiter).Put("/", handler.UpdateComment)
			r.With(middleware.RequireAuth, middleware.RequireRole(roles.RoleAdmin), userLimiter).Delete("/", handler.DeleteComment)
		})
	})
}
