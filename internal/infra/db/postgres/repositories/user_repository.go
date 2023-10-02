package repository

import (
	"context"

	"github.com/Caixetadev/snippet/internal/core/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *userRepository {
	return &userRepository{db: db}
}

func (ur *userRepository) Create(ctx context.Context, user *domain.User) error {
	_, err := ur.db.Exec(
		ctx,
		"INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4)",
		user.ID,
		user.Name,
		user.Email,
		user.Password,
	)
	return err
}

func (ur *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	line, err := ur.db.Query(ctx, "SELECT id, name, email, password FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}

	defer line.Close()

	var user domain.User

	if line.Next() {
		if err = line.Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
			return nil, err
		}
	} else {
		return nil, pgx.ErrNoRows
	}

	return &user, nil
}

func (ur *userRepository) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := ur.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
