package mocks

import (
	"time"

	"github.com/Caixetadev/snippet/internal/entity"
	"github.com/stretchr/testify/mock"
)

type Maker struct {
	mock.Mock
}

func (m *Maker) CreateToken(user *entity.User, expiry time.Duration) (string, *entity.Payload, error) {
	args := m.Called(user, expiry)
	return args.String(0), args.Get(1).(*entity.Payload), args.Error(2)
}

func (m *Maker) VerifyToken(token string) (*entity.Payload, error) {
	args := m.Called(token)
	return args.Get(0).(*entity.Payload), args.Error(1)
}
