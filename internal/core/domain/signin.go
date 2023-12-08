package domain

type SigninRequest struct {
	Email    string `json:"email" validate:"required,email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SigninResponse struct {
	AccessToken string `json:"accessToken"`
}
