package comment

import (
	"Project2-v7/internal/shared/middleware"
	"Project2-v7/internal/shared/roles"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *CommentHandler) {
	r.Route("/comments", func(r chi.Router) {
		r.Get("/", handler.GetAllComments)
		r.Post("/", handler.CreateComment)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handler.GetCommentById)
			r.Put("/", handler.UpdateComment)
			r.With(middleware.RequireRole(roles.RoleAdmin)).Delete("/", handler.DeleteComment)
		})
	})
}
