package services

import (
	"context"
	"errors"

	"github.com/Caixetadev/snippet/internal/core/domain"
	"github.com/Caixetadev/snippet/internal/validation"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository domain.UserRepository
	validation     validation.Validator
}

func NewUserService(userRepository domain.UserRepository, validation validation.Validator) *UserService {
	return &UserService{
		userRepository: userRepository,
		validation:     validation,
	}
}

func (su *UserService) Create(ctx context.Context, input *domain.User) error {
	if err := su.validation.Validate(input); err != nil {
		// return nil, errors.BadRequest("")
		return errors.New("BadRequest")
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		// return nil, errors.InternalServerError("")
	}

	user := domain.User{
		ID:       uuid.New(),
		Name:     input.Name,
		Email:    input.Email,
		Password: string(encryptedPassword),
	}

	err = su.userRepository.Create(ctx, &user)
	if err != nil {
		return err
	}

	return nil
}

func (ls *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return ls.userRepository.GetUserByEmail(ctx, email)
}

func (ls *UserService) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	return ls.userRepository.UserExistsByEmail(ctx, email)
}