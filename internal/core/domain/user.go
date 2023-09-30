package domain

import (
	"context"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id,omitempty"`
	Name     string    `json:"name" validate:"required" binding:"required"`
	Email    string    `json:"email" validate:"required,email" binding:"required"`
	Password string    `json:"password" validate:"required" binding:"required"`
}

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) error // Create(user *User) error

	UserExistsByEmail(ctx context.Context, email string) (bool, error)
}

type SignupService interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) error // Create(user *User) error
	UserExistsByEmail(ctx context.Context, email string) (bool, error)
}
