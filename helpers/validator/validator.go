package validator

import (
	validatorLib "github.com/go-playground/validator/v10"
)

type validator struct {
	validatorLib *validatorLib.Validate
}

type Interface interface {
	ValidateStruct(data interface{}) error
}

func Init() Interface {
	return &validator{validatorLib: validatorLib.New()}
}

func (v *validator) ValidateStruct(data interface{}) error {
	return v.validatorLib.Struct(data)
}
