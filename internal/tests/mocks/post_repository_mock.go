package mocks

import (
	"context"

	"github.com/Caixetadev/snippet/internal/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type PostRepository struct {
	mock.Mock
}

func (ps *PostRepository) Create(ctx context.Context, post *entity.Post) error {
	args := ps.Called(ctx, post)
	return args.Error(0)
}

func (ps *PostRepository) GetPosts(ctx context.Context, id uuid.UUID) ([]*entity.Post, error) {
	args := ps.Called(ctx, id)
	return args.Get(0).([]*entity.Post), args.Error(1)
}

func (ps *PostRepository) GetPostByID(ctx context.Context, id uuid.UUID) (*entity.Post, error) {
	args := ps.Called(ctx, id)
	return args.Get(0).(*entity.Post), args.Error(1)
}

func (ps *PostRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := ps.Called(ctx, id)
	return args.Error(0)
}
