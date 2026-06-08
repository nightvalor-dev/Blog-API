package category

import "time"

type Category struct {
	CategoryId   int       `json:"category_id" db:"category_id"`
	CategoryName string    `json:"category_name" db:"category_name"`
	Description  string    `json:"description" db:"description"`
	CreatedAt    time.Time `json:"created_at"db:"created_at"`
}
