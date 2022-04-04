package validator

import (
	"github.com/VlasovArtem/hob/src/common/int-errors"
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

	result := validator.ValidateCreateRequest(createUserRequest).(*int_errors.ErrResponse)

	assert.NotNil(t, result)
	response := result.Response.(*int_errors.ErrorResponseObject)
	assert.Equal(t, "Create User Request Validation Error", response.Message)
	assert.Equal(t, []string{"email should not be empty"}, response.Details)
}

func Test_WithCreateUserRequest_WithEmptyPassword(t *testing.T) {
	validator := NewUserRequestValidator()

	createUserRequest := mocks.GenerateCreateUserRequest()
	createUserRequest.Password = ""

	result := validator.ValidateCreateRequest(createUserRequest).(*int_errors.ErrResponse)

	assert.NotNil(t, result)
	response := result.Response.(*int_errors.ErrorResponseObject)

	assert.Equal(t, "Create User Request Validation Error", response.Message)
	assert.Equal(t, []string{"password should not be empty"}, response.Details)
}

func Test_WithCreateUserRequest_WithAllErrors(t *testing.T) {
	validator := NewUserRequestValidator()

	createUserRequest := mocks.GenerateCreateUserRequest()
	createUserRequest.Email = ""
	createUserRequest.Password = ""

	result := validator.ValidateCreateRequest(createUserRequest).(*int_errors.ErrResponse)

	assert.NotNil(t, result)
	response := result.Response.(*int_errors.ErrorResponseObject)

	assert.Equal(t, "Create User Request Validation Error", response.Message)
	assert.ElementsMatch(t, []string{"email should not be empty", "password should not be empty"}, response.Details)
}
