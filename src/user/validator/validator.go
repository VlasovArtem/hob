package validator

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/int-errors"
	baseValidator "github.com/VlasovArtem/hob/src/common/validator"
	userModel "github.com/VlasovArtem/hob/src/user/model"
)

type UserRequestValidatorObject struct {
	baseValidator.BaseValidator
}

func (u *UserRequestValidatorObject) Initialize(factory dependency.DependenciesProvider) any {
	return NewUserRequestValidator()
}

func NewUserRequestValidator() UserRequestValidator {
	return &UserRequestValidatorObject{baseValidator.NewBaseValidator()}
}

type UserRequestValidator interface {
	ValidateCreateRequest(request userModel.CreateUserRequest) int_errors.ErrorResponse
	ValidateUpdateRequest(request userModel.UpdateUserRequest) int_errors.ErrorResponse
}

func (u *UserRequestValidatorObject) ValidateCreateRequest(request userModel.CreateUserRequest) int_errors.ErrorResponse {
	return u.ValidateStringFieldNotEmpty(request.Email, "email should not be empty").
		ValidateStringFieldNotEmpty(request.Password, "password should not be empty").
		Result("Create User Request Validation Error")
}

func (u *UserRequestValidatorObject) ValidateUpdateRequest(request userModel.UpdateUserRequest) int_errors.ErrorResponse {
	return u.ValidateStringFieldNotEmpty(request.Password, "password should not be empty").
		Result("Create User Request Validation Error")
}
