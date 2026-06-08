package user

import "Project2-v7/internal/shared/roles"

type CreateUserRequest struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Bio          string `json:"bio"`
	PasswordHash string `json:"password_hash"`
}

type UpdateUserRequest struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Bio          string `json:"bio"`
	PasswordHash string `json:"password_hash"`
}

type UserResponse struct {
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	Phone        string     `json:"phone"`
	Bio          string     `json:"bio"`
	AssignedRole roles.Role `json:"role"`
}

type UserAuthDetails struct {
	Id           int        `json:"id"`
	PasswordHash string     `json:"password_hash"`
	AssignedRole roles.Role `json:"assigned_role"`
}

type ChangeRoleRequest struct {
	Role roles.Role `json:"role" validate:"required,oneof=admin user premium_user dba"`
}
