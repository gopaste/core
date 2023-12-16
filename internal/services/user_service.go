package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Caixetadev/snippet/internal/entity"
	"github.com/Caixetadev/snippet/internal/token"
	"github.com/Caixetadev/snippet/pkg/validation"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository entity.UserRepository
	validation     validation.Validator
	passwordHasher entity.PasswordHasher
}

func NewUserService(
	userRepository entity.UserRepository,
	validation validation.Validator,
	passwordHasher entity.PasswordHasher,
) *UserService {
	return &UserService{
		userRepository: userRepository,
		validation:     validation,
		passwordHasher: passwordHasher,
	}
}

func (us *UserService) Create(ctx context.Context, input *entity.User) (*entity.User, error) {
	if err := us.validation.Validate(input); err != nil {
		return nil, entity.BadRequest
	}

	encryptedPassword, err := us.passwordHasher.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, entity.ServerError
	}

	user := entity.NewUser(input.Name, input.Email, string(encryptedPassword))

	err = us.userRepository.Create(ctx, user)
	if err != nil {
		return nil, entity.ServerError
	}

	return user, nil
}

func (us *UserService) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := us.userRepository.GetUserByEmail(ctx, email)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, entity.Unauthorized
		}

		return nil, entity.ServerError
	}

	return user, nil
}

func (us *UserService) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	return us.userRepository.UserExistsByEmail(ctx, email)
}

func (us *UserService) CreateAccessToken(user *entity.User, secret string, expiry int) (accesstoken string, err error) {
	return token.CreateAccessToken(user, secret, expiry)
}

func (us *UserService) CompareHashAndPassword(passwordInDatabase, passwordRequest string) error {
	return us.passwordHasher.CompareHashAndPassword([]byte(passwordInDatabase), []byte(passwordRequest))
}

func (us *UserService) StoreVerificationData(ctx context.Context, userID uuid.UUID, email string, code string) error {
	verificationData := &entity.VerificationData{
		UserID: userID,
		Email:  email,
		Code:   code,
	}

	return us.userRepository.StoreVerificationData(ctx, verificationData)
}

func (us *UserService) VerifyCodeToResetPassword(ctx context.Context, code string) (string, bool, error) {
	userID, valid, err := us.userRepository.VerifyCodeToResetPassword(ctx, code)
	if !valid {
		return "", false, entity.NewHttpError("Invalid or expired token", "The token provided for password recovery is invalid or has expired.", http.StatusUnauthorized)
	}

	if err != nil {
		fmt.Println(err)
		return "", false, entity.ServerError
	}

	return userID, true, nil
}

func (us *UserService) UpdatePassword(ctx context.Context, password, passwordConfirmation, id string) error {
	if password != passwordConfirmation {
		return entity.BadRequest
	}

	encryptedPassword, err := us.passwordHasher.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return entity.ServerError
	}

	return us.userRepository.UpdatePassword(ctx, string(encryptedPassword), id)
}
