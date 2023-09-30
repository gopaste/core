package services

import (
	"context"
	"errors"
	"testing"

	"github.com/Caixetadev/snippet/internal/core/domain"
	apperr "github.com/Caixetadev/snippet/internal/core/error"
	"github.com/Caixetadev/snippet/internal/mocks"
	"github.com/Caixetadev/snippet/internal/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	validatorv10 "github.com/go-playground/validator/v10"
)

func TestCreate(t *testing.T) {
	t.Run("should create user", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)

		userService := NewUserService(mockRepo, validation.NewValidator(validatorv10.New()), &domain.BcryptPasswordHasher{})

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

		userService := NewUserService(mockRepo, validation.NewValidator(validatorv10.New()), &domain.BcryptPasswordHasher{})

		ctx := context.TODO()
		input := &domain.User{
			Name:     "John",
			Email:    "john@example.com",
			Password: "123",
		}

		err := userService.Create(ctx, input)

		assert.Equal(t, apperr.ServerError, err)
	})
}
