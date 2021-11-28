package handler

import (
	helperModel "common/errors"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"test/testhelper"
	"testing"
	"user/mocks"
	"user/model"
)

var (
	userService   *mocks.UserService
	userValidator *mocks.UserRequestValidator
)

func generateHandler() UserHandler {
	userService = new(mocks.UserService)
	userValidator = new(mocks.UserRequestValidator)

	return NewUserHandler(userService, userValidator)
}

func Test_AddUser(t *testing.T) {
	handler := generateHandler()

	request := mocks.GenerateCreateUserRequest()

	userValidator.On("ValidateCreateRequest", request).Return(nil)
	userService.On("Add", request).Return(mocks.GenerateUserResponse(), nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/user").
		WithMethod("POST").
		WithHandler(handler.Add()).
		WithBody(request)

	content := testRequest.Verify(t, http.StatusCreated)

	actualResponse := model.UserResponse{}

	json.Unmarshal(content, &actualResponse)

	assert.Equal(t, model.UserResponse{
		Id:        actualResponse.Id,
		FirstName: "First Name",
		LastName:  "Last Name",
		Email:     "mail@mai.com",
	}, actualResponse)
}

func Test_AddUserWithInvalidRequest(t *testing.T) {
	handler := generateHandler()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/user").
		WithMethod("POST").
		WithHandler(handler.Add())

	testRequest.Verify(t, http.StatusBadRequest)
}

func Test_AddUserWithMissingDetails(t *testing.T) {
	handler := generateHandler()

	request := mocks.GenerateCreateUserRequest()

	error := helperModel.NewWithDetails("error", "details")
	userValidator.On("ValidateCreateRequest", request).Return(error)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/user").
		WithMethod("POST").
		WithHandler(handler.Add()).
		WithBody(request)

	response := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, *error.(*helperModel.ErrorResponseObject), testhelper.ReadErrorResponse(response))
}

func Test_Add_WithErrorFromService(t *testing.T) {
	handler := generateHandler()

	request := mocks.GenerateCreateUserRequest()

	userValidator.On("ValidateCreateRequest", request).Return(nil)

	err := errors.New("error")

	userService.On("Add", request).Return(model.UserResponse{}, err)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/user").
		WithMethod("POST").
		WithHandler(handler.Add()).
		WithBody(request)

	response := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, helperModel.ErrorResponseObject{Error: "error"}, testhelper.ReadErrorResponse(response))
}

func Test_FindById(t *testing.T) {
	handler := generateHandler()

	request := mocks.GenerateCreateUserRequest()
	expected := request.ToEntity().ToResponse()

	userService.On("FindById", expected.Id).Return(expected, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/user/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", expected.Id.String())

	content := testRequest.Verify(t, http.StatusOK)

	actual := model.UserResponse{}

	err := json.Unmarshal(content, &actual)

	assert.Nil(t, err)

	assert.Equal(t, expected, actual)
}

func Test_FindById_WithErrorFromService(t *testing.T) {
	handler := generateHandler()

	userService.On("FindById", mock.Anything).Return(model.UserResponse{}, errors.New("test"))

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/user/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", uuid.New().String())

	content := testRequest.Verify(t, http.StatusNotFound)

	assert.Equal(t, "test\n", string(content))
}

func Test_FindByIdWithMissingParameter(t *testing.T) {
	handler := generateHandler()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/user/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById())

	content := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "parameter 'id' not found\n", string(content))
}

func Test_FindById_WithInvalidUUID(t *testing.T) {
	handler := generateHandler()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/user/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", "id")

	content := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "invalid UUID length: 2\n", string(content))
}
