package services

import (
	"context"
	"errors"
	"testing"

	"github.com/Caixetadev/snippet/internal/core/domain"
	apperr "github.com/Caixetadev/snippet/internal/core/error"
	"github.com/Caixetadev/snippet/internal/mocks"
	"github.com/Caixetadev/snippet/internal/validation"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	validatorv10 "github.com/go-playground/validator/v10"
)

func TestCreate(t *testing.T) {
	t.Run("should create user", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)

		userService := NewUserService(
			mockRepo,
			validation.NewValidator(validatorv10.New()),
			&domain.BcryptPasswordHasher{},
		)

		ctx := context.TODO()
		input := &domain.User{
			Name:     "John",
			Email:    "john@example.com",
			Password: "password",
		}

		mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(nil)

		err := userService.Create(ctx, input)

		assert.Nil(t, err)

		mockRepo.AssertCalled(t, "Create", ctx, mock.AnythingOfType("*domain.User"))
	})

	t.Run("should return BadRequest when validation fails", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)

		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

		userService := NewUserService(mockRepo, validation.NewValidator(validatorv10.New()), nil)

		ctx := context.TODO()
		input := &domain.User{
			Name:  "John",
			Email: "john@example.com",
		}

		err := userService.Create(ctx, input)

		assert.Equal(t, apperr.BadRequest, err)
	})

	t.Run("should return ServerError when userRepository fails", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)

		mockRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New(""))

		userService := NewUserService(
			mockRepo,
			validation.NewValidator(validatorv10.New()),
			&domain.BcryptPasswordHasher{},
		)

		ctx := context.TODO()
		input := &domain.User{
			Name:     "John",
			Email:    "john@example.com",
			Password: "123",
		}

		err := userService.Create(ctx, input)

		assert.Equal(t, apperr.ServerError, err)
	})

	t.Run("should encrypt password using BcryptPasswordHasher", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)

		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

		userService := NewUserService(
			mockRepo,
			validation.NewValidator(validatorv10.New()),
			&domain.BcryptPasswordHasher{},
		)

		ctx := context.TODO()
		input := &domain.User{
			Name:     "John",
			Email:    "john@example.com",
			Password: "password",
		}

		err := userService.Create(ctx, input)

		assert.Nil(t, err)

		assert.NotEqual(t, "password", input.Password)
	})

	t.Run("should return ServerError when password encryption fails", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)

		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

		mockPasswordHasher := new(mocks.PasswordHasher)
		mockPasswordHasher.On("GenerateFromPassword", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("int")).
			Return([]byte(""), errors.New("error"))

		userService := NewUserService(mockRepo, validation.NewValidator(validatorv10.New()), mockPasswordHasher)

		ctx := context.TODO()
		input := &domain.User{
			Name:     "John",
			Email:    "john@example.com",
			Password: "password",
		}

		err := userService.Create(ctx, input)

		assert.NotNil(t, err)
		assert.Equal(t, apperr.ServerError, err)
	})
}

func TestUserExistsByEmail(t *testing.T) {
	t.Run("should return true if the user with the given email exists", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)

		mockRepo.On("UserExistsByEmail", mock.Anything, "test@example.com").Return(true, nil)

		userService := NewUserService(mockRepo, validation.NewValidator(validatorv10.New()), nil)

		ctx := context.TODO()
		email := "test@example.com"

		exists, err := userService.UserExistsByEmail(ctx, email)

		assert.Nil(t, err)
		assert.True(t, exists)
	})

	t.Run("should return false if the user with the specified email does not exist", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)

		mockRepo.On("UserExistsByEmail", mock.Anything, "test@example.com").Return(false, nil)

		userService := NewUserService(mockRepo, validation.NewValidator(validatorv10.New()), nil)

		ctx := context.TODO()
		email := "test@example.com"

		exists, err := userService.UserExistsByEmail(ctx, email)

		assert.Nil(t, err)
		assert.False(t, exists)
	})

	t.Run("should return an error when the userRepository fails", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)
		mockRepo.On("UserExistsByEmail", mock.Anything, "test@example.com").Return(false, errors.New("database error"))

		userService := NewUserService(
			mockRepo,
			validation.NewValidator(validatorv10.New()),
			&domain.BcryptPasswordHasher{},
		)

		ctx := context.TODO()
		email := "test@example.com"

		exists, err := userService.UserExistsByEmail(ctx, email)

		assert.NotNil(t, err)
		assert.False(t, exists)
	})
}

func TestGetUserByEmail(t *testing.T) {
	t.Run("should return a user with the correct email", func(t *testing.T) {
		repoMock := new(mocks.UserRepository)

		expectedUser := &domain.User{
			Name:     "test",
			Email:    "test@example.com",
			Password: "123",
		}

		repoMock.On("GetUserByEmail", mock.Anything, "test@example.com").Return(expectedUser, nil)

		userService := &UserService{
			userRepository: repoMock,
		}

		ctx := context.TODO()
		email := "test@example.com"
		user, err := userService.GetUserByEmail(ctx, email)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)

		repoMock.AssertCalled(t, "GetUserByEmail", mock.Anything, "test@example.com")

		repoMock.AssertExpectations(t)
	})

	t.Run("should return error 404 if the user with the specified email does not exist", func(t *testing.T) {
		repoMock := new(mocks.UserRepository)

		repoMock.On("GetUserByEmail", mock.Anything, "test@example.com").Return((*domain.User)(nil), pgx.ErrNoRows)

		userService := &UserService{
			userRepository: repoMock,
		}

		ctx := context.TODO()
		email := "test@example.com"
		user, err := userService.GetUserByEmail(ctx, email)

		assert.Nil(t, user)
		assert.Equal(t, apperr.NotFound, err)

		repoMock.AssertCalled(t, "GetUserByEmail", mock.Anything, "test@example.com")

		repoMock.AssertExpectations(t)
	})

	t.Run("should return ServerError when the GetUserByEmail repository fails", func(t *testing.T) {
		repoMock := new(mocks.UserRepository)

		repoMock.On("GetUserByEmail", mock.Anything, "test@example.com").Return((*domain.User)(nil), errors.New("error"))

		userService := &UserService{
			userRepository: repoMock,
		}

		ctx := context.TODO()
		email := "test@example.com"
		user, err := userService.GetUserByEmail(ctx, email)

		assert.Nil(t, user)
		assert.Equal(t, apperr.ServerError, err)

		repoMock.AssertCalled(t, "GetUserByEmail", mock.Anything, "test@example.com")

		repoMock.AssertExpectations(t)
	})
}
