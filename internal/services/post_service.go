package services

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Caixetadev/snippet/internal/entity"
	"github.com/Caixetadev/snippet/internal/utils"
	"github.com/Caixetadev/snippet/pkg/typesystem"
	"github.com/Caixetadev/snippet/pkg/validation"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type PostRepository interface {
	Insert(ctx context.Context, post *entity.PostInput) error
	FindAll(ctx context.Context, id uuid.UUID, limit int, offset int) ([]*entity.PostOutput, error)
	Delete(ctx context.Context, id uuid.UUID) error
	CountUserPosts(ctx context.Context, id uuid.UUID) (int, error)
	CountPostsInSearch(ctx context.Context, query string) (int, error)
	FindOneByID(ctx context.Context, id uuid.UUID) (*entity.PostOutput, error)
	Update(ctx context.Context, post *entity.PostUpdateInput) error
	Search(ctx context.Context, query string, limit int, offset int) ([]*entity.PostOutput, error)
}

type PostService struct {
	postRepo   PostRepository
	validation validation.Validator
}

func NewPostService(postRepo PostRepository, validation validation.Validator) *PostService {
	return &PostService{postRepo: postRepo, validation: validation}
}

func (ps *PostService) Create(ctx context.Context, input *entity.PostInput) error {
	err := ps.validation.Validate(input)
	if err != nil {
		return typesystem.BadRequest
	}

	post := entity.NewPost(input.UserID, input.Title, input.Content, input.Password, input.IsPrivate)

	if len(*post.UserID) == 0 {
		post.UserID = nil
	}

	err = ps.postRepo.Insert(ctx, post)
	if err != nil {
		return typesystem.ServerError
	}

	return nil
}

func (ps *PostService) GetPosts(ctx context.Context, id uuid.UUID, pageStr string) ([]*entity.PostOutput, *entity.PaginationInfo, error) {
	count, err := ps.postRepo.CountUserPosts(ctx, id)
	if err != nil {
		return nil, nil, typesystem.ServerError
	}

	pageResponse, err := utils.CalculePagination(count, pageStr)
	if err != nil {
		return nil, nil, err
	}

	posts, err := ps.postRepo.FindAll(ctx, id, pageResponse.Limit, pageResponse.Offset)
	if err != nil {
		return nil, nil, typesystem.ServerError
	}

	nextPage, prevPage := "", ""
	if len(posts) == pageResponse.Limit {
		nextPage = fmt.Sprintf("/post/all?page=%d", pageResponse.Page+1)
	}
	if pageResponse.Page > 1 {
		prevPage = fmt.Sprintf("/post/all?page=%d", pageResponse.Page-1)
	}

	paginationInfo := &entity.PaginationInfo{
		Next:  utils.StringToPtr(nextPage),
		Prev:  utils.StringToPtr(prevPage),
		Pages: pageResponse.Total,
		Count: count,
	}

	return posts, paginationInfo, nil
}

func (ps *PostService) DeletePost(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	post, err := ps.postRepo.FindOneByID(ctx, id)
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
	postInDatabase, err := ps.postRepo.FindOneByID(ctx, id)
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

func (ps *PostService) SearchPost(ctx context.Context, query string, pageStr string) ([]*entity.PostOutput, *entity.PaginationInfo, error) {
	count, err := ps.postRepo.CountPostsInSearch(ctx, query)
	if err != nil {
		return nil, nil, typesystem.ServerError
	}

	pageResponse, err := utils.CalculePagination(count, pageStr)
	if err != nil {
		return nil, nil, err
	}

	posts, err := ps.postRepo.Search(ctx, query, pageResponse.Limit, pageResponse.Offset)
	if err != nil {
		return nil, nil, typesystem.ServerError
	}

	nextPage, prevPage := "", ""
	if len(posts) == pageResponse.Limit {
		nextPage = fmt.Sprintf("/post/search?q=%s&page=%d", query, pageResponse.Page+1)
	}
	if pageResponse.Page > 1 {
		prevPage = fmt.Sprintf("/post/search?q=%s&page=%d", query, pageResponse.Page-1)
	}

	paginationInfo := &entity.PaginationInfo{
		Next:  utils.StringToPtr(nextPage),
		Prev:  utils.StringToPtr(prevPage),
		Pages: pageResponse.Total,
		Count: count,
	}

	return posts, paginationInfo, nil
}

func (ps *PostService) GetPost(ctx context.Context, id uuid.UUID, password string) (*entity.PostOutput, error) {
	post, err := ps.postRepo.FindOneByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, typesystem.NotFound
		}
		return nil, typesystem.ServerError
	}

	if post.IsPrivate {
		fmt.Println(post.Password, password)
		if post.Password != password {
			return nil, typesystem.Unauthorized
		}
	}

	return post, nil
}
