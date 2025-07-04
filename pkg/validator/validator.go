package validator

import (
	"github.com/go-playground/validator/v10"
)

type Validator interface {
	Validate(i interface{}) error
}

type DefaultValidator struct {
	validate *validator.Validate
}

func NewValidator() Validator {
	return &DefaultValidator{
		validate: validator.New(),
	}
}

func (v *DefaultValidator) Validate(i interface{}) error {
	return v.validate.Struct(i)
}
