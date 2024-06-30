package validator

import (
	"fmt"

	validatorLib "github.com/go-playground/validator/v10"
)

type validator struct {
	validatorLib *validatorLib.Validate
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type Interface interface {
	ValidateStruct(data interface{}) error
	GetValidationErrors(err error) []ValidationError
}

func Init() Interface {
	return &validator{validatorLib: validatorLib.New()}
}

func (v *validator) ValidateStruct(data interface{}) error {
	return v.validatorLib.Struct(data)
}

func (v *validator) GetValidationErrors(err error) []ValidationError {
	var errors []ValidationError
	for _, err := range err.(validatorLib.ValidationErrors) {
		errors = append(errors, ValidationError{
			Field:   err.Field(),
			Message: v.messageForTag(err),
		})
	}
	return errors
}

func (v *validator) messageForTag(err validatorLib.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("The %s field is required.", err.Field())
	case "min":
		return fmt.Sprintf("The %s field must be at least %s characters.", err.Field(), err.Param())
	}
	return fmt.Sprintf("The %s field is invalid.", err.Field())
}
