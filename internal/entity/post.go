package entity

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

func NewPost(userID *string, title string, content string) *Post {
	uuidGenerator := UUIDGeneratorImpl{}

	return &Post{
		ID:      uuidGenerator.Generate(),
		UserID:  userID,
		Title:   title,
		Content: content,
	}
}

type PostRepository interface {
	Create(ctx context.Context, post *Post) error
	GetPosts(ctx context.Context, id uuid.UUID) ([]*Post, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetPostByID(ctx context.Context, id uuid.UUID) (*Post, error)
}

type PostService interface {
	Create(ctx context.Context, post *Post) error
	GetPosts(ctx context.Context, id uuid.UUID) ([]*Post, error)
	DeletePost(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}
