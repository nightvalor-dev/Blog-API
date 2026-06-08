package main

import (
	"Project2-v7/api"
	"Project2-v7/config"
	"Project2-v7/internal/modules/auth"
	"Project2-v7/internal/modules/blog"
	"Project2-v7/internal/modules/category"
	"Project2-v7/internal/modules/comment"
	"Project2-v7/internal/modules/likes"
	"Project2-v7/internal/modules/media"
	"Project2-v7/internal/modules/tag"
	"Project2-v7/internal/modules/user"
	cloudinary "Project2-v7/internal/shared/cloudinary"
	db "Project2-v7/internal/shared/db"
	"Project2-v7/internal/shared/middleware/logger"
	"Project2-v7/internal/shared/redis"
	"fmt"
	"log"
	"net/http"
)

func main() {
	cfg := config.Load()

	lo, err := logger.NewLogger("logs/app.log")
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer lo.Close()

	pool, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	redisClient, err := redis.NewRedisClient(cfg.RedisAddr)
	if err != nil {
		log.Fatalf("redis: %v", err)
	}
	appCache := redis.NewCache(redisClient)

	cldClient, err := cloudinary.NewClient(cfg)
	if err != nil {
		log.Fatalf("cloudinary: %v", err)
	}

	cloudinaryService := cloudinary.NewService(cldClient)

	// repositories
	userRepo := user.NewUserRepository(pool)
	blogRepo := blog.NewBlogRepository(pool)
	categoryRepo := category.NewCatRepository(pool)
	tagRepo := tag.NewTagRepository(pool)
	commentRepo := comment.NewCommentRepository(pool)
	likeRepo := likes.NewLikeRepository(pool)
	mediaRepo := media.NewMediaRepository(pool)

	// services
	authService := auth.NewAuthService(userRepo)
	blogService := blog.NewBlogService(blogRepo, categoryRepo, tagRepo, mediaRepo, appCache)
	categoryService := category.NewCategoryService(categoryRepo)
	tagService := tag.NewTagService(tagRepo)
	commentService := comment.NewCommentService(commentRepo)
	userService := user.NewUserService(userRepo)
	likeService := likes.NewLikeService(likeRepo)
	mediaService := media.NewMediaService(mediaRepo, cloudinaryService)

	// handlers
	authHandler := auth.NewAuthHandler(authService)
	blogHandler := blog.NewBlogHandler(blogService)
	categoryHandler := category.NewCategoryHandler(categoryService)
	tagHandler := tag.NewTagHandler(tagService)
	commentHandler := comment.NewCommentHandler(commentService)
	userHandler := user.NewUserHandler(userService)
	likeHandler := likes.NewLikeHandler(likeService)
	mediaHandler := media.NewMediaHandler(mediaService)

	router := api.NewRouter(lo, authHandler, blogHandler, categoryHandler, tagHandler,
		commentHandler, userHandler, likeHandler, mediaHandler)

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("server running on %s", addr)
	if err = http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
