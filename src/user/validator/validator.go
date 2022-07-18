package validator

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	baseValidator "github.com/VlasovArtem/hob/src/common/validator"
	userModel "github.com/VlasovArtem/hob/src/user/model"
)

type UserRequestValidatorStr struct {
	baseValidator.BaseValidator
}

func (u *UserRequestValidatorStr) GetRequiredDependencies() []dependency.Requirements {
	return []dependency.Requirements{}
}

func (u *UserRequestValidatorStr) Initialize(factory dependency.DependenciesProvider) any {
	return NewUserRequestValidator()
}

func NewUserRequestValidator() UserRequestValidator {
	return &UserRequestValidatorStr{baseValidator.NewBaseValidator()}
}

type UserRequestValidator interface {
	ValidateCreateRequest(request userModel.CreateUserRequest) error
	ValidateUpdateRequest(request userModel.UpdateUserRequest) error
}

func (u *UserRequestValidatorStr) ValidateCreateRequest(request userModel.CreateUserRequest) error {
	return u.ValidateStringFieldNotEmpty(request.Email, "email should not be empty").
		ValidateStringFieldNotEmpty(request.Password, "password should not be empty").
		Result("Create User Request Validation Error")
}

func (u *UserRequestValidatorStr) ValidateUpdateRequest(request userModel.UpdateUserRequest) error {
	return u.ValidateStringFieldNotEmpty(request.Password, "password should not be empty").
		Result("Create User Request Validation Error")
}
