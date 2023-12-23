package validation

import (
	validatorv10 "github.com/go-playground/validator/v10"
)

// Validator defines the validation interface
type Validator interface {
	Validate(obj interface{}) error
}

// validator implements Validator
type validator struct {
	validatorv10 *validatorv10.Validate
}

// NewValidator creates a new Validator
func NewValidator(validatorv10 *validatorv10.Validate) Validator {
	return &validator{validatorv10}
}

// Validate implements Validator
func (v *validator) Validate(obj interface{}) error {
	return v.validatorv10.Struct(obj)
}
