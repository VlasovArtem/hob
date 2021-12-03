package validator

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/errors"
	baseValidator "github.com/VlasovArtem/hob/src/common/validator"
	userModel "github.com/VlasovArtem/hob/src/user/model"
)

type UserRequestValidatorObject struct {
	baseValidator.BaseValidator
}

func (u *UserRequestValidatorObject) Initialize(factory dependency.DependenciesFactory) interface{} {
	return NewUserRequestValidator()
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
