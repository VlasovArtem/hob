package int_errors

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"reflect"
)

var errNotFoundType = reflect.TypeOf(ErrNotFound{})

var errResponseType = reflect.TypeOf(ErrResponse{})

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

type ErrResponse struct {
	Response ErrorResponse
}

func NewErrResponse(build ErrorResponseBuilder) error {
	return &ErrResponse{build.Build()}
}

func (e ErrResponse) Error() string {
	if marshal, err := json.Marshal(e.Response); err != nil {
		log.Err(err).Msg("")
		return err.Error()
	} else {
		return string(marshal)
	}
}

func (e ErrResponse) Is(err error) bool {
	return reflect.TypeOf(err) == errResponseType
}

type ErrorResponseObject struct {
	Message string
	Details []string
}

func NewBuilder() ErrorResponseBuilder {
	return &ErrorResponseObject{}
}

func New() ErrorResponse {
	return &ErrorResponseObject{}
}

func NewWithDetails(error string, messages ...string) ErrorResponse {
	return &ErrorResponseObject{
		Message: error,
		Details: messages,
	}
}

type ErrorResponseBuilder interface {
	ErrorResponse
	WithMessage(error string) ErrorResponseBuilder
	WithDetail(message string) ErrorResponseBuilder
	Build() ErrorResponse
}

type ErrorResponse interface {
	HasErrors() bool
}

func (e *ErrorResponseObject) WithDetail(detail string) ErrorResponseBuilder {
	e.Details = append(e.Details, detail)

	return e
}

func (e *ErrorResponseObject) WithMessage(message string) ErrorResponseBuilder {
	e.Message = message

	return e
}

func (e *ErrorResponseObject) Build() ErrorResponse {
	if e.Message == "" && len(e.Details) == 0 {
		return nil
	}
	return e
}

func (e *ErrorResponseObject) BuildWithMessage(message string) ErrorResponse {
	e.Message = message
	if e.Message == "" && len(e.Details) == 0 {
		return nil
	}
	return e
}

func (e *ErrorResponseObject) HasErrors() bool {
	return e.Message != "" || len(e.Details) != 0
}
