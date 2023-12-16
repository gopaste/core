package entity

import (
	"context"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"-"`
	Name     string    `json:"name"         validate:"required"       binding:"required"`
	Email    string    `json:"email"        validate:"required,email" binding:"required"`
	Password string    `json:"password"     validate:"required"       binding:"required"`
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

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) error
	UserExistsByEmail(ctx context.Context, email string) (bool, error)
	StoreVerificationData(ctx context.Context, verificationData *VerificationData) error
	UpdatePassword(ctx context.Context, password string, id string) error
	VerifyCodeToResetPassword(ctx context.Context, code string) (string, bool, error)
}

type UserService interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) (*User, error)
	UserExistsByEmail(ctx context.Context, email string) (bool, error)
	CreateAccessToken(user *User, secret string, expiry int) (string, error)
	CompareHashAndPassword(passwordInDatabase, passwordRequest string) error
	StoreVerificationData(ctx context.Context, userID uuid.UUID, email string, code string) error
	VerifyCodeToResetPassword(ctx context.Context, code string) (string, bool, error)
	UpdatePassword(ctx context.Context, password string, passwordConfirmation string, id string) error
}
