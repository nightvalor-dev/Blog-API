package media

import "time"

type Media struct {
	MediaId   int       `json:"media_id" db:"media_id"`
	BlogId    int       `json:"blog_id" db:"blog_id"`
	Url       string    `json:"url" db:"url"`
	PublicId  string    `json:"public_id" db:"public_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
