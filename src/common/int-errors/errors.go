package int_errors

import (
	"fmt"
	"reflect"
)

var errNotFoundType = reflect.TypeOf(ErrNotFound{})

type ErrNotFound struct {
	message string
}

func NewErrNotFound(message string, args ...any) error {
	return &ErrNotFound{fmt.Sprintf(message, args...)}
}

func (e ErrNotFound) Error() string {
	return e.message
}

func (e ErrNotFound) Is(err error) bool {
	return reflect.TypeOf(err) == errNotFoundType
}

type ErrorResponseObject struct {
	Error    string
	Messages []string
}

func NewBuilder() ErrorResponseBuilder {
	return &ErrorResponseObject{}
}

func New() ErrorResponse {
	return &ErrorResponseObject{}
}

func NewWithDetails(error string, messages ...string) ErrorResponse {
	return &ErrorResponseObject{
		Error:    error,
		Messages: messages,
	}
}
func NewWithError(err error, messages ...string) ErrorResponse {
	return &ErrorResponseObject{
		Error:    err.Error(),
		Messages: messages,
	}
}

type ErrorResponseBuilder interface {
	AddError(error string)
	AddErrorIfMessagesExists(error string)
	AddMessage(message string)
	Build() ErrorResponse
}

type ErrorResponse interface {
	HasErrors() bool
}

func (e *ErrorResponseObject) AddMessage(message string) {
	e.Messages = append(e.Messages, message)
}

func (e *ErrorResponseObject) AddError(error string) {
	e.Error = error
}

func (e *ErrorResponseObject) AddErrorIfMessagesExists(error string) {
	if len(e.Messages) != 0 {
		e.Error = error
	}
}

func (e *ErrorResponseObject) Build() ErrorResponse {
	if e.Error == "" && len(e.Messages) == 0 {
		return nil
	}
	return e
}

func (e *ErrorResponseObject) HasErrors() bool {
	return e.Error != "" || len(e.Messages) != 0
}
