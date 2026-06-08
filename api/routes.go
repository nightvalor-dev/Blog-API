package api

import (
	"Project2-v7/internal/modules/auth"
	"Project2-v7/internal/modules/blog"
	"Project2-v7/internal/modules/category"
	"Project2-v7/internal/modules/comment"
	"Project2-v7/internal/modules/likes"
	"Project2-v7/internal/modules/media"
	"Project2-v7/internal/modules/tag"
	"Project2-v7/internal/modules/user"
	"Project2-v7/internal/shared/middleware"
	"Project2-v7/internal/shared/middleware/logger"
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(
	lo *logger.Logger,
	authHandler *auth.AuthHandler,
	blogHandler *blog.BlogHandler,
	categoryHandler *category.CategoryHandler,
	tagHandler *tag.TagHandler,
	commentHandler *comment.CommentHandler,
	userHandler *user.UserHandler,
	likeHandler *likes.LikeHandler,
	mediaHandler *media.MediaHandler, // ← add
) http.Handler {
	r := chi.NewRouter()
	r.Use(logger.LoggerMiddleware(lo))
	r.Use(chiMiddleware.Recoverer)

	auth.RegisterRoutes(r, authHandler)
	blog.RegisterRoutes(r, blogHandler)
	category.RegisterRoutes(r, categoryHandler)
	tag.RegisterRoutes(r, tagHandler)

	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth)
		comment.RegisterRoutes(r, commentHandler)
		user.RegisterRoutes(r, userHandler)
		likes.RegisterRoutes(r, likeHandler)
		media.RegisterRoutes(r, mediaHandler)
	})

	return r
}
