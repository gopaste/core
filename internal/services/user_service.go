package services

import (
	"context"
	"fmt"
	"net/http"
	"time"

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
	tokenMaker     token.Maker
}

func NewUserService(
	userRepository entity.UserRepository,
	validation validation.Validator,
	passwordHasher entity.PasswordHasher,
	tokenMaker token.Maker,
) *UserService {
	return &UserService{
		userRepository: userRepository,
		validation:     validation,
		passwordHasher: passwordHasher,
		tokenMaker:     tokenMaker,
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

func (us *UserService) CreateAccessToken(user *entity.User, expiry time.Duration) (string, *entity.Payload, error) {
	token, payload, err := us.tokenMaker.CreateToken(user, expiry)
	if err != nil {
		return "", nil, err
	}

	return token, payload, nil
}

func (us *UserService) CreateRefreshToken(ctx context.Context, user *entity.User, expiry time.Duration) (string, *entity.Payload, error) {
	token, payload, err := us.tokenMaker.CreateToken(user, expiry)
	if err != nil {
		return "", nil, err
	}

	return token, payload, nil
}

func (us *UserService) CreateSession(ctx context.Context, payload *entity.Payload, token string) error {
	session := entity.NewSession(ctx, payload, token)

	err := us.userRepository.CreateSession(ctx, session)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) CompareHashAndPassword(passwordInDatabase, passwordRequest string) error {
	return us.passwordHasher.CompareHashAndPassword([]byte(passwordInDatabase), []byte(passwordRequest))
}

func (us *UserService) StoreVerificationData(ctx context.Context, userID uuid.UUID, email string, code string) error {
	verificationData := entity.NewVerificationData(userID, email, code)

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

func (us *UserService) GetSession(ctx context.Context, id string) (*entity.Session, error) {
	user, err := us.userRepository.GetSession(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) VerifyToken(ctx context.Context, token string) (*entity.Payload, error) {
	refreshToken, err := us.tokenMaker.VerifyToken(token)
	if err != nil {
		return nil, entity.NewHttpError(err.Error(), "Token verification failed", http.StatusUnauthorized)
	}

	return refreshToken, nil
}
