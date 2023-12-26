package services

import (
	"fmt"

	"github.com/Caixetadev/snippet/internal/entity"
	"github.com/Caixetadev/snippet/pkg/typesystem"
	"github.com/Caixetadev/snippet/pkg/validation"
	"github.com/google/uuid"
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

func (ps *PostService) GetPosts(ctx context.Context, id uuid.UUID) ([]*entity.Post, error) {
	posts, err := ps.postRepo.GetPosts(ctx, id)
	if err != nil {
		return nil, typesystem.ServerError
	}

	return posts, nil
}

func (ps *PostService) DeletePost(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	post, err := ps.postRepo.GetPostByID(ctx, id)
	if err != nil {
		fmt.Println(err)
		return typesystem.ServerError
	}

	if *post.UserID != userID.String() {
		return typesystem.Unauthorized
	}

	err = ps.postRepo.Delete(ctx, id)
	if err != nil {
		fmt.Println(err)
		return typesystem.ServerError
	}

	return nil
}
