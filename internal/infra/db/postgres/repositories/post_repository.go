package repository

import (
	"context"
	"database/sql"

	"github.com/Caixetadev/snippet/internal/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postRepository struct {
	db *pgxpool.Pool
}

func NewPostRepository(db *pgxpool.Pool) *postRepository {
	return &postRepository{db: db}
}

func (pr *postRepository) Create(ctx context.Context, post *entity.Post) error {
	_, err := pr.db.Exec(ctx, "INSERT INTO posts (id, user_id, title, content) VALUES ($1, $2, $3, $4)", post.ID, post.UserID, post.Title, post.Content)

	return err
}

func (pr *postRepository) GetPosts(ctx context.Context, id uuid.UUID) ([]*entity.Post, error) {
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

func (pr *postRepository) GetPostByID(ctx context.Context, id uuid.UUID) (*entity.Post, error) {
	line, err := pr.db.Query(ctx, "SELECT id, user_id, title, content FROM posts WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	defer line.Close()

	var post entity.Post

	if line.Next() {
		if err := line.Scan(&post.ID, &post.UserID, &post.Title, &post.Content); err != nil {
			return nil, err
		}
	} else {
		return nil, sql.ErrNoRows
	}

	return &post, nil
}

func (pr *postRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := pr.db.Exec(ctx, "DELETE FROM posts WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
