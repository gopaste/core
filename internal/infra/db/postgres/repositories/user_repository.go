package repository

import (
	"github.com/Caixetadev/snippet/internal/core/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *userRepository {
	return &userRepository{db: db}
}

func (ur *userRepository) GetUserByEmail(email string) (*domain.User, error) {
	return nil, nil
}
