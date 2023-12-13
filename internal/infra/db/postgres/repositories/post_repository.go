package repository

import (
	"context"

	"github.com/Caixetadev/snippet/internal/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postRepository struct {
	db *pgxpool.Pool
}

func NewPostRepository(db *pgxpool.Pool) *postRepository {
	return &postRepository{db: db}
}

func (pr *postRepository) Create(ctx context.Context, post *entity.Post) error {
	_, err := pr.db.Exec(ctx, "INSERT INTO posts (post_id, user_id, title, content) VALUES ($1, $2, $3, $4)", post.ID, post.UserID, post.Title, post.Content)

	return err
}

func (pr *postRepository) GetPosts(ctx context.Context, id string) ([]*entity.Post, error) {
	line, err := pr.db.Query(ctx, "SELECT title, content FROM posts WHERE user_id = $1", id)
	if err != nil {
		return nil, err
	}

	defer line.Close()

	var posts []*entity.Post

	for line.Next() {
		post := &entity.Post{}
		if err := line.Scan(&post.Title, &post.Content); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}
