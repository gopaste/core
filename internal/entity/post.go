package entity

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID `json:"id"`
	UserID    *string   `json:"-"`
	Title     string    `json:"title" validate:"required" binding:"required"`
	Content   string    `json:"content" validate:"required" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
}

type PostUpdateInput struct {
	ID      uuid.UUID `json:"-"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
}

type PaginationInfo struct {
	Next  *string `json:"next"`
	Prev  *string `json:"prev"`
	Pages int     `json:"pages"`
	Count int     `json:"count"`
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
	GetPosts(ctx context.Context, id uuid.UUID, limit int, offset int) ([]*Post, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Count(ctx context.Context, id uuid.UUID) (int, error)
	GetPostByID(ctx context.Context, id uuid.UUID) (*Post, error)
	Update(ctx context.Context, post *PostUpdateInput) error
}

type PostService interface {
	Create(ctx context.Context, post *Post) error
	GetPosts(ctx context.Context, id uuid.UUID, page string) ([]*Post, *PaginationInfo, error)
	DeletePost(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	UpdatePost(ctx context.Context, post *PostUpdateInput, userID uuid.UUID, id uuid.UUID) error
}
