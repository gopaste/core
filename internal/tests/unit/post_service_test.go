package unit

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/Caixetadev/snippet/internal/entity"
	"github.com/Caixetadev/snippet/internal/services"
	"github.com/Caixetadev/snippet/internal/tests/mocks"
	"github.com/Caixetadev/snippet/pkg/typesystem"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type PostServiceTestSuite struct {
	suite.Suite
	mocksRepo           *mocks.PostRepository
	validation          *mocks.Validator
	postService         *services.PostService
	mocksPasswordHasher *mocks.PasswordHasher
}

func (suite *PostServiceTestSuite) SetupTest() {
	suite.mocksRepo = new(mocks.PostRepository)
	suite.validation = new(mocks.Validator)
	suite.mocksPasswordHasher = new(mocks.PasswordHasher)
	suite.postService = services.NewPostService(suite.mocksRepo, suite.validation, suite.mocksPasswordHasher)
}

func TestPostServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PostServiceTestSuite))
}

func (suite *PostServiceTestSuite) TestCreate() {
	ctx := context.TODO()

	userID := "22c15b0d-5445-4c84-a52a-40888798d1d0"

	input := &entity.PostInput{
		UserID:  &userID,
		Title:   "Title",
		Content: "Body",
	}

	suite.validation.On("Validate", mock.Anything).Return(nil).Once()
	suite.mocksRepo.On("Insert", ctx, mock.AnythingOfType("*entity.PostInput")).Return(nil).Once()

	err := suite.postService.Create(ctx, input)

	suite.NoError(err)

	suite.mocksRepo.AssertExpectations(suite.T())
	suite.validation.AssertExpectations(suite.T())
}

func (suite *PostServiceTestSuite) TestCreate_UserAnonymous() {
	ctx := context.TODO()

	userID := ""

	input := &entity.PostInput{
		UserID:  &userID,
		Title:   "Title",
		Content: "Body",
	}

	suite.validation.On("Validate", mock.Anything).Return(nil).Once()
	suite.mocksRepo.On("Insert", ctx, mock.Anything).Return(nil).Once()

	err := suite.postService.Create(ctx, input)

	suite.NoError(err)

	suite.mocksRepo.AssertExpectations(suite.T())
	suite.validation.AssertExpectations(suite.T())
}

func (suite *PostServiceTestSuite) TestCreate_ValidationError() {
	ctx := context.TODO()

	userID := "22c15b0d-5445-4c84-a52a-40888798d1d0"

	input := &entity.PostInput{
		UserID: &userID,
	}

	suite.validation.On("Validate", mock.Anything).Return(errors.New("error")).Once()

	err := suite.postService.Create(ctx, input)

	suite.Equal(typesystem.BadRequest, err)

	suite.validation.AssertExpectations(suite.T())
}

func (suite *PostServiceTestSuite) TestCreate_PostRepositoryError() {
	ctx := context.TODO()

	userID := "22c15b0d-5445-4c84-a52a-40888798d1d0"

	input := &entity.PostInput{
		UserID:  &userID,
		Title:   "Title",
		Content: "Body",
	}

	suite.validation.On("Validate", mock.Anything).Return(nil).Once()
	suite.mocksRepo.On("Insert", ctx, mock.Anything).Return(errors.New("error")).Once()

	err := suite.postService.Create(ctx, input)

	suite.Equal(typesystem.ServerError, err)

	suite.mocksRepo.AssertExpectations(suite.T())
	suite.validation.AssertExpectations(suite.T())
}

func (suite *PostServiceTestSuite) TestGetPosts() {
	ctx := context.TODO()

	userID := uuid.New()
	userIDStr := userID.String()
	page := "1"

	output := []*entity.PostOutput{
		{
			ID:      uuid.New(),
			UserID:  &userIDStr,
			Title:   "Title",
			Content: "Body",
		},
	}

	suite.mocksRepo.On("CountUserPosts", ctx, userID).Return(10, nil).Once()
	suite.mocksRepo.On("FindAll", ctx, mock.Anything).Return(output, nil).Once()

	posts, _, err := suite.postService.GetPosts(ctx, userID, page)

	suite.NoError(err)
	suite.Equal(output, posts)

	suite.mocksRepo.AssertExpectations(suite.T())
}

