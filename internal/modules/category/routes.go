package category

import (
	"Project2-v7/internal/shared/middleware"
	"Project2-v7/internal/shared/roles"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *CategoryHandler) {
	r.Route("/categories", func(r chi.Router) {
		// Public
		r.Get("/", handler.GetAllCategory)
		r.Get("/{id}", handler.GetCategoryById)

		// Admin/DBA only
		r.With(middleware.RequireAuth, middleware.RequireRole(roles.RoleAdmin, roles.RoleDBA)).Post("/", handler.CreateCategory)
		r.With(middleware.RequireAuth, middleware.RequireRole(roles.RoleAdmin, roles.RoleDBA)).Put("/{id}", handler.UpdateCategory)
		r.With(middleware.RequireAuth, middleware.RequireRole(roles.RoleAdmin, roles.RoleDBA)).Delete("/{id}", handler.DeleteCategory)
	})
}
