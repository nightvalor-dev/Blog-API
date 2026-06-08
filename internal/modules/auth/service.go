package auth

import (
	"Project2-v7/internal/modules/user"
	"Project2-v7/pkg/jwt"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo user.UserRepository
}

func NewAuthService(repo user.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Register(req RegisterRequest) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	registeredUser := user.CreateUserRequest{
		Username:     req.Username,
		Email:        req.Email,
		Phone:        req.Phone,
		PasswordHash: string(hash),
	}

	return s.repo.Create(registeredUser)
}

func (s *AuthService) Login(req LoginRequest) (LoginResponse, error) {
	u, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return LoginResponse{}, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password))
	if err != nil {
		return LoginResponse{}, errors.New("invalid credentials")
	}

	token, err := jwt.GenerateToken(u.Id, u.AssignedRole)
	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{Token: token}, nil
}
