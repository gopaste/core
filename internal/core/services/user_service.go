package services

import (
	"context"

	"github.com/Caixetadev/snippet/internal/core/domain"
	apperr "github.com/Caixetadev/snippet/internal/core/error"
	"github.com/Caixetadev/snippet/internal/validation"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository domain.UserRepository
	validation     validation.Validator
	passwordHasher domain.PasswordHasher
}

func NewUserService(
	userRepository domain.UserRepository,
	validation validation.Validator,
	passwordHasher domain.PasswordHasher,
) *UserService {
	return &UserService{
		userRepository: userRepository,
		validation:     validation,
		passwordHasher: passwordHasher,
	}
}

func (su *UserService) Create(ctx context.Context, input *domain.User) error {
	if err := su.validation.Validate(input); err != nil {
		return apperr.BadRequest
	}

	encryptedPassword, err := su.passwordHasher.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return apperr.ServerError
	}

	input.ID = uuid.New()
	input.Password = string(encryptedPassword)

	err = su.userRepository.Create(ctx, input)
	if err != nil {
		return apperr.ServerError
	}

	return nil
}

func (ls *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := ls.userRepository.GetUserByEmail(ctx, email)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperr.NotFound
		}

		return nil, apperr.ServerError
	}

	return user, nil
}

func (ls *UserService) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	return ls.userRepository.UserExistsByEmail(ctx, email)
}
