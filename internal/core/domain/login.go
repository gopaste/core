package domain

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
}

type LoginService interface {
	GetUserByEmail(email string) (*User, error)
}
