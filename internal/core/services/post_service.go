package services

import (
	"github.com/Caixetadev/snippet/internal/core/domain"
	"github.com/google/uuid"

	apperr "github.com/Caixetadev/snippet/internal/core/error"
	"github.com/Caixetadev/snippet/pkg/validation"
	"golang.org/x/net/context"
)

type PostService struct {
	postRepo   domain.PostRepository
	validation validation.Validator
}

func NewPostService(postRepo domain.PostRepository, validation validation.Validator) *PostService {
	return &PostService{postRepo: postRepo, validation: validation}
}

func (ps *PostService) Create(ctx context.Context, input *domain.Post) error {
	err := ps.validation.Validate(input)

	if err != nil {
		return apperr.BadRequest
	}

	input.ID = uuid.New()

	if len(*input.UserID) == 0 {
		input.UserID = nil
	}

	err = ps.postRepo.Create(ctx, input)
	if err != nil {
		return apperr.ServerError
	}

	return nil
}
