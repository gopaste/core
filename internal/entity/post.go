package entity

import (
	"time"

	"github.com/Caixetadev/snippet/internal/utils"
)

type PostInput struct {
	ID          string    `json:"id"`
	UserID      *string   `json:"-"`
	Title       string    `json:"title" validate:"required" binding:"required"`
	Content     string    `json:"content,omitempty" validate:"required" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
	Password    string    `json:"password,omitempty"`
	HasPassword bool      `json:"has_password"`
}

type PostOutput struct {
	ID          string    `json:"id"`
	UserID      *string   `json:"-"`
	Title       string    `json:"title" validate:"required" binding:"required"`
	Content     string    `json:"content,omitempty" validate:"required" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
	Password    string    `json:"-"`
	HasPassword bool      `json:"has_password"`
}

type PostUpdateInput struct {
	ID      string `json:"-"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type GetPostInput struct {
	Password string `json:"password,omitempty"`
}

type PaginationInfo struct {
	Next  *string `json:"next"`
	Prev  *string `json:"prev"`
	Pages int     `json:"pages"`
	Count int     `json:"count"`
}

func NewPost(userID *string, title string, content string, password string, hasPassword bool) *PostInput {
	return &PostInput{
		ID:          utils.GenerateRandomString(8),
		UserID:      userID,
		Title:       title,
		Content:     content,
		Password:    password,
		HasPassword: hasPassword,
	}
}
