package user

import "context"

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(req CreateUserRequest) error {
	return s.repo.Create(req)
}

func (s *UserService) GetAll() ([]UserResponse, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserService) GetById(id int) (UserResponse, error) {
	user, err := s.repo.GetById(id)
	if err != nil {
		return UserResponse{}, err
	}
	return user, nil
}

func (s *UserService) Update(id int, req UpdateUserRequest) error {
	return s.repo.Update(id, req)
}

func (s *UserService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *UserService) ChangeRole(ctx context.Context, id int, req ChangeRoleRequest) error {
	return s.repo.UpdateRole(ctx, id, req.Role)
}
