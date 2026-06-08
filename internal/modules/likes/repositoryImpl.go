package likes

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type likeRepository struct {
	db *pgxpool.Pool
}

func NewLikeRepository(db *pgxpool.Pool) LikeRepository {
	return &likeRepository{db: db}
}

func (l *likeRepository) Add(blogId, userId int) error {
	query := "INSERT INTO blog_likes (blog_id, user_id) VALUES ($1, $2)"
	_, err := l.db.Exec(context.Background(), query, blogId, userId)
	return err
}

func (l *likeRepository) Remove(blogId, userId int) error {
	query := "DELETE FROM blog_likes WHERE blog_id = $1 AND user_id = $2"
	_, err := l.db.Exec(context.Background(), query, blogId, userId)
	return err
}

func (l *likeRepository) Count(blogId int) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM blog_likes WHERE blog_id = $1"
	err := l.db.QueryRow(context.Background(), query, blogId).Scan(&count)
	return count, err
}

func (l *likeRepository) Exists(blogId, userId int) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM blog_likes WHERE blog_id = $1 AND user_id = $2)"
	err := l.db.QueryRow(context.Background(), query, blogId, userId).Scan(&exists)
	return exists, err
}
