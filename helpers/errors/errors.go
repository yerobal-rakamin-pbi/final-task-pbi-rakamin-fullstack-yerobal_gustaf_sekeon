package errors

import (
	"net/http"
	"reflect"
)

type Errors struct {
	Type    string
	Code    int64
	Message interface{}
}

const (
	NotFoundType            = "HTTPStatusNotFound"
	InternalServerErrorType = "HTTPStatusInternalServerError"
	BadRequestType          = "HTTPStatusBadRequest"
	UnauthorizedType        = "HTTPStatusUnauthorized"
	RequestTimeoutType      = "HTTPStatusRequestTimeout"
	UnprocessableEntityType = "HTTPStatusUnprocessableEntity"
	ConflictType            = "HTTPStatusConflict"
)

func (e *Errors) Error() string {
	return e.Message.(string)
}

func NewWithCode(code int64, message interface{}, errType string) error {
	errors := &Errors{
		Type:    errType,
		Code:    code,
		Message: message,
	}

	return errors
}

func NotFound(entity string) error {
	return NewWithCode(http.StatusNotFound, entity, NotFoundType)
}

func InternalServerError(message string) error {
	return NewWithCode(http.StatusInternalServerError, message, InternalServerErrorType)
}

func BadRequest(message string) error {
	return NewWithCode(http.StatusBadRequest, message, BadRequestType)
}

func Unauthorized(message string) error {
	return NewWithCode(http.StatusUnauthorized, message, UnauthorizedType)
}

func RequestTimeout(message string) error {
	return NewWithCode(http.StatusRequestTimeout, message, RequestTimeoutType)
}

func ValidationError(message interface{}) error {
	return NewWithCode(http.StatusUnprocessableEntity, message, UnprocessableEntityType)
}

func Conflict(message string) error {
	return NewWithCode(http.StatusConflict, message, ConflictType)
}

func GetType(err error) string {
	if err == nil {
		return "HTTPStatusOK"
	}

	if reflect.TypeOf(err).String() == "*errors.Errors" {
		return err.(*Errors).Type
	}

	return InternalServerErrorType
}

func GetCode(err error) int64 {
	if err == nil {
		return 200
	}

	if reflect.TypeOf(err).String() == "*errors.Errors" {
		return err.(*Errors).Code
	}

	return 500
}

func GetMessage(err error) interface{} {
	if err == nil {
		return "OK"
	}

	if reflect.TypeOf(err).String() == "*errors.Errors" {
		return err.(*Errors).Message
	}

	return err.Error()
}
