package repository

import (
	"context"

	"github.com/Caixetadev/snippet/internal/core/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postRepository struct {
	db *pgxpool.Pool
}

func NewPostRepository(db *pgxpool.Pool) *postRepository {
	return &postRepository{db: db}
}

func (pr *postRepository) Create(ctx context.Context, post *domain.Post) error {
	_, err := pr.db.Exec(ctx, "INSERT INTO posts (post_id, user_id, title, content) VALUES ($1, $2, $3, $4)", post.ID, post.UserID, post.Title, post.Content)

	return err
}
