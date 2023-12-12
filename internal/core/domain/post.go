package domain

import (
	"context"

	"github.com/google/uuid"
)

type Post struct {
	ID      uuid.UUID `json:"-" `
	UserID  *string   `json:"-"`
	Title   string    `json:"title" validate:"required" binding:"required"`
	Content string    `json:"content" validate:"required" binding:"required"`
}

type PostRepository interface {
	Create(ctx context.Context, post *Post) error
}

type PostService interface {
	Create(ctx context.Context, post *Post) error
}
