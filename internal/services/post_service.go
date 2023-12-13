package services

import (
	"github.com/google/uuid"

	"github.com/Caixetadev/snippet/internal/entity"
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
		return entity.BadRequest
	}

	input.ID = uuid.New()

	if len(*input.UserID) == 0 {
		input.UserID = nil
	}

	err = ps.postRepo.Create(ctx, input)
	if err != nil {
		return entity.ServerError
	}

	return nil
}
