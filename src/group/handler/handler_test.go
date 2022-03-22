package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	int_errors "github.com/VlasovArtem/hob/src/common/int-errors"
	"github.com/VlasovArtem/hob/src/group/mocks"
	"github.com/VlasovArtem/hob/src/group/model"
	"github.com/VlasovArtem/hob/src/test/testhelper"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var (
	groupService *mocks.GroupService
)

func handlerGenerator() GroupHandler {
	groupService = new(mocks.GroupService)

	return NewGroupHandler(groupService)
}

func Test_Add_WithNotValidRequest(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/group").
		WithMethod("POST").
		WithHandler(handler.Add())

	testRequest.Verify(t, http.StatusBadRequest)
}

func Test_Add(t *testing.T) {
	handler := handlerGenerator()

	request := mocks.GenerateCreateGroupRequest()

	groupService.On("Add", request).Return(request.ToEntity().ToDto(), nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/group").
		WithMethod("POST").
		WithBody(request).
		WithHandler(handler.Add())

	body := testRequest.Verify(t, http.StatusCreated)

	actual := model.GroupDto{}

	json.Unmarshal(body, &actual)

	assert.Equal(t,
		model.GroupDto{
			Id:      actual.Id,
			Name:    "name",
			OwnerId: actual.OwnerId,
		}, actual)
}

func Test_Add_WithErrorFromService(t *testing.T) {
	handler := handlerGenerator()

	request := mocks.GenerateCreateGroupRequest()

	groupService.On("Add", request).Return(model.GroupDto{}, errors.New("error"))

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/group").
		WithMethod("POST").
		WithBody(request).
		WithHandler(handler.Add())

	actual := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "error\n", string(actual))
}

func Test_FindById(t *testing.T) {
	handler := handlerGenerator()

	groupDto := mocks.GenerateGroupDto()

	groupService.On("FindById", groupDto.Id).Return(groupDto, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/group/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", groupDto.Id.String())

	body := testRequest.Verify(t, http.StatusOK)

	var responses model.GroupDto
	json.Unmarshal(body, &responses)

	assert.Equal(t, groupDto, responses)
}

func Test_FindById_WithErrorFromService(t *testing.T) {
	tests := []struct {
		err        error
		statusCode int
	}{
		{
			err:        errors.New("error"),
			statusCode: http.StatusBadRequest,
		},
		{
			err:        int_errors.NewErrNotFound("error %s", "test"),
			statusCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		handler := handlerGenerator()

		id := uuid.New()

		groupService.On("FindById", id).Return(model.GroupDto{}, test.err)

		testRequest := testhelper.NewTestRequest().
			WithURL("https://test.com/api/v1/group/{id}").
			WithMethod("GET").
			WithHandler(handler.FindById()).
			WithVar("id", id.String())

		body := testRequest.Verify(t, test.statusCode)

		assert.Equal(t, fmt.Sprintf("%s\n", test.err.Error()), string(body))
	}
}

func Test_FindById_WithInvalidId(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/group/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", "id")

	body := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "the id is not valid id\n", string(body))
}

func Test_FindById_WithMissingId(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/group/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById())

	body := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "parameter 'id' not found\n", string(body))
}

func Test_FindByUserId(t *testing.T) {
	handler := handlerGenerator()

	groupDto := mocks.GenerateGroupDto()

	groups := []model.GroupDto{groupDto}
	groupService.On("FindByUserId", groupDto.Id).Return(groups)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/group/user/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByUserId()).
		WithVar("id", groupDto.Id.String())

	body := testRequest.Verify(t, http.StatusOK)

	var responses []model.GroupDto
	json.Unmarshal(body, &responses)

	assert.Equal(t, groups, responses)
}

func Test_FindByUserId_WithEmptyResponse(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	groupService.On("FindByUserId", id).Return([]model.GroupDto{})

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/group/user/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByUserId()).
		WithVar("id", id.String())

	body := testRequest.Verify(t, http.StatusOK)

	var responses []model.GroupDto
	json.Unmarshal(body, &responses)

	assert.Equal(t, []model.GroupDto{}, responses)
}

func Test_FindByUserId_WithInvalidId(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/group/user/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByUserId()).
		WithVar("id", "id")

	body := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "the id is not valid id\n", string(body))
}

func Test_FindByUserId_WithMissingId(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/group/user/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByUserId())

	body := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "parameter 'id' not found\n", string(body))
}
