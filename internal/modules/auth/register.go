package auth

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required,len=10,numeric"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}
