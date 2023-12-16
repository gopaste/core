package entity

import (
	"time"

	"github.com/google/uuid"
)

type SigninRequest struct {
	Email    string `json:"email" validate:"required,email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ResetPasswordRequest struct {
	Password             string `json:"password" binding:"required"`
	PasswordConfirmation string `json:"passwordConfirmation" binding:"required"`
}

type VerificationData struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Email     string    `json:"email" validate:"required"`
	Code      string    `json:"code" validate:"required"`
	ExpiresAt time.Time `json:"expiresat"`
}

func NewVerificationData(userID uuid.UUID, email string, code string) VerificationData {
	uuidGenerator := UUIDGeneratorImpl{}

	return VerificationData{
		ID:        uuidGenerator.Generate(),
		UserID:    userID,
		Email:     email,
		Code:      code,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email" binding:"required"`
}

type SigninResponse struct {
	AccessToken string `json:"accessToken"`
}

type SignupResponse struct {
	AccessToken string `json:"acessToken"`
}
