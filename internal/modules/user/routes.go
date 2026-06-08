package user

import (
	"Project2-v7/internal/shared/middleware"
	"Project2-v7/internal/shared/roles"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, h *UserHandler) {
	r.Route("/users", func(r chi.Router) {
		r.With(middleware.RequireRole(roles.RoleAdmin, roles.RoleDBA)).Get("/", h.GetAllUsers)
		r.With(middleware.RequireRole(roles.RoleAdmin)).Post("/", h.CreateUser)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetUserById)
			r.Put("/", h.UpdateUser)
			r.With(middleware.RequireRole(roles.RoleAdmin, roles.RoleDBA)).Delete("/", h.DeleteUser)
			r.With(middleware.RequireRole(roles.RoleAdmin)).Patch("/role", h.ChangeUserRole)
		})
	})
}
