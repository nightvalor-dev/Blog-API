package comment

import "time"

type Comment struct {
	CommentId int       `db:"comment_id" json:"comment_id"`
	Content   string    `db:"content" json:"content"`
	BlogId    int       `db:"blog_id" json:"blog_id"`
	UserId    int       `db:"user_id" json:"user_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
