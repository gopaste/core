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

func (ps *PostRepository) Insert(ctx context.Context, post *entity.PostInput) error {
	args := ps.Called(ctx, post)
	return args.Error(0)
}

func (ps *PostRepository) FindAll(ctx context.Context, id uuid.UUID, limit, offset int) ([]*entity.PostOutput, error) {
	args := ps.Called(ctx, id)
	return args.Get(0).([]*entity.PostOutput), args.Error(1)
}

func (ps *PostRepository) FindOneByID(ctx context.Context, id string) (*entity.PostOutput, error) {
	args := ps.Called(ctx, id)
	return args.Get(0).(*entity.PostOutput), args.Error(1)
}

func (ps *PostRepository) Delete(ctx context.Context, id string) error {
	args := ps.Called(ctx, id)
	return args.Error(0)
}

func (ps *PostRepository) Update(ctx context.Context, post *entity.PostUpdateInput) error {
	args := ps.Called(ctx, post)
	return args.Error(0)
}

func (ps *PostRepository) CountPostsInSearch(ctx context.Context, query string) (int, error) {
	args := ps.Called(ctx, query)
	return args.Int(0), args.Error(1)
}

func (ps *PostRepository) CountUserPosts(ctx context.Context, id uuid.UUID) (int, error) {
	args := ps.Called(ctx, id)
	return args.Int(0), args.Error(1)
}

func (ps *PostRepository) Search(ctx context.Context, query string, limit, offset int) ([]*entity.PostOutput, error) {
	args := ps.Called(ctx, query)
	return args.Get(0).([]*entity.PostOutput), args.Error(1)
}
