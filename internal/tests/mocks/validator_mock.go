package mocks

import (
	"github.com/stretchr/testify/mock"
)

type Validator struct {
	mock.Mock
}

func (v *Validator) Validate(obj interface{}) error {
	arg := v.Called(obj)
	return arg.Error(0)
}
