package validator

import (
	"github.com/VlasovArtem/hob/src/common/errors"
	"github.com/VlasovArtem/hob/src/user/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_WithCreateUserRequest(t *testing.T) {
	validator := NewUserRequestValidator()

	createUserRequest := mocks.GenerateCreateUserRequest()

	assert.Nil(t, validator.ValidateCreateRequest(createUserRequest))
}

func Test_WithCreateUserRequest_WithEmptyEmail(t *testing.T) {
	validator := NewUserRequestValidator()

	createUserRequest := mocks.GenerateCreateUserRequest()
	createUserRequest.Email = ""

	result := validator.ValidateCreateRequest(createUserRequest).(*errors.ErrorResponseObject)

	assert.NotNil(t, result)
	assert.Equal(t, "Create User Request Validation Error", result.Error)
	assert.Equal(t, []string{"email should not be empty"}, result.Messages)
}

func Test_WithCreateUserRequest_WithEmptyPassword(t *testing.T) {
	validator := NewUserRequestValidator()

	createUserRequest := mocks.GenerateCreateUserRequest()
	createUserRequest.Password = ""

	result := validator.ValidateCreateRequest(createUserRequest).(*errors.ErrorResponseObject)

	assert.NotNil(t, result)
	assert.Equal(t, "Create User Request Validation Error", result.Error)
	assert.Equal(t, []string{"password should not be empty"}, result.Messages)
}

func Test_WithCreateUserRequest_WithAllErrors(t *testing.T) {
	validator := NewUserRequestValidator()

	createUserRequest := mocks.GenerateCreateUserRequest()
	createUserRequest.Email = ""
	createUserRequest.Password = ""

	result := validator.ValidateCreateRequest(createUserRequest).(*errors.ErrorResponseObject)

	assert.NotNil(t, result)
	assert.Equal(t, "Create User Request Validation Error", result.Error)
	assert.ElementsMatch(t, []string{"email should not be empty", "password should not be empty"}, result.Messages)
}
