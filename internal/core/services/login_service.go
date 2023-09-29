package services

import (
	"github.com/Caixetadev/snippet/internal/core/domain"
	"github.com/Caixetadev/snippet/internal/validation"
)

type loginService struct {
	userRepository domain.UserRepository
	validation     validation.Validator
}

func NewLoginService(userRepository domain.UserRepository, validation validation.Validator) *loginService {
	return &loginService{
		userRepository: userRepository,
		validation:     validation,
	}
}

func (ls *loginService) GetUserByEmail(email string) (*domain.User, error) {
	return nil, nil
}
