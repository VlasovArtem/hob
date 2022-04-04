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
	Result(error string) error
}

func (v *Validator) ValidateStringFieldNotEmpty(value string, message string) BaseValidator {
	if value == "" {
		v.errors.WithDetail(message)
	}
	return v
}

func (v *Validator) Result(error string) error {
	if v.errors.HasErrors() {
		v.errors.WithMessage(error)
		return int_errors.NewErrResponse(v.errors)
	}

	return nil
}
