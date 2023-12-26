package unit

import (
	"context"
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
	mocksRepo   *mocks.PostRepository
	validation  *mocks.Validator
	postService *services.PostService
}

func (suite *PostServiceTestSuite) SetupTest() {
	suite.mocksRepo = new(mocks.PostRepository)
	suite.validation = new(mocks.Validator)

	suite.postService = services.NewPostService(suite.mocksRepo, suite.validation)
}

func TestPostServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PostServiceTestSuite))
}

func (suite *PostServiceTestSuite) TestCreate() {
	ctx := context.TODO()

	userID := "22c15b0d-5445-4c84-a52a-40888798d1d0"

	input := &entity.Post{
		UserID:  &userID,
		Title:   "Title",
		Content: "Body",
	}

	suite.validation.On("Validate", mock.Anything).Return(nil).Once()
	suite.mocksRepo.On("Create", ctx, mock.Anything).Return(nil).Once()

	err := suite.postService.Create(ctx, input)

	suite.NoError(err)

	suite.mocksRepo.AssertExpectations(suite.T())
	suite.validation.AssertExpectations(suite.T())
}

func (suite *PostServiceTestSuite) TestCreate_UserAnonymous() {
	ctx := context.TODO()

	userID := ""

	input := &entity.Post{
		UserID:  &userID,
		Title:   "Title",
		Content: "Body",
	}

	suite.validation.On("Validate", mock.Anything).Return(nil).Once()
	suite.mocksRepo.On("Create", ctx, mock.Anything).Return(nil).Once()

	err := suite.postService.Create(ctx, input)

	suite.NoError(err)

	suite.mocksRepo.AssertExpectations(suite.T())
	suite.validation.AssertExpectations(suite.T())
}

func (suite *PostServiceTestSuite) TestCreate_ValidationError() {
	ctx := context.TODO()

	userID := "22c15b0d-5445-4c84-a52a-40888798d1d0"

	input := &entity.Post{
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

	input := &entity.Post{
		UserID:  &userID,
		Title:   "Title",
		Content: "Body",
	}

	suite.validation.On("Validate", mock.Anything).Return(nil).Once()
	suite.mocksRepo.On("Create", ctx, mock.Anything).Return(errors.New("error")).Once()

	err := suite.postService.Create(ctx, input)

	suite.Equal(typesystem.ServerError, err)

	suite.mocksRepo.AssertExpectations(suite.T())
	suite.validation.AssertExpectations(suite.T())
}

func (suite *PostServiceTestSuite) TestGetPosts() {
	ctx := context.TODO()

	userID := uuid.New()
	userIDStr := userID.String()

	output := []*entity.Post{
		{
			ID:      uuid.New(),
			UserID:  &userIDStr,
			Title:   "Title",
			Content: "Body",
		},
	}

	suite.mocksRepo.On("GetPosts", ctx, mock.Anything).Return(output, nil).Once()

	posts, err := suite.postService.GetPosts(ctx, userID)

	suite.NoError(err)
	suite.Equal(output, posts)

	suite.mocksRepo.AssertExpectations(suite.T())
}

func (suite *PostServiceTestSuite) TestGetPosts_PostRepositoryError() {
	ctx := context.TODO()

	userID := uuid.New()

	suite.mocksRepo.On("GetPosts", ctx, mock.Anything).Return([]*entity.Post{}, errors.New("error")).Once()
	posts, err := suite.postService.GetPosts(ctx, userID)

	suite.Equal(typesystem.ServerError, err)
	suite.Nil(posts)

	suite.mocksRepo.AssertExpectations(suite.T())
}
