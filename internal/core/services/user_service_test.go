package services

import (
	"context"
	"errors"
	"testing"

	"github.com/Caixetadev/snippet/internal/core/domain"
	apperr "github.com/Caixetadev/snippet/internal/core/error"
	"github.com/Caixetadev/snippet/internal/mocks"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserServiceTestSuite struct {
	suite.Suite
	mockRepo           *mocks.UserRepository
	mockPasswordHasher *mocks.PasswordHasher
	userService        *UserService
	validation         *mocks.Validator
}

func (suite *UserServiceTestSuite) SetupTest() {
	suite.mockRepo = new(mocks.UserRepository)
	suite.mockPasswordHasher = new(mocks.PasswordHasher)
	suite.validation = new(mocks.Validator)

	suite.userService = NewUserService(
		suite.mockRepo,
		suite.validation,
		suite.mockPasswordHasher,
	)
}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

func (suite *UserServiceTestSuite) TestCreate() {
	ctx := context.TODO()
	input := &domain.User{
		Name:     "John",
		Email:    "john@example.com",
		Password: "password",
	}

	suite.validation.On("Validate", input).Return(nil)
	suite.mockPasswordHasher.On("GenerateFromPassword", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("int")).Return([]byte("password_hashed"), nil)
	suite.mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(nil)

	err := suite.userService.Create(ctx, input)

	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), input.ID)
	assert.Equal(suite.T(), "password_hashed", input.Password)

	suite.mockRepo.AssertCalled(suite.T(), "Create", ctx, mock.AnythingOfType("*domain.User"))
	suite.validation.AssertCalled(suite.T(), "Validate", input)
	suite.mockPasswordHasher.AssertCalled(suite.T(), "GenerateFromPassword", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("int"))
}

func (suite *UserServiceTestSuite) TestValidationFails() {
	ctx := context.TODO()
	input := &domain.User{
		Name:  "John",
		Email: "john@example.com",
	}

	suite.validation.On("Validate", input).Return(errors.New("error"))
	suite.mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	err := suite.userService.Create(ctx, input)

	assert.Equal(suite.T(), apperr.BadRequest, err)

	suite.mockRepo.AssertNotCalled(suite.T(), "Create", ctx, mock.AnythingOfType("*domain.User"))

	suite.validation.AssertNumberOfCalls(suite.T(), "Validate", 1)
}

func (suite *UserServiceTestSuite) TestCreateWithUserRepositoryFailure() {
	ctx := context.TODO()
	input := &domain.User{
		Name:     "John",
		Email:    "john@example.com",
		Password: "password",
	}

	inputGenerateFromPassword := input.Password

	suite.mockRepo.On("Create", ctx, input).Return(errors.New("error"))

	suite.validation.On("Validate", input).Return(nil)

	suite.mockPasswordHasher.On("GenerateFromPassword", []byte(inputGenerateFromPassword), 10).Return([]byte("password_hashed"), nil)

	err := suite.userService.Create(ctx, input)

	assert.Equal(suite.T(), apperr.ServerError, err)

	suite.validation.AssertCalled(suite.T(), "Validate", input)
	suite.mockRepo.AssertCalled(suite.T(), "Create", ctx, input)
	suite.mockPasswordHasher.AssertCalled(suite.T(), "GenerateFromPassword", []byte(inputGenerateFromPassword), 10)
}

func (suite *UserServiceTestSuite) TestCreate_ErrorOnHash() {
	ctx := context.TODO()
	input := &domain.User{
		Name:     "John",
		Email:    "john@example.com",
		Password: "password",
	}

	suite.validation.On("Validate", mock.Anything).Return(nil)
	suite.mockPasswordHasher.On("GenerateFromPassword", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("int")).Return([]byte(""), errors.New("error"))
	suite.mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(nil)

	err := suite.userService.Create(ctx, input)

	suite.Equal(apperr.ServerError, err)
}

func (suite *UserServiceTestSuite) TestGetUserByEmailExists() {
	ctx := context.TODO()
	input := "jhon@example.com"
	output := &domain.User{
		Name:     "John",
		Email:    "john@example.com",
		Password: "password",
	}

	suite.mockRepo.On("GetUserByEmail", ctx, input).Return(output, nil)

	user, err := suite.userService.GetUserByEmail(ctx, input)

	suite.Assert().Nil(err)
	suite.Assert().Equal(output, user)

	suite.mockRepo.AssertCalled(suite.T(), "GetUserByEmail", ctx, input)
}

func (suite *UserServiceTestSuite) TestGetUserByEmailNotExists() {
	ctx := context.TODO()
	input := "jhon@example.com"

	suite.mockRepo.On("GetUserByEmail", ctx, input).Return(&domain.User{}, pgx.ErrNoRows)

	user, err := suite.userService.GetUserByEmail(ctx, input)

	suite.Assert().Nil(user)
	suite.Assert().Equal(err, apperr.NotFound)

	suite.mockRepo.AssertCalled(suite.T(), "GetUserByEmail", ctx, input)
}

func (suite *UserServiceTestSuite) TestGetUserByEmailWithRepositoryFails() {
	ctx := context.TODO()
	input := "jhon@example.com"

	suite.mockRepo.On("GetUserByEmail", ctx, input).Return(&domain.User{}, errors.New("error"))

	user, err := suite.userService.GetUserByEmail(ctx, input)

	suite.Assert().Nil(user)
	suite.Assert().Equal(err, apperr.ServerError)

	suite.mockRepo.AssertCalled(suite.T(), "GetUserByEmail", ctx, input)
}

func (suite *UserServiceTestSuite) TestUserExistsByEmail() {
	ctx := context.TODO()

	input := "user@gmail.com"

	suite.mockRepo.On("UserExistsByEmail", ctx, input).Return(true, nil)

	exists, err := suite.userService.UserExistsByEmail(ctx, input)

	suite.True(exists)
	suite.Nil(err)

	suite.mockRepo.AssertCalled(suite.T(), "UserExistsByEmail", ctx, input)
}

func (suite *UserServiceTestSuite) TestUserExistsByEmail_Error() {
	ctx := context.TODO()

	input := "user@gmail.com"

	suite.mockRepo.On("UserExistsByEmail", ctx, input).Return(false, errors.New("error"))

	exists, err := suite.userService.UserExistsByEmail(ctx, input)

	suite.False(exists)
	suite.NotNil(err)

	suite.mockRepo.AssertCalled(suite.T(), "UserExistsByEmail", ctx, input)
}
