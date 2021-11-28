package validator

import "common/errors"

type Validator struct {
	errors errors.ErrorResponse
}

func NewBaseValidator() BaseValidator {
	return &Validator{errors.New()}
}

type BaseValidator interface {
	ValidateStringFieldNotEmpty(value string, message string) BaseValidator
	Result(error string) errors.ErrorResponse
}

func (v *Validator) ValidateStringFieldNotEmpty(value string, message string) BaseValidator {
	if value == "" {
		v.errors.AddMessage(message)
	}
	return v
}

func (v *Validator) Result(error string) errors.ErrorResponse {
	v.errors.AddErrorIfMessagesExists(error)

	return v.errors.Result()
}
