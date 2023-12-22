package services

import (
	"github.com/Caixetadev/snippet/internal/entity"
	"github.com/Caixetadev/snippet/pkg/typesystem"
	"github.com/Caixetadev/snippet/pkg/validation"
	"golang.org/x/net/context"
)

type PostService struct {
	postRepo   entity.PostRepository
	validation validation.Validator
}

func NewPostService(postRepo entity.PostRepository, validation validation.Validator) *PostService {
	return &PostService{postRepo: postRepo, validation: validation}
}

func (ps *PostService) Create(ctx context.Context, input *entity.Post) error {
	err := ps.validation.Validate(input)

	if err != nil {
		return typesystem.BadRequest
	}

	post := entity.NewPost(input.UserID, input.Title, input.Content)

	if len(*post.UserID) == 0 {
		post.UserID = nil
	}

	err = ps.postRepo.Create(ctx, post)
	if err != nil {
		return typesystem.ServerError
	}

	return nil
}

func (ps *PostService) GetPosts(ctx context.Context, id string) ([]*entity.Post, error) {
	if len(id) == 0 {
		return nil, typesystem.Unauthorized
	}

	posts, err := ps.postRepo.GetPosts(ctx, id)
	if err != nil {
		return nil, typesystem.ServerError
	}

	return posts, nil
}
