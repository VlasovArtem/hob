package handler

import (
	"encoding/json"
	"errors"
	helperModel "github.com/VlasovArtem/hob/src/common/int-errors"
	"github.com/VlasovArtem/hob/src/test/testhelper"
	"github.com/VlasovArtem/hob/src/user/mocks"
	"github.com/VlasovArtem/hob/src/user/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
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

	actualResponse := model.UserDto{}

	json.Unmarshal(content, &actualResponse)

	assert.Equal(t, model.UserDto{
		Id:        actualResponse.Id,
		FirstName: "First Name",
		LastName:  "Last Name",
		Email:     "mail@mail.com",
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

	userService.On("Add", request).Return(model.UserDto{}, err)

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
	expected := request.ToEntity().ToDto()

	userService.On("FindById", expected.Id).Return(expected, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/user/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", expected.Id.String())

	content := testRequest.Verify(t, http.StatusOK)

	actual := model.UserDto{}

	err := json.Unmarshal(content, &actual)

	assert.Nil(t, err)

	assert.Equal(t, expected, actual)
}

func Test_FindById_WithErrorFromService(t *testing.T) {
	handler := generateHandler()

	userService.On("FindById", mock.Anything).Return(model.UserDto{}, errors.New("test"))

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/user/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", uuid.New().String())

	content := testRequest.Verify(t, http.StatusBadRequest)

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

func Test_Update(t *testing.T) {
	handler := generateHandler()

	request := mocks.GenerateUpdateUserRequest()
	id := uuid.New()

	userValidator.On("ValidateUpdateRequest", request).Return(nil)
	userService.On("Update", id, request).Return(nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/user/{id}").
		WithMethod("PUT").
		WithHandler(handler.Update()).
		WithBody(request).
		WithVar("id", id.String())

	testRequest.Verify(t, http.StatusOK)
}

func Test_Update_WithInvalidRequest(t *testing.T) {
	handler := generateHandler()

	id := uuid.New()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/user/{id}").
		WithMethod("PUT").
		WithHandler(handler.Update()).
		WithVar("id", id.String())

	testRequest.Verify(t, http.StatusBadRequest)

	userService.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func Test_Update_WithInvalidUUID(t *testing.T) {
	handler := generateHandler()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/user/{id}").
		WithMethod("PUT").
		WithHandler(handler.Update()).
		WithVar("id", "id")

	testRequest.Verify(t, http.StatusBadRequest)

	userService.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func Test_Update_WithMissingIdParameter(t *testing.T) {
	handler := generateHandler()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/user/{id}").
		WithMethod("PUT").
		WithHandler(handler.Update())

	testRequest.Verify(t, http.StatusBadRequest)

	userService.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func Test_Update_WithMissingDetails(t *testing.T) {
	handler := generateHandler()

	request := mocks.GenerateUpdateUserRequest()

	errorResponse := helperModel.NewWithDetails("error", "details")
	userValidator.On("ValidateUpdateRequest", request).Return(errorResponse)
	id := uuid.New()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/user/{id}").
		WithMethod("PUT").
		WithHandler(handler.Update()).
		WithBody(request).
		WithVar("id", id.String())

	response := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, *errorResponse.(*helperModel.ErrorResponseObject), testhelper.ReadErrorResponse(response))
	userService.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func Test_Update_WithErrorFromService(t *testing.T) {
	handler := generateHandler()

	request := mocks.GenerateUpdateUserRequest()
	id := uuid.New()

	userValidator.On("ValidateUpdateRequest", request).Return(nil)
	userService.On("Update", id, request).Return(errors.New("error"))

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/user/{id}").
		WithMethod("PUT").
		WithHandler(handler.Update()).
		WithBody(request).
		WithVar("id", id.String())

	response := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, helperModel.ErrorResponseObject{Error: "error"}, testhelper.ReadErrorResponse(response))
}
