package mocks

import (
	"context"

	"github.com/Caixetadev/snippet/internal/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type UserRepository struct {
	mock.Mock
}

func (m *UserRepository) Insert(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *UserRepository) FindOneByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *UserRepository) FindOneByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	args := m.Called(ctx, id)
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

func (m *UserRepository) UpdatePassword(ctx context.Context, password string, id uuid.UUID) error {
	args := m.Called(ctx, password, id)
	return args.Error(0)
}

func (m *UserRepository) VerifyCodeToResetPassword(ctx context.Context, code string) (entity.VerificationData, error) {
	args := m.Called(ctx, code)
	return args.Get(0).(entity.VerificationData), args.Error(1)
}

func (m *UserRepository) CreateSession(ctx context.Context, session *entity.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *UserRepository) GetRefreshTokenByToken(ctx context.Context, token string) (*entity.RefreshToken, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(*entity.RefreshToken), args.Error(1)
}

func (m *UserRepository) GetSession(ctx context.Context, id uuid.UUID) (*entity.Session, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entity.Session), args.Error(1)
}
