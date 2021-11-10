package handler

import (
	helperModel "common/model"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"test"
	"test/testhelper"
	"testing"
	"user/model"
	"user/service"
)

var userService, handler = func() (service.UserService, UserHandler) {
	userService := service.NewUserService()

	return userService, NewUserHandler(userService)
}()

func Test_AddUser(t *testing.T) {
	request := test.GetCreateUserRequest()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/user").
		WithMethod("POST").
		WithHandler(handler.AddUser()).
		WithBody(request)

	content := testRequest.Verify(t, http.StatusCreated)

	actualResponse := model.UserResponse{}

	json.Unmarshal(content, &actualResponse)

	assert.Equal(t, test.GetUserResponse(actualResponse.Id, actualResponse.Email), actualResponse)
}

func Test_AddUserWithInvalidRequest(t *testing.T) {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/user").
		WithMethod("POST").
		WithHandler(handler.AddUser())

	testRequest.Verify(t, http.StatusBadRequest)
}

func Test_AddUserWithMissingDetails(t *testing.T) {
	request := test.GetCreateUserRequest()

	request.Email = ""
	request.Password = []byte{}

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/user").
		WithMethod("POST").
		WithHandler(handler.AddUser()).
		WithBody(request)

	response := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, helperModel.ErrorResponse{
		Error: "Create User Request Validation Error",
		Messages: []string{
			"email should not be empty",
			"password should not be empty",
		},
	}, testhelper.ReadErrorResponse(response))
}

func Test_AddUserWithExistingEmail(t *testing.T) {
	request := test.GetCreateUserRequest()

	request.Email = "newemail@mail.com"

	_, err := userService.AddUser(request)

	assert.Nil(t, err)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/user").
		WithMethod("POST").
		WithHandler(handler.AddUser()).
		WithBody(request)

	response := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, helperModel.ErrorResponse{
		Error: fmt.Sprintf("user with '%s' already exists", request.Email),
	}, testhelper.ReadErrorResponse(response))
}

func Test_FindById(t *testing.T) {
	request := test.GetCreateUserRequest()

	user, err := userService.AddUser(request)

	assert.Nil(t, err)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/user/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", user.Id.String())

	content := testRequest.Verify(t, http.StatusOK)

	actual := model.UserResponse{}

	err = json.Unmarshal(content, &actual)

	assert.Nil(t, err)

	assert.Equal(t, user, actual)
}

func Test_FindByIdWithMissingParameter(t *testing.T) {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/user/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById())

	content := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "parameter 'id' not found\n", string(content))
}

func Test_FindByIdInvalid(t *testing.T) {
	type args struct {
		id   string
		code int
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "invalid uuid",
			args: args{
				id:   "id",
				code: http.StatusBadRequest,
			},
		},
		{
			name: "with not exists",
			args: args{
				id:   uuid.New().String(),
				code: http.StatusNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testRequest := testhelper.NewTestRequest().
				WithURL("https://test.com/api/v1/user/{id}").
				WithMethod("GET").
				WithHandler(handler.FindById()).
				WithVar("id", tt.args.id)

			testRequest.Verify(t, tt.args.code)
		})
	}
}
