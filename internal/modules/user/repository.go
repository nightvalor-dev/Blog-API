package user

import (
	"Project2-v7/internal/shared/roles"
	"context"
)

type UserRepository interface {
	Create(user CreateUserRequest) error
	GetAll() ([]UserResponse, error)
	GetById(id int) (UserResponse, error)
	Update(id int, newUser UpdateUserRequest) error
	Delete(id int) error

	GetUserByEmail(email string) (UserAuthDetails, error)
	UpdateRole(ctx context.Context, id int, role roles.Role) error
}
