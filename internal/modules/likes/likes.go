package likes

import "time"

type Likes struct {
	BlogId  int       `json:"blog_id" db:"blog_id"`
	UserId  int       `json:"user_id" db:"user_id"`
	LikedAt time.Time `json:"liked_at" db:"liked_at"`
}