func (suite *PostServiceTestSuite) TestGetPosts_PostRepositoryError() {
	ctx := context.TODO()

	userID := uuid.New()
	page := "1"

	suite.mocksRepo.On("FindAll", ctx, mock.Anything).Return([]*entity.PostOutput{}, errors.New("error")).Once()
	suite.mocksRepo.On("CountUserPosts", ctx, userID).Return(10, nil).Once()
	posts, _, err := suite.postService.GetPosts(ctx, userID, page)

	suite.Equal(typesystem.ServerError, err)
	suite.Nil(posts)

	suite.mocksRepo.AssertExpectations(suite.T())
}

func (suite *PostServiceTestSuite) TestDeletePost() {
	ctx := context.TODO()

	userID := uuid.New()
	userIDStr := userID.String()
	postID := uuid.New()

	output := &entity.PostOutput{
		ID:      postID,
		UserID:  &userIDStr,
		Title:   "Title",
		Content: "Body",
	}

	suite.mocksRepo.On("FindOneByID", ctx, postID).Return(output, nil).Once()
	suite.mocksRepo.On("Delete", ctx, postID).Return(nil).Once()

	err := suite.postService.DeletePost(ctx, postID, userID)

	suite.NoError(err)

	suite.mocksRepo.AssertExpectations(suite.T())
}

func (suite *PostServiceTestSuite) TestDeletePost_UserNonAuth() {
	ctx := context.TODO()

	userIDNon := uuid.New()

	userID := uuid.New()
	userIDStr := userID.String()
	postID := uuid.New()

	output := &entity.PostOutput{
		ID:      postID,
		UserID:  &userIDStr,
		Title:   "Title",
		Content: "Body",
	}

	suite.mocksRepo.On("FindOneByID", ctx, postID).Return(output, nil).Once()

	err := suite.postService.DeletePost(ctx, postID, userIDNon)

	suite.Equal(typesystem.Unauthorized, err)

	suite.mocksRepo.AssertExpectations(suite.T())
}

func (suite *PostServiceTestSuite) TestDeletePost_PostNotFound() {
	ctx := context.TODO()

	userID := uuid.New()
	postID := uuid.New()

	suite.mocksRepo.On("FindOneByID", ctx, postID).Return(&entity.PostOutput{}, sql.ErrNoRows).Once()

	err := suite.postService.DeletePost(ctx, postID, userID)

	suite.Equal(typesystem.NotFound, err)

	suite.mocksRepo.AssertExpectations(suite.T())
}

func (suite *PostServiceTestSuite) TestDeletePost_GetPostsRepositoryError() {
	ctx := context.TODO()

	userID := uuid.New()
	postID := uuid.New()

	suite.mocksRepo.On("FindOneByID", ctx, postID).Return(&entity.PostOutput{}, errors.New("error")).Once()

	err := suite.postService.DeletePost(ctx, postID, userID)

	suite.Equal(typesystem.ServerError, err)

	suite.mocksRepo.AssertExpectations(suite.T())
}

func (suite *PostServiceTestSuite) TestDeletePost_DeleteRepositoryError() {
	ctx := context.TODO()

	userID := uuid.New()
	userIDStr := userID.String()
	postID := uuid.New()

	output := &entity.PostOutput{
		ID:      postID,
		UserID:  &userIDStr,
		Title:   "Title",
		Content: "Body",
	}

	suite.mocksRepo.On("FindOneByID", ctx, postID).Return(output, nil).Once()
	suite.mocksRepo.On("Delete", ctx, postID).Return(errors.New("error")).Once()

	err := suite.postService.DeletePost(ctx, postID, userID)

	suite.Equal(typesystem.ServerError, err)

	suite.mocksRepo.AssertExpectations(suite.T())
}
