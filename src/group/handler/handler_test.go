package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/common"
	interrors "github.com/VlasovArtem/hob/src/common/int-errors"
	"github.com/VlasovArtem/hob/src/group/mocks"
	"github.com/VlasovArtem/hob/src/group/model"
	"github.com/VlasovArtem/hob/src/test/testhelper"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type GroupHandlerTestSuite struct {
	testhelper.MockTestSuite[GroupHandler]
	groupService *mocks.GroupService
}

func TestGroupHandlerTestSuite(t *testing.T) {
	testingSuite := &GroupHandlerTestSuite{}
	testingSuite.TestObjectGenerator = func() GroupHandler {
		testingSuite.groupService = new(mocks.GroupService)
		return NewGroupHandler(testingSuite.groupService)
	}

	suite.Run(t, testingSuite)
}

func (g *GroupHandlerTestSuite) Test_Add_WithNotValidRequest() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/group").
		WithMethod("POST").
		WithHandler(g.TestO.Add())

	testRequest.Verify(g.T(), http.StatusBadRequest)
}

func (g *GroupHandlerTestSuite) Test_Add() {
	request := mocks.GenerateCreateGroupRequest()

	g.groupService.On("Add", request).Return(request.ToEntity().ToDto(), nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/group").
		WithMethod("POST").
		WithBody(request).
		WithHandler(g.TestO.Add())

	body := testRequest.Verify(g.T(), http.StatusCreated)

	actual := model.GroupDto{}

	json.Unmarshal(body, &actual)

	assert.Equal(g.T(),
		model.GroupDto{
			Id:      actual.Id,
			Name:    "name",
			OwnerId: actual.OwnerId,
		}, actual)
}

func (g *GroupHandlerTestSuite) Test_Add_WithErrorFromService() {
	request := mocks.GenerateCreateGroupRequest()

	g.groupService.On("Add", request).Return(model.GroupDto{}, errors.New("error"))

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/group").
		WithMethod("POST").
		WithBody(request).
		WithHandler(g.TestO.Add())

	actual := testRequest.Verify(g.T(), http.StatusBadRequest)

	assert.Equal(g.T(), "error\n", string(actual))
}

func (g *GroupHandlerTestSuite) Test_AddBatch() {
	request := mocks.GenerateCreateGroupBatchRequest(2)

	addBatchResponse := common.MapSlice(request.Groups, func(r model.CreateGroupRequest) model.GroupDto {
		return r.ToEntity().ToDto()
	})

	g.groupService.On("AddBatch", request).Return(addBatchResponse, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/group/batch").
		WithMethod("POST").
		WithBody(request).
		WithHandler(g.TestO.AddBatch())

	body := testRequest.Verify(g.T(), http.StatusCreated)

	var actual []model.GroupDto

	err := json.Unmarshal(body, &actual)
	assert.Nil(g.T(), err)
	assert.Equal(g.T(), addBatchResponse, actual)
}

func (g *GroupHandlerTestSuite) Test_AddBatch_WithErrorFromService() {
	request := mocks.GenerateCreateGroupBatchRequest(1)

	errorResponse := interrors.NewBuilder().
		WithDetail("message").
		WithMessage("error")
	err := interrors.NewErrResponse(errorResponse)

	g.groupService.On("AddBatch", mock.Anything).Return([]model.GroupDto{}, err)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/group/batch").
		WithMethod("POST").
		WithBody(request).
		WithHandler(g.TestO.AddBatch())

	actual := testRequest.Verify(g.T(), http.StatusBadRequest)

	expected, err := json.Marshal(errorResponse.Build())

	assert.Nil(g.T(), err)
	assert.Equal(g.T(), append(expected, []byte("\n")...), actual)
}

func (g *GroupHandlerTestSuite) Test_FindById() {
	groupDto := mocks.GenerateGroupDto()

	g.groupService.On("FindById", groupDto.Id).Return(groupDto, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/group/{id}").
		WithMethod("GET").
		WithHandler(g.TestO.FindById()).
		WithVar("id", groupDto.Id.String())

	body := testRequest.Verify(g.T(), http.StatusOK)

	var responses model.GroupDto
	json.Unmarshal(body, &responses)

	assert.Equal(g.T(), groupDto, responses)
}

func (g *GroupHandlerTestSuite) Test_FindById_WithErrorFromService() {
	tests := []struct {
		err        error
		statusCode int
	}{
		{
			err:        errors.New("error"),
			statusCode: http.StatusBadRequest,
		},
		{
			err:        interrors.NewErrNotFound("error %s", "test"),
			statusCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		g.TestO = g.TestObjectGenerator()

		id := uuid.New()

		g.groupService.On("FindById", id).Return(model.GroupDto{}, test.err)

		testRequest := testhelper.NewTestRequest().
			WithURL("https://test.com/api/v1/group/{id}").
			WithMethod("GET").
			WithHandler(g.TestO.FindById()).
			WithVar("id", id.String())

		body := testRequest.Verify(g.T(), test.statusCode)

		assert.Equal(g.T(), fmt.Sprintf("%s\n", test.err.Error()), string(body))
	}
}

func (g *GroupHandlerTestSuite) Test_FindById_WithInvalidId() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/group/{id}").
		WithMethod("GET").
		WithHandler(g.TestO.FindById()).
		WithVar("id", "id")

	body := testRequest.Verify(g.T(), http.StatusBadRequest)

	assert.Equal(g.T(), "the id is not valid id\n", string(body))
}

func (g *GroupHandlerTestSuite) Test_FindById_WithMissingId() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/group/{id}").
		WithMethod("GET").
		WithHandler(g.TestO.FindById())

	body := testRequest.Verify(g.T(), http.StatusBadRequest)

	assert.Equal(g.T(), "parameter 'id' not found\n", string(body))
}

func (g *GroupHandlerTestSuite) Test_FindByUserId() {
	groupDto := mocks.GenerateGroupDto()

	groups := []model.GroupDto{groupDto}
	g.groupService.On("FindByUserId", groupDto.Id).Return(groups)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/group/user/{id}").
		WithMethod("GET").
		WithHandler(g.TestO.FindByUserId()).
		WithVar("id", groupDto.Id.String())

	body := testRequest.Verify(g.T(), http.StatusOK)

	var responses []model.GroupDto
	json.Unmarshal(body, &responses)

	assert.Equal(g.T(), groups, responses)
}

func (g *GroupHandlerTestSuite) Test_FindByUserId_WithEmptyResponse() {
	id := uuid.New()

	g.groupService.On("FindByUserId", id).Return([]model.GroupDto{})

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/group/user/{id}").
		WithMethod("GET").
		WithHandler(g.TestO.FindByUserId()).
		WithVar("id", id.String())

	body := testRequest.Verify(g.T(), http.StatusOK)

	var responses []model.GroupDto
	json.Unmarshal(body, &responses)

	assert.Equal(g.T(), []model.GroupDto{}, responses)
}

func (g *GroupHandlerTestSuite) Test_FindByUserId_WithInvalidId() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/group/user/{id}").
		WithMethod("GET").
		WithHandler(g.TestO.FindByUserId()).
		WithVar("id", "id")

	body := testRequest.Verify(g.T(), http.StatusBadRequest)

	assert.Equal(g.T(), "the id is not valid id\n", string(body))
}

func (g *GroupHandlerTestSuite) Test_FindByUserId_WithMissingId() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/group/user/{id}").
		WithMethod("GET").
		WithHandler(g.TestO.FindByUserId())

	body := testRequest.Verify(g.T(), http.StatusBadRequest)

	assert.Equal(g.T(), "parameter 'id' not found\n", string(body))
}
