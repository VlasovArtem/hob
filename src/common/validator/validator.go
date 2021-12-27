package validator

import "github.com/VlasovArtem/hob/src/common/int-errors"

type Validator struct {
	errors int_errors.ErrorResponseBuilder
}

func NewBaseValidator() BaseValidator {
	return &Validator{int_errors.NewBuilder()}
}

type BaseValidator interface {
	ValidateStringFieldNotEmpty(value string, message string) BaseValidator
	Result(error string) int_errors.ErrorResponse
}

func (v *Validator) ValidateStringFieldNotEmpty(value string, message string) BaseValidator {
	if value == "" {
		v.errors.AddMessage(message)
	}
	return v
}

func (v *Validator) Result(error string) int_errors.ErrorResponse {
	v.errors.AddErrorIfMessagesExists(error)

	return v.errors.Build()
}
