package tag

import "time"

type Tag struct {
	TagId     int       `json:"tag_id" db:"tag_id"`
	TagName   string    `json:"tag_name" db:"tag_name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
