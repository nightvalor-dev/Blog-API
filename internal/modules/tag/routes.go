package tag

import (
	"Project2-v7/internal/shared/middleware"
	"Project2-v7/internal/shared/roles"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *TagHandler) {
	r.Route("/tags", func(r chi.Router) {
		// Public
		r.Get("/", handler.GetAllTags)
		r.Get("/{id}", handler.GetTagById)

		// Admin/DBA only
		r.With(middleware.RequireAuth, middleware.RequireRole(roles.RoleAdmin, roles.RoleDBA)).Post("/", handler.CreateTag)
		r.With(middleware.RequireAuth, middleware.RequireRole(roles.RoleAdmin, roles.RoleDBA)).Put("/{id}", handler.UpdateTag)
		r.With(middleware.RequireAuth, middleware.RequireRole(roles.RoleAdmin, roles.RoleDBA)).Delete("/{id}", handler.DeleteTag)
	})
}
