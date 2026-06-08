package blog

import (
	"Project2-v7/internal/shared/middleware"
	"Project2-v7/internal/shared/roles"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *BlogHandler) {
	r.Route("/blogs", func(r chi.Router) {
		r.Get("/", handler.GetAllBlogs)
		r.Get("/{id}", handler.GetBlogById)

		r.With(middleware.RequireAuth).Post("/", handler.CreateBlog)
		r.With(middleware.RequireAuth, middleware.RequireRole(roles.RoleAdmin)).Put("/{id}", handler.UpdateBlog)
		r.With(middleware.RequireAuth, middleware.RequireRole(roles.RoleAdmin)).Delete("/{id}", handler.DeleteBlog)
	})
}
