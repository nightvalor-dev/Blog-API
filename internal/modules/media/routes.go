package media

import (
	"Project2-v7/internal/shared/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *MediaHandler) {
	r.Route("/blogs/{id}/media", func(r chi.Router) {
		r.Use(middleware.RequireAuth)
		r.Post("/", handler.Upload)
		r.Get("/", handler.GetByBlogId)
		r.Delete("/{mid}", handler.Delete)
	})
}
