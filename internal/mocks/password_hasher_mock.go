package mocks

import "github.com/stretchr/testify/mock"

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
