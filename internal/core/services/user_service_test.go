package services

import (
	"context"
	"testing"

	"github.com/Caixetadev/snippet/internal/core/domain"
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
}
