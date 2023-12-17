package mocks

import (
	"context"

	"github.com/Caixetadev/snippet/internal/entity"
	"github.com/stretchr/testify/mock"
)

type UserRepository struct {
	mock.Mock
}

func (m *UserRepository) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *UserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *UserRepository) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *UserRepository) StoreVerificationData(ctx context.Context, verification *entity.VerificationData) error {
	args := m.Called(ctx, verification)
	return args.Error(0)
}

func (m *UserRepository) UpdatePassword(ctx context.Context, password string, id string) error {
	args := m.Called(ctx, password, id)
	return args.Error(0)
}

func (m *UserRepository) VerifyCodeToResetPassword(ctx context.Context, code string) (string, bool, error) {
	args := m.Called(ctx, code)
	return args.String(0), args.Bool(1), args.Error(2)
}

type PasswordHasher struct {
	mock.Mock
}

func (m *PasswordHasher) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	args := m.Called(password, cost)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *PasswordHasher) CompareHashAndPassword(hashedPassword, password []byte) error {
	args := m.Called(hashedPassword, password)
	return args.Error(1)
}

type Validator struct {
	mock.Mock
}

func (v *Validator) Validate(obj interface{}) error {
	arg := v.Called(obj)
	return arg.Error(0)
}
