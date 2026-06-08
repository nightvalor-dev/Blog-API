package blog

import "time"

type Blog struct {
	BlogId      int        `json:"blog_id" db:"blog_id"`
	Title       string     `json:"title" db:"title"`
	Content     string     `json:"content" db:"content"`
	Status      BlogStatus `json:"status" db:"status"`
	UserId      int        `json:"user_id" db:"user_id"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	PublishedAt time.Time  `json:"published_at" db:"published_at"`
}

// if category already exists then return the categoryId for this table
// if the category doesn't exist then make a new one and return the categoryId for this table
