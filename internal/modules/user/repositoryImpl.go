package user

import (
	"Project2-v7/internal/shared/roles"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (ur *userRepository) Create(user CreateUserRequest) error {
	query := `INSERT INTO users (username, email, phone, bio, password_hash, role) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := ur.db.Exec(context.Background(), query,
		user.Username, user.Email, user.Phone,
		user.Bio, user.PasswordHash, roles.RoleUser)
	return err
}

func (ur *userRepository) GetAll() ([]UserResponse, error) {
	query := `SELECT username, email, phone, bio, role FROM users`
	rows, err := ur.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []UserResponse
	for rows.Next() {
		var entity UserResponse
		err = rows.Scan(&entity.Username, &entity.Email,
			&entity.Phone, &entity.Bio, &entity.AssignedRole)

		if err != nil {
			return nil, err
		}

		result = append(result, entity)
	}
	return result, nil
}

func (ur *userRepository) GetById(id int) (UserResponse, error) {
	var entity UserResponse
	query := `SELECT username, email, phone, bio, role FROM users WHERE user_id = $1`
	err := ur.db.QueryRow(context.Background(), query, id).Scan(&entity.Username, &entity.Email,
		&entity.Phone, &entity.Bio, &entity.AssignedRole)
	if err != nil {
		return UserResponse{}, err
	}

	return entity, nil
}

func (ur *userRepository) Update(id int, newUser UpdateUserRequest) error {
	args := []any{}
	argIdx := 1
	query := "UPDATE users SET "
	sep := ""

	if newUser.Username != "" {
		query += fmt.Sprintf("%s username = $%d", sep, argIdx)
		args = append(args, newUser.Username)
		argIdx++
		sep = ","
	}

	if newUser.Email != "" {
		query += fmt.Sprintf("%s email = $%d", sep, argIdx)
		args = append(args, newUser.Email)
		argIdx++
		sep = ","
	}

	if newUser.Phone != "" {
		query += fmt.Sprintf("%s phone = $%d", sep, argIdx)
		args = append(args, newUser.Phone)
		argIdx++
		sep = ","
	}

	if newUser.Bio != "" {
		query += fmt.Sprintf("%s bio = $%d", sep, argIdx)
		args = append(args, newUser.Bio)
		argIdx++
		sep = ","
	}

	if newUser.PasswordHash != "" {
		query += fmt.Sprintf("%s password_hash = $%d", sep, argIdx)
		args = append(args, newUser.PasswordHash)
		argIdx++
		sep = ","
	}

	if len(args) == 0 {
		return nil
	}

	query += fmt.Sprintf(" WHERE user_id = $%d", argIdx)
	args = append(args, id)

	_, err := ur.db.Exec(context.Background(), query, args...)
	return err
}

func (ur *userRepository) Delete(id int) error {
	query := "DELETE FROM users WHERE user_id = $1"
	_, err := ur.db.Exec(context.Background(), query, id)
	return err
}

func (ur *userRepository) GetUserByEmail(email string) (UserAuthDetails, error) {
	var entity UserAuthDetails
	query := `SELECT user_id, password_hash, role FROM users WHERE email = $1`
	err := ur.db.QueryRow(context.Background(), query, email).Scan(
		&entity.Id, &entity.PasswordHash, &entity.AssignedRole,
	)
	if err != nil {
		return UserAuthDetails{}, err
	}

	return entity, nil
}

func (r *userRepository) UpdateRole(ctx context.Context, id int, role roles.Role) error {
	query := `UPDATE users SET role = $1, updated_at = NOW() WHERE user_id = $2`
	result, err := r.db.Exec(ctx, query, role, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
