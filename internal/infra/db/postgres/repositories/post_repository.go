package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Caixetadev/snippet/internal/entity"
	"github.com/Caixetadev/snippet/internal/services"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ services.PostRepository = (*postRepository)(nil)

type postRepository struct {
	db *pgxpool.Pool
}

func NewPostRepository(db *pgxpool.Pool) *postRepository {
	return &postRepository{db: db}
}

func (pr *postRepository) Insert(ctx context.Context, post *entity.PostInput) error {
	query := "INSERT INTO posts (id, user_id, title, content, password, is_private) VALUES ($1, $2, $3, $4, $5, $6)"

	_, err := pr.db.Exec(ctx, query, post.ID, post.UserID, post.Title, post.Content, post.Password, post.IsPrivate)

	return err
}

func (pr *postRepository) FindAll(ctx context.Context, id uuid.UUID, limit int, offset int) ([]*entity.PostOutput, error) {
	query := `
		SELECT id, title, created_at, is_private
		FROM posts
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3;
	`

	line, err := pr.db.Query(ctx, query, id, limit, offset)
	if err != nil {
		return nil, err
	}

	defer line.Close()

	var posts []*entity.PostOutput

	for line.Next() {
		post := &entity.PostOutput{}
		if err := line.Scan(&post.ID, &post.Title, &post.CreatedAt, &post.IsPrivate); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (pr *postRepository) CountUserPosts(ctx context.Context, id uuid.UUID) (int, error) {
	var count int

	query := "SELECT COUNT(*) FROM posts WHERE user_id = $1"

	err := pr.db.QueryRow(ctx, query, id).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (pr *postRepository) CountPostsInSearch(ctx context.Context, query string) (int, error) {
	var count int

	querySql := "SELECT COUNT(*) FROM posts WHERE title ILIKE '%' || $1 || '%' OR content ILIKE '%' || $1 || '%'"

	err := pr.db.QueryRow(ctx, querySql, query).
		Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (pr *postRepository) FindOneByID(ctx context.Context, id uuid.UUID) (*entity.PostOutput, error) {
	query := "SELECT id, user_id, title, content, created_at, password, is_private FROM posts WHERE id = $1"

	line, err := pr.db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}

	defer line.Close()

	var post entity.PostOutput

	if line.Next() {
		if err := line.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.Password, &post.IsPrivate); err != nil {
			return nil, err
		}
	} else {
		return nil, sql.ErrNoRows
	}

	return &post, nil
}

func (pr *postRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM posts WHERE id = $1"

	_, err := pr.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

// TODO: refactor this function
func (pr *postRepository) Update(ctx context.Context, post *entity.PostUpdateInput) error {
	query := "UPDATE posts SET"
	var args []interface{}

	if post.Title != "" {
		query += fmt.Sprintf(" title = $%d,", len(args)+1)
		args = append(args, post.Title)
	}

	if post.Content != "" {
		query += fmt.Sprintf(" content = $%d,", len(args)+1)
		args = append(args, post.Content)
	}

	query = strings.TrimSuffix(query, ",")

	query += fmt.Sprintf(" WHERE id = $%d", len(args)+1)
	args = append(args, post.ID)

	_, err := pr.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (pr *postRepository) Search(ctx context.Context, q string, limit int, offset int) ([]*entity.PostOutput, error) {
	query := `
		SELECT id, title, content, is_private, created_at
		FROM posts
		WHERE title ILIKE '%' || $1 || '%' OR content ILIKE '%' || $1 || '%'
		ORDER BY created_at DESC, id DESC
		LIMIT $2 OFFSET $3;
	`

	line, err := pr.db.Query(ctx, query, q, limit, offset)
	if err != nil {
		return nil, err
	}

	defer line.Close()

	var posts []*entity.PostOutput

	for line.Next() {
		post := &entity.PostOutput{}
		if err := line.Scan(&post.ID, &post.Title, &post.Content, &post.IsPrivate, &post.CreatedAt); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}
