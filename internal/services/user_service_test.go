package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Caixetadev/snippet/internal/entity"
	"github.com/Caixetadev/snippet/internal/mocks"
	"github.com/Caixetadev/snippet/pkg/typesystem"
	"github.com/google/uuid"
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
	tokenMaker         *mocks.Maker
}

func (suite *UserServiceTestSuite) SetupTest() {
	suite.mockRepo = new(mocks.UserRepository)
	suite.mockPasswordHasher = new(mocks.PasswordHasher)
	suite.validation = new(mocks.Validator)
	suite.tokenMaker = new(mocks.Maker)

	suite.userService = NewUserService(
		suite.mockRepo,
		suite.validation,
		suite.mockPasswordHasher,
		suite.tokenMaker,
	)
}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

func (suite *UserServiceTestSuite) TestCreate() {
	ctx := context.TODO()
	input := &entity.User{
		Name:     "John",
		Email:    "john@example.com",
		Password: "password",
	}

	suite.validation.On("Validate", input).Return(nil)
	suite.mockPasswordHasher.On("GenerateFromPassword", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("int")).Return([]byte("password_hashed"), nil)
	suite.mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.User")).Return(nil)

	user, err := suite.userService.Create(ctx, input)

	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), user.ID)
	assert.Equal(suite.T(), "password_hashed", user.Password)

	suite.mockRepo.AssertCalled(suite.T(), "Create", ctx, mock.AnythingOfType("*entity.User"))
	suite.validation.AssertCalled(suite.T(), "Validate", input)
	suite.mockPasswordHasher.AssertCalled(suite.T(), "GenerateFromPassword", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("int"))
}

func (suite *UserServiceTestSuite) TestValidationFails() {
	ctx := context.TODO()
	input := &entity.User{
		Name:  "John",
		Email: "john@example.com",
	}

	suite.validation.On("Validate", input).Return(errors.New("error"))
	suite.mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	_, err := suite.userService.Create(ctx, input)

	assert.Equal(suite.T(), typesystem.BadRequest, err)

	suite.mockRepo.AssertNotCalled(suite.T(), "Create", ctx, mock.AnythingOfType("*entity.User"))

	suite.validation.AssertNumberOfCalls(suite.T(), "Validate", 1)
}

func (suite *UserServiceTestSuite) TestCreateWithUserRepositoryFailure() {
	ctx := context.TODO()
	input := &entity.User{
		Name:     "John",
		Email:    "john@example.com",
		Password: "password",
	}

	suite.mockRepo.On("Create", ctx, mock.Anything).Return(errors.New("error"))

	suite.validation.On("Validate", input).Return(nil)

	suite.mockPasswordHasher.On("GenerateFromPassword", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("int")).Return([]byte("password_hashed"), nil)

	createdUser, err := suite.userService.Create(ctx, input)

	assert.Equal(suite.T(), typesystem.ServerError, err)
	assert.Nil(suite.T(), createdUser)

	suite.validation.AssertCalled(suite.T(), "Validate", input)
	suite.mockRepo.AssertCalled(suite.T(), "Create", ctx, mock.Anything)
	suite.mockPasswordHasher.AssertCalled(suite.T(), "GenerateFromPassword", mock.Anything, mock.AnythingOfType("int"))
}

func (suite *UserServiceTestSuite) TestCreate_ErrorOnHash() {
	ctx := context.TODO()
	input := &entity.User{
		Name:     "John",
		Email:    "john@example.com",
		Password: "password",
	}

	suite.validation.On("Validate", mock.Anything).Return(nil)
	suite.mockPasswordHasher.On("GenerateFromPassword", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("int")).Return([]byte(""), errors.New("error"))
	suite.mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.User")).Return(nil)

	_, err := suite.userService.Create(ctx, input)

	suite.Equal(typesystem.ServerError, err)
}

