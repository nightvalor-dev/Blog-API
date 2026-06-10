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
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(
	lo *middleware.Logger,
	authHandler *auth.AuthHandler,
	blogHandler *blog.BlogHandler,
	categoryHandler *category.CategoryHandler,
	tagHandler *tag.TagHandler,
	commentHandler *comment.CommentHandler,
	userHandler *user.UserHandler,
	likeHandler *likes.LikeHandler,
	mediaHandler *media.MediaHandler,
	authLimiter *middleware.RedisRateLimiter, // tight: 10 req/min for login & register
	generalLimiter *middleware.RedisRateLimiter, // normal: 100 req/min for everything else
) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.LoggerMiddleware(lo))
	r.Use(chiMiddleware.Recoverer)

	auth.RegisterRoutes(r, authHandler, authLimiter)
	blog.RegisterRoutes(r, blogHandler, generalLimiter)
	category.RegisterRoutes(r, categoryHandler, generalLimiter)
	tag.RegisterRoutes(r, tagHandler, generalLimiter)

	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth)
		comment.RegisterRoutes(r, commentHandler, generalLimiter)
		user.RegisterRoutes(r, userHandler, generalLimiter)
		likes.RegisterRoutes(r, likeHandler, generalLimiter)
		media.RegisterRoutes(r, mediaHandler, generalLimiter)
	})

	return r
}
