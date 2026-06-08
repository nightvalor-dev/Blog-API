package user

import (
	"Project2-v7/internal/shared/roles"
	"time"
)

type User struct {
	UserId       int        `db:"user_id"       json:"user_id"`
	Username     string     `db:"username"      json:"username"`
	Email        string     `db:"email"         json:"email"`
	Phone        string     `db:"phone"         json:"phone"`
	Bio          string     `db:"bio"           json:"bio"`
	PasswordHash string     `db:"password_hash" json:"-"`
	AssignedRole roles.Role `db:"role"          json:"role"`
	CreatedAt    time.Time  `db:"created_at"    json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"    json:"updated_at"`
}
