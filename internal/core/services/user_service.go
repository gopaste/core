package services

import (
	"context"

	"github.com/Caixetadev/snippet/internal/core/domain"
	apperr "github.com/Caixetadev/snippet/internal/core/error"
	"github.com/Caixetadev/snippet/internal/token"
	"github.com/Caixetadev/snippet/pkg/validation"
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

func (us *UserService) Create(ctx context.Context, input *domain.User) (*domain.User, error) {
	if err := us.validation.Validate(input); err != nil {
		return nil, apperr.BadRequest
	}

	encryptedPassword, err := us.passwordHasher.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperr.ServerError
	}

	user := domain.NewUser(input.Name, input.Email, string(encryptedPassword))

	err = us.userRepository.Create(ctx, user)
	if err != nil {
		return nil, apperr.ServerError
	}

	return user, nil
}

func (us *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := us.userRepository.GetUserByEmail(ctx, email)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperr.Unauthorized
		}

		return nil, apperr.ServerError
	}

	return user, nil
}

func (us *UserService) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	return us.userRepository.UserExistsByEmail(ctx, email)
}

func (us *UserService) CreateAccessToken(user *domain.User, secret string, expiry int) (accesstoken string, err error) {
	return token.CreateAccessToken(user, secret, expiry)
}

func (us *UserService) CompareHashAndPassword(passwordInDatabase, passwordRequest string) error {
	return us.passwordHasher.CompareHashAndPassword([]byte(passwordInDatabase), []byte(passwordRequest))
}