func (suite *UserServiceTestSuite) TestGetUserByEmailExists() {
	ctx := context.TODO()
	input := "jhon@example.com"
	output := &entity.User{
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

	suite.mockRepo.On("GetUserByEmail", ctx, input).Return(&entity.User{}, pgx.ErrNoRows)

	user, err := suite.userService.GetUserByEmail(ctx, input)

	suite.Assert().Nil(user)
	suite.Assert().Equal(err, typesystem.Unauthorized)

	suite.mockRepo.AssertCalled(suite.T(), "GetUserByEmail", ctx, input)
}

func (suite *UserServiceTestSuite) TestGetUserByEmailWithRepositoryFails() {
	ctx := context.TODO()
	input := "jhon@example.com"

	suite.mockRepo.On("GetUserByEmail", ctx, input).Return(&entity.User{}, errors.New("error"))

	user, err := suite.userService.GetUserByEmail(ctx, input)

	suite.Assert().Nil(user)
	suite.Assert().Equal(err, typesystem.ServerError)

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

func (suite *UserServiceTestSuite) TestUpdatePassword() {
	ctx := context.TODO()

	id := uuid.New()
	password := "password"
	passwordConfirm := "password"

	passwordHashed := "password_hashed"

	suite.mockPasswordHasher.On("GenerateFromPassword", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("int")).Return([]byte(passwordHashed), nil)
	suite.mockRepo.On("UpdatePassword", ctx, passwordHashed, id.String()).Return(nil)

	err := suite.userService.UpdatePassword(ctx, password, passwordConfirm, id.String())

	suite.Nil(err)
	suite.mockRepo.AssertCalled(suite.T(), "UpdatePassword", ctx, passwordHashed, id.String())
	suite.mockPasswordHasher.AssertCalled(suite.T(), "GenerateFromPassword", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("int"))
}

func (suite *UserServiceTestSuite) TestUpdatePassword_PasswordNotMatch() {
	ctx := context.TODO()

	id := uuid.New()
	password := "password"
	passwordConfirm := "not_password"

	err := suite.userService.UpdatePassword(ctx, password, passwordConfirm, id.String())

	suite.Assert().Equal(err, typesystem.BadRequest)
}

func (suite *UserServiceTestSuite) TestUpdatePassword_HashError() {
	ctx := context.TODO()

	id := uuid.New()
	password := "password"
	passwordConfirm := "password"

	suite.mockPasswordHasher.On("GenerateFromPassword", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("int")).Return([]byte(""), errors.New("error"))

	err := suite.userService.UpdatePassword(ctx, password, passwordConfirm, id.String())

	suite.Equal(typesystem.ServerError, err)
}

func (suite *UserServiceTestSuite) TestUpdatePassword_RepoError() {
	ctx := context.TODO()

	id := uuid.New()
	password := "password"
	passwordConfirm := "password"

	passwordHashed := "password_hashed"

	suite.mockPasswordHasher.On("GenerateFromPassword", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("int")).Return([]byte(passwordHashed), nil)
	suite.mockRepo.On("UpdatePassword", ctx, passwordHashed, id.String()).Return(errors.New("error"))

	err := suite.userService.UpdatePassword(ctx, password, passwordConfirm, id.String())

	suite.Assert().Equal(typesystem.ServerError, err)
	suite.mockRepo.AssertCalled(suite.T(), "UpdatePassword", ctx, passwordHashed, id.String())
	suite.mockPasswordHasher.AssertCalled(suite.T(), "GenerateFromPassword", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("int"))
}

func (suite *UserServiceTestSuite) TestVerifyCodeToResetPassword() {
	ctx := context.TODO()

	output := entity.VerificationData{
		Code:      "code",
		ExpiresAt: time.Now().Add(time.Minute * 5),
		Email:     "email",
		UserID:    uuid.New(),
	}

	suite.mockRepo.On("VerifyCodeToResetPassword", ctx, "code").Return(output, nil)

	userID, err := suite.userService.VerifyCodeToResetPassword(ctx, "code")

	suite.Nil(err)
	suite.Equal(output.UserID.String(), userID)
}

func (suite *UserServiceTestSuite) TestVerifyCodeToResetPassword_CodeExpired() {
	ctx := context.TODO()

	output := entity.VerificationData{
		Code:      "code",
		ExpiresAt: time.Now().Add(-1),
		Email:     "email",
		UserID:    uuid.New(),
	}

	suite.mockRepo.On("VerifyCodeToResetPassword", ctx, "code").Return(output, nil)

	userID, err := suite.userService.VerifyCodeToResetPassword(ctx, "code")

	suite.Equal("", userID)
	suite.NotNil(err)
	suite.Equal(err, typesystem.TokenExpiredError)
}

func (suite *UserServiceTestSuite) TestVerifyCodeToResetPassword_Error() {
	ctx := context.TODO()

	suite.mockRepo.On("VerifyCodeToResetPassword", ctx, "code").Return(entity.VerificationData{}, errors.New("error"))

	userID, err := suite.userService.VerifyCodeToResetPassword(ctx, "code")

	suite.Equal("", userID)
	suite.NotNil(err)
	suite.Equal(err, typesystem.ServerError)
}

func (suite *UserServiceTestSuite) TestVerifyToken() {
	ctx := context.TODO()

	output := &entity.Payload{}

	suite.tokenMaker.On("VerifyToken", "token").Return(output, nil)

	refreshToken, err := suite.userService.VerifyToken(ctx, "token")

	suite.Nil(err)
	suite.NotNil(refreshToken)
}

func (suite *UserServiceTestSuite) TestVerifyToken_Invalid() {
	ctx := context.TODO()

	output := &entity.Payload{}

	suite.tokenMaker.On("VerifyToken", "token").Return(output, errors.New("error"))

	accessToken, err := suite.userService.VerifyToken(ctx, "token")

	suite.Error(err)
	suite.Nil(accessToken)
}
