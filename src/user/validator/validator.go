package validator

import (
	"common/dependency"
	"common/errors"
	baseValidator "common/validator"
	userModel "user/model"
)

type UserRequestValidatorObject struct {
	baseValidator.BaseValidator
}

func (u *UserRequestValidatorObject) Initialize(factory dependency.DependenciesFactory) {
	factory.Add(NewUserRequestValidator())
}

func NewUserRequestValidator() UserRequestValidator {
	return &UserRequestValidatorObject{baseValidator.NewBaseValidator()}
}

type UserRequestValidator interface {
	ValidateCreateRequest(request userModel.CreateUserRequest) errors.ErrorResponse
}

func (u *UserRequestValidatorObject) ValidateCreateRequest(request userModel.CreateUserRequest) errors.ErrorResponse {
	return u.ValidateStringFieldNotEmpty(request.Email, "email should not be empty").
		ValidateStringFieldNotEmpty(request.Password, "password should not be empty").
		Result("Create User Request Validation Error")
}
