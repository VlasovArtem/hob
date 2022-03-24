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
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type UserHandlerTestSuite struct {
	testhelper.MockTestSuite[UserHandler]
	userService   *mocks.UserService
	userValidator *mocks.UserRequestValidator
}

func TestUserHandlerTestSuite(t *testing.T) {
	ts := &UserHandlerTestSuite{}
	ts.TestObjectGenerator = func() UserHandler {
		ts.userService = new(mocks.UserService)
		ts.userValidator = new(mocks.UserRequestValidator)

		return NewUserHandler(ts.userService, ts.userValidator)
	}

	suite.Run(t, ts)
}

func (u *UserHandlerTestSuite) Test_AddUser() {
	request := mocks.GenerateCreateUserRequest()

	u.userValidator.On("ValidateCreateRequest", request).Return(nil)
	u.userService.On("Add", request).Return(mocks.GenerateUserResponse(), nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/user").
		WithMethod("POST").
		WithHandler(u.TestO.Add()).
		WithBody(request)

	content := testRequest.Verify(u.T(), http.StatusCreated)

	actualResponse := model.UserDto{}

	json.Unmarshal(content, &actualResponse)

	assert.Equal(u.T(), model.UserDto{
		Id:        actualResponse.Id,
		FirstName: "First Name",
		LastName:  "Last Name",
		Email:     "mail@mail.com",
	}, actualResponse)
}

func (u *UserHandlerTestSuite) Test_AddUserWithInvalidRequest() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/user").
		WithMethod("POST").
		WithHandler(u.TestO.Add())

	testRequest.Verify(u.T(), http.StatusBadRequest)
}

func (u *UserHandlerTestSuite) Test_AddUserWithMissingDetails() {
	request := mocks.GenerateCreateUserRequest()

	error := helperModel.NewWithDetails("error", "details")
	u.userValidator.On("ValidateCreateRequest", request).Return(error)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/user").
		WithMethod("POST").
		WithHandler(u.TestO.Add()).
		WithBody(request)

	response := testRequest.Verify(u.T(), http.StatusBadRequest)

	assert.Equal(u.T(), *error.(*helperModel.ErrorResponseObject), testhelper.ReadErrorResponse(response))
}

func (u *UserHandlerTestSuite) Test_Add_WithErrorFromService() {
	request := mocks.GenerateCreateUserRequest()

	u.userValidator.On("ValidateCreateRequest", request).Return(nil)

	err := errors.New("error")

	u.userService.On("Add", request).Return(model.UserDto{}, err)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/user").
		WithMethod("POST").
		WithHandler(u.TestO.Add()).
		WithBody(request)

	response := testRequest.Verify(u.T(), http.StatusBadRequest)

	assert.Equal(u.T(), []byte("error\n"), response)
}

func (u *UserHandlerTestSuite) Test_FindById() {
	request := mocks.GenerateCreateUserRequest()
	expected := request.ToEntity().ToDto()

	u.userService.On("FindById", expected.Id).Return(expected, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/user/{id}").
		WithMethod("GET").
		WithHandler(u.TestO.FindById()).
		WithVar("id", expected.Id.String())

	content := testRequest.Verify(u.T(), http.StatusOK)

	actual := model.UserDto{}

	err := json.Unmarshal(content, &actual)

	assert.Nil(u.T(), err)

	assert.Equal(u.T(), expected, actual)
}

func (u *UserHandlerTestSuite) Test_FindById_WithErrorFromService() {
	u.userService.On("FindById", mock.Anything).Return(model.UserDto{}, errors.New("test"))

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/user/{id}").
		WithMethod("GET").
		WithHandler(u.TestO.FindById()).
		WithVar("id", uuid.New().String())

	content := testRequest.Verify(u.T(), http.StatusBadRequest)

	assert.Equal(u.T(), "test\n", string(content))
}

func (u *UserHandlerTestSuite) Test_FindByIdWithMissingParameter() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/user/{id}").
		WithMethod("GET").
		WithHandler(u.TestO.FindById())

	content := testRequest.Verify(u.T(), http.StatusBadRequest)

	assert.Equal(u.T(), "parameter 'id' not found\n", string(content))
}

func (u *UserHandlerTestSuite) Test_FindById_WithInvalidUUID() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/user/{id}").
		WithMethod("GET").
		WithHandler(u.TestO.FindById()).
		WithVar("id", "id")

	content := testRequest.Verify(u.T(), http.StatusBadRequest)

	assert.Equal(u.T(), "the id is not valid id\n", string(content))
}

func (u *UserHandlerTestSuite) Test_Update() {
	id, request := mocks.GenerateUpdateUserRequest()

	u.userValidator.On("ValidateUpdateRequest", request).Return(nil)
	u.userService.On("Update", id, request).Return(nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/user/{id}").
		WithMethod("PUT").
		WithHandler(u.TestO.Update()).
		WithBody(request).
		WithVar("id", id.String())

	testRequest.Verify(u.T(), http.StatusOK)
}

func (u *UserHandlerTestSuite) Test_Update_WithInvalidRequest() {
	id := uuid.New()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/user/{id}").
		WithMethod("PUT").
		WithHandler(u.TestO.Update()).
		WithVar("id", id.String())

	testRequest.Verify(u.T(), http.StatusBadRequest)

	u.userService.AssertNotCalled(u.T(), "Update", mock.Anything, mock.Anything)
}

func (u *UserHandlerTestSuite) Test_Update_WithInvalidUUID() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/user/{id}").
		WithMethod("PUT").
		WithHandler(u.TestO.Update()).
		WithVar("id", "id")

	testRequest.Verify(u.T(), http.StatusBadRequest)

	u.userService.AssertNotCalled(u.T(), "Update", mock.Anything, mock.Anything)
}

func (u *UserHandlerTestSuite) Test_Update_WithMissingIdParameter() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/user/{id}").
		WithMethod("PUT").
		WithHandler(u.TestO.Update())

	testRequest.Verify(u.T(), http.StatusBadRequest)

	u.userService.AssertNotCalled(u.T(), "Update", mock.Anything, mock.Anything)
}

func (u *UserHandlerTestSuite) Test_Update_WithMissingDetails() {
	id, request := mocks.GenerateUpdateUserRequest()

	errorResponse := helperModel.NewWithDetails("error", "details")
	u.userValidator.On("ValidateUpdateRequest", request).Return(errorResponse)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/user/{id}").
		WithMethod("PUT").
		WithHandler(u.TestO.Update()).
		WithBody(request).
		WithVar("id", id.String())

	response := testRequest.Verify(u.T(), http.StatusBadRequest)

	assert.Equal(u.T(), *errorResponse.(*helperModel.ErrorResponseObject), testhelper.ReadErrorResponse(response))
	u.userService.AssertNotCalled(u.T(), "Update", mock.Anything, mock.Anything)
}

func (u *UserHandlerTestSuite) Test_Update_WithErrorFromService() {
	id, request := mocks.GenerateUpdateUserRequest()

	u.userValidator.On("ValidateUpdateRequest", request).Return(nil)
	u.userService.On("Update", id, request).Return(errors.New("error"))

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/user/{id}").
		WithMethod("PUT").
		WithHandler(u.TestO.Update()).
		WithBody(request).
		WithVar("id", id.String())

	response := testRequest.Verify(u.T(), http.StatusBadRequest)

	assert.Equal(u.T(), []byte("error\n"), response)
}
