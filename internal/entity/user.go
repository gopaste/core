package entity

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"-"`
	Name     string    `json:"name"         validate:"required"       binding:"required"`
	Email    string    `json:"email"        validate:"required,email" binding:"required"`
	Password string    `json:"password,omitempty"     validate:"required"       binding:"required"`
}

func NewUser(name, email, password string) *User {
	uuidGenerator := UUIDGeneratorImpl{}

	return &User{
		ID:       uuidGenerator.Generate(),
		Name:     name,
		Email:    email,
		Password: password,
	}
}

type UserService interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	Create(ctx context.Context, user *User) (*User, error)
	UserExistsByEmail(ctx context.Context, email string) (bool, error)
	CreateAccessToken(user *User, expiry time.Duration) (string, *Payload, error)
	CreateRefreshToken(ctx context.Context, user *User, expiry time.Duration) (string, *Payload, error)
	CompareHashAndPassword(passwordInDatabase, passwordRequest string) error
	StoreVerificationData(ctx context.Context, userID uuid.UUID, email string, code string) error
	VerifyCodeToResetPassword(ctx context.Context, code string) (uuid.UUID, error)
	UpdatePassword(ctx context.Context, password string, passwordConfirmation string, id uuid.UUID) error
	VerifyToken(ctx context.Context, token string) (*Payload, error)
	GetSession(ctx context.Context, id uuid.UUID) (*Session, error)
	CreateSession(ctx context.Context, payload *Payload, token string) error
	RevokeRefreshToken(ctx context.Context, token string) error
}
