package entity

import (
	"context"
	"time"

	"github.com/Caixetadev/snippet/internal/utils"
	"github.com/gin-gonic/gin"
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

type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
}

func NewRefreshToken(user *User) *RefreshToken {
	uuidGenerator := UUIDGeneratorImpl{}

	return &RefreshToken{
		ID:        uuidGenerator.Generate(),
		UserID:    user.ID,
		Token:     utils.GenerateRandomString(8),
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}
}

type Session struct {
	ID           uuid.UUID
	Name         string
	RefreshToken string
	UserAgent    string
	ClientIp     string
	IsBlocked    bool
	ExpiresAt    time.Time
}

func NewSession(ctx context.Context, payload *Payload, refreshToken string) *Session {
	return &Session{
		ID:           payload.ID,
		Name:         payload.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.(*gin.Context).Request.UserAgent(),
		ClientIp:     ctx.(*gin.Context).ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    payload.ExpiredAt,
	}
}

func NewVerificationData(userID uuid.UUID, email string, code string) *VerificationData {
	uuidGenerator := UUIDGeneratorImpl{}

	return &VerificationData{
		ID:        uuidGenerator.Generate(),
		UserID:    userID,
		Email:     email,
		Code:      code,
		ExpiresAt: time.Now().Add(time.Hour * 15),
	}
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email" binding:"required"`
}

type SigninResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type SignupResponse struct {
	AccessToken string `json:"acessToken"`
}
