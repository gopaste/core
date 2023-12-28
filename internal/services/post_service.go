package services

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/Caixetadev/snippet/internal/entity"
	"github.com/Caixetadev/snippet/internal/utils"
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

func (ps *PostService) GetPosts(ctx context.Context, id uuid.UUID, pageStr string) ([]*entity.Post, *entity.PaginationInfo, error) {
	if pageStr == "" {
		pageStr = "1"
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return nil, nil, typesystem.BadRequest
	}

	count, err := ps.postRepo.Count(ctx, id)
	if err != nil {
		return nil, nil, typesystem.ServerError
	}

	limit := 10

	totalPages := (count + limit - 1) / limit

	if totalPages < page {
		return nil, nil, typesystem.NotFound
	}

	offset := (page - 1) * limit

	posts, err := ps.postRepo.GetPosts(ctx, id, limit, offset)
	if err != nil {
		return nil, nil, typesystem.ServerError
	}

	nextPage, prevPage := "", ""
	if len(posts) == limit {
		nextPage = fmt.Sprintf("/post/all?page=%d", page+1)
	}
	if page > 1 {
		prevPage = fmt.Sprintf("/post/all?page=%d", page-1)
	}

	paginationInfo := &entity.PaginationInfo{
		Next:  utils.StringToPtr(nextPage),
		Prev:  utils.StringToPtr(prevPage),
		Pages: totalPages,
		Count: count,
	}

	return posts, paginationInfo, nil
}

func (ps *PostService) DeletePost(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	post, err := ps.postRepo.GetPostByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return typesystem.NotFound
		}
		return typesystem.ServerError
	}

	if *post.UserID != userID.String() {
		return typesystem.Unauthorized
	}

	err = ps.postRepo.Delete(ctx, id)
	if err != nil {
		return typesystem.ServerError
	}

	return nil
}

func (ps *PostService) UpdatePost(ctx context.Context, post *entity.PostUpdateInput, userID uuid.UUID, id uuid.UUID) error {
	postInDatabase, err := ps.postRepo.GetPostByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return typesystem.NotFound
		}
		return typesystem.ServerError
	}

	if *postInDatabase.UserID != userID.String() {
		return typesystem.Unauthorized
	}

	post.ID = postInDatabase.ID

	err = ps.postRepo.Update(ctx, post)
	if err != nil {
		return typesystem.ServerError
	}

	return nil
}
