package likes

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *LikeHandler) {
	r.Post("/blogs/{id}/like", handler.Like)
	r.Delete("/blogs/{id}/like", handler.Unlike)
	r.Get("/blogs/{id}/likes", handler.GetCount)
}
