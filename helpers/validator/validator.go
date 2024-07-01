package validator

import (
	"fmt"
	"strings"

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
	GetValidationErrors(err error) ([]ValidationError, string)
}

func Init() Interface {
	return &validator{validatorLib: validatorLib.New()}
}

func (v *validator) ValidateStruct(data interface{}) error {
	return v.validatorLib.Struct(data)
}

func (v *validator) GetValidationErrors(err error) ([]ValidationError, string) {
	var errors []ValidationError
	var errorList []string
	for i, err := range err.(validatorLib.ValidationErrors) {
		errorList = append(errorList, fmt.Sprintf("%d.%s", i+1, v.messageForTag(err)))
		errors = append(errors, ValidationError{
			Field:   err.Field(),
			Message: errorList[i],
		})
	}

	errorMsg := strings.Join(errorList, " ")
	return errors, errorMsg
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
