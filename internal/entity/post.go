package entity

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID `json:"id"`
	UserID    *string   `json:"-"`
	Title     string    `json:"title" validate:"required" binding:"required"`
	Content   string    `json:"content,omitempty" validate:"required" binding:"required"`
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
