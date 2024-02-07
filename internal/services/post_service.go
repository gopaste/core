package services

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Caixetadev/snippet/internal/entity"
	"github.com/Caixetadev/snippet/internal/pagination"
	"github.com/Caixetadev/snippet/pkg/passwordhash"
	"github.com/Caixetadev/snippet/pkg/typesystem"
	"github.com/Caixetadev/snippet/pkg/validation"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

var (
	ErrAccountRequired = typesystem.NewHttpError(
		"Cannot create a private post without an account. Please log in or create an account.",
		"[Error: Account required for private post]",
		http.StatusUnauthorized,
	)
	ErrDeleteAndViewConflict = typesystem.NewHttpError(
		"Cannot have both delete_after_view and expiration_at together",
		"[Error: delete_after_view_conflict]",
		http.StatusBadRequest,
	)
	ErrPasswordLength = typesystem.NewHttpError(
		"Password should have a minimum of 3 characters.",
		"[Error: password_length]",
		http.StatusBadRequest,
	)
)

type PostRepository interface {
	Insert(ctx context.Context, post *entity.PostInput) error
	FindAll(ctx context.Context, id uuid.UUID, page int) ([]*entity.PostOutput, int, error)
	FindAllPublics(ctx context.Context, page int) ([]*entity.PostOutput, int, error)
	Delete(ctx context.Context, id string) error
	CountUserPosts(ctx context.Context, id uuid.UUID) (int, error)
	CountAllPostsPublics(ctx context.Context) (int, error)
	CountPostsInSearch(ctx context.Context, query string) (int, error)
	FindOneByID(ctx context.Context, id string) (*entity.PostOutput, error)
	Update(ctx context.Context, post *entity.PostUpdateInput) error
	Search(ctx context.Context, query string, page int) ([]*entity.PostOutput, int, error)
}

type PostService struct {
	postRepo       PostRepository
	validation     validation.Validator
	passwordHasher passwordhash.PasswordHasher
}

func NewPostService(
	postRepo PostRepository,
	validation validation.Validator,
	passwordHasher passwordhash.PasswordHasher,
) *PostService {
	return &PostService{postRepo: postRepo, validation: validation, passwordHasher: passwordHasher}
}

func (ps *PostService) Create(ctx context.Context, input *entity.PostInput) error {
	err := ps.validation.Validate(input)
	if err != nil {
		return typesystem.BadRequest
	}

	if *input.UserID == "" && input.Visibility == entity.Private {
		return ErrAccountRequired
	}

	if input.DeleteAfterView && !input.ExpirationAt.IsZero() {
		return ErrDeleteAndViewConflict
	}

	if input.HasPassword {
		if len(input.Password) < 3 {
			return ErrPasswordLength
		}

		encryptedPassword, err := ps.passwordHasher.GenerateFromPassword([]byte(input.Password), 10)
		if err != nil {
			return typesystem.ServerError
		}

		input.Password = string(encryptedPassword)
	} else {
		input.Password = ""
	}

	post := entity.NewPost(
		input.UserID,
		input.Title,
		input.Content,
		input.Password,
		input.HasPassword,
		input.Visibility,
		input.ExpirationAt,
		input.DeleteAfterView,
	)

	if len(*post.UserID) == 0 {
		post.UserID = nil
	}

	err = ps.postRepo.Insert(ctx, post)
	if err != nil {
		fmt.Println(err.Error())
		return typesystem.ServerError
	}

	return nil
}

func (ps *PostService) GetPosts(
	ctx context.Context,
	id uuid.UUID,
	pageStr string,
) ([]*entity.PostOutput, *entity.PaginationInfo, error) {
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return nil, nil, typesystem.ServerError
	}

	posts, count, err := ps.postRepo.FindAll(ctx, id, page)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, typesystem.NotFound
		}
		return nil, nil, typesystem.ServerError
	}

	paginationInfo, err := pagination.GeneratePaginationInfo(count, page, "/post/user/all")
	if err != nil {
		return nil, nil, err
	}

	return posts, paginationInfo, nil
}

func (ps *PostService) GetAllPublics(
	ctx context.Context,
	pageStr string,
) ([]*entity.PostOutput, *entity.PaginationInfo, error) {
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return nil, nil, typesystem.ServerError
	}

	posts, count, err := ps.postRepo.FindAllPublics(ctx, page)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, typesystem.NotFound
		}
		return nil, nil, typesystem.ServerError
	}

	paginationInfo, err := pagination.GeneratePaginationInfo(count, page, "/post/all")
	if err != nil {
		return nil, nil, err
	}

	return posts, paginationInfo, nil
}

func (ps *PostService) DeletePost(ctx context.Context, id string, userID uuid.UUID) error {
	post, err := ps.postRepo.FindOneByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return typesystem.NotFound
		}
		return typesystem.ServerError
	}

	if *post.UserID != userID.String() {
		return typesystem.Forbidden
	}

	err = ps.postRepo.Delete(ctx, id)
	if err != nil {
		return typesystem.ServerError
	}

	return nil
}

func (ps *PostService) UpdatePost(
	ctx context.Context,
	post *entity.PostUpdateInput,
	userID uuid.UUID,
	id string,
) error {
	postInDatabase, err := ps.postRepo.FindOneByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return typesystem.NotFound
		}
		return typesystem.ServerError
	}

	if *postInDatabase.UserID != userID.String() {
		return typesystem.Forbidden
	}

	post.ID = postInDatabase.ID

	err = ps.postRepo.Update(ctx, post)
	if err != nil {
		return typesystem.ServerError
	}

	return nil
}

func (ps *PostService) SearchPost(
	ctx context.Context,
	query string,
	pageStr string,
) ([]*entity.PostOutput, *entity.PaginationInfo, error) {
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return nil, nil, typesystem.ServerError
	}

	posts, count, err := ps.postRepo.Search(ctx, query, page)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, typesystem.NotFound
		}
		return nil, nil, typesystem.ServerError
	}

	paginationInfo, err := pagination.GeneratePaginationInfo(count, page, fmt.Sprintf("/post/search?q=%s&", query))
	if err != nil {
		return nil, nil, err
	}

	return posts, paginationInfo, nil
}

func (ps *PostService) GetPost(
	ctx context.Context,
	id string,
	userID string,
	password string,
) (*entity.PostOutput, error) {
	post, err := ps.postRepo.FindOneByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, typesystem.NotFound
		}
		return nil, typesystem.ServerError
	}

	if !post.ExpirationAt.IsZero() && time.Now().After(post.ExpirationAt) {
		go ps.postRepo.Delete(ctx, post.ID)

		return nil, typesystem.NotFound
	}

	if post.HasPassword {
		err := ps.passwordHasher.CompareHashAndPassword([]byte(post.Password), []byte(password))
		if err != nil {
			return nil, typesystem.Unauthorized
		}
	}

	if post.DeleteAfterView {
		defer ps.postRepo.Delete(ctx, post.ID)
	}

	if post.Visibility == entity.Private {
		if *post.UserID != userID {
			return nil, typesystem.NotFound
		}
	}

	return post, nil
}
