package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/common"
	int_errors "github.com/VlasovArtem/hob/src/common/int-errors"
	groupModel "github.com/VlasovArtem/hob/src/group/model"
	"github.com/VlasovArtem/hob/src/house/mocks"
	"github.com/VlasovArtem/hob/src/house/model"
	"github.com/VlasovArtem/hob/src/test"
	"github.com/VlasovArtem/hob/src/test/testhelper"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type HouseHandlerTestSuite struct {
	testhelper.MockTestSuite[HouseHandler]
	houses *mocks.HouseService
}

func TestHouseHandlerTestSuite(t *testing.T) {
	testingSuite := &HouseHandlerTestSuite{}
	testingSuite.TestObjectGenerator = func() HouseHandler {
		testingSuite.houses = new(mocks.HouseService)
		return NewHouseHandler(testingSuite.houses)
	}

	suite.Run(t, testingSuite)
}

func (h *HouseHandlerTestSuite) Test_Add_WithNotValidRequest() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/house").
		WithMethod("POST").
		WithHandler(h.TestO.Add())

	testRequest.Verify(h.T(), http.StatusBadRequest)
}

func (h *HouseHandlerTestSuite) Test_Add() {
	request := mocks.GenerateCreateHouseRequest()

	h.houses.On("Add", request).Return(request.ToEntity(test.CountryObject).ToDto(), nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/house").
		WithMethod("POST").
		WithBody(request).
		WithHandler(h.TestO.Add())

	body := testRequest.Verify(h.T(), http.StatusCreated)

	actual := model.HouseDto{}

	json.Unmarshal(body, &actual)

	assert.Equal(h.T(),
		model.HouseDto{
			Id:          actual.Id,
			Name:        "Test House",
			CountryCode: "UA",
			City:        "City",
			StreetLine1: "StreetLine1",
			StreetLine2: "StreetLine2",
			UserId:      actual.UserId,
			Groups:      []groupModel.GroupDto{},
		}, actual)
}

func (h *HouseHandlerTestSuite) Test_Add_WithErrorFromService() {
	request := mocks.GenerateCreateHouseRequest()

	h.houses.On("Add", request).Return(model.HouseDto{}, errors.New("error"))

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/house").
		WithMethod("POST").
		WithBody(request).
		WithHandler(h.TestO.Add())

	actual := testRequest.Verify(h.T(), http.StatusBadRequest)

	assert.Equal(h.T(), "error\n", string(actual))
}

func (h *HouseHandlerTestSuite) Test_AddBatch() {
	request := mocks.GenerateCreateHouseBatchRequest(2)

	serviceResponse := common.MapSlice(request.Houses, func(house model.CreateHouseRequest) model.HouseDto {
		return house.ToEntity(test.CountryObject).ToDto()
	})

	h.houses.On("AddBatch", request).Return(serviceResponse, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/house/batch").
		WithMethod("POST").
		WithBody(request).
		WithHandler(h.TestO.AddBatch())

	body := testRequest.Verify(h.T(), http.StatusCreated)

	var actual []model.HouseDto

	err := json.Unmarshal(body, &actual)
	assert.Nil(h.T(), err)

	assert.Equal(h.T(),
		[]model.HouseDto{
			{
				Id:          serviceResponse[0].Id,
				Name:        "House Name #0",
				CountryCode: "UA",
				City:        "City",
				StreetLine1: "StreetLine1",
				StreetLine2: "StreetLine2",
				UserId:      request.Houses[0].UserId,
				Groups:      []groupModel.GroupDto{},
			},
			{
				Id:          serviceResponse[1].Id,
				Name:        "House Name #1",
				CountryCode: "UA",
				City:        "City",
				StreetLine1: "StreetLine1",
				StreetLine2: "StreetLine2",
				UserId:      request.Houses[1].UserId,
				Groups:      []groupModel.GroupDto{},
			},
		}, actual)
}

func (h *HouseHandlerTestSuite) Test_AddBatch_WithErrorFromService() {
	request := mocks.GenerateCreateHouseBatchRequest(1)

	errorResponse := int_errors.NewBuilder().
		WithDetail("message").
		WithMessage("error")
	err := int_errors.NewErrResponse(errorResponse)

	h.houses.On("AddBatch", request).Return([]model.HouseDto{}, err)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/house/batch").
		WithMethod("POST").
		WithBody(request).
		WithHandler(h.TestO.AddBatch())

	actual := testRequest.Verify(h.T(), http.StatusBadRequest)

	expected, err := json.Marshal(errorResponse.Build())

	assert.Nil(h.T(), err)
	assert.Equal(h.T(), append(expected, []byte("\n")...), actual)
}

func (h *HouseHandlerTestSuite) Test_FindById() {
	houseResponse := mocks.GenerateHouseResponse()

	h.houses.On("FindById", houseResponse.Id).Return(houseResponse, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/house/{id}").
		WithMethod("GET").
		WithHandler(h.TestO.FindById()).
		WithVar("id", houseResponse.Id.String())

	body := testRequest.Verify(h.T(), http.StatusOK)

	var responses model.HouseDto
	json.Unmarshal(body, &responses)

	assert.Equal(h.T(), houseResponse, responses)
}

func (h *HouseHandlerTestSuite) Test_FindById_WithErrorFromService() {
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
		id := uuid.New()

		h.houses.On("FindById", id).Return(model.HouseDto{}, test.err)

		testRequest := testhelper.NewTestRequest().
			WithURL("https://test.com/api/v1/house/{id}").
			WithMethod("GET").
			WithHandler(h.TestO.FindById()).
			WithVar("id", id.String())

		body := testRequest.Verify(h.T(), test.statusCode)

		assert.Equal(h.T(), fmt.Sprintf("%s\n", test.err.Error()), string(body))
	}
}

func (h *HouseHandlerTestSuite) Test_FindById_WithInvalidId() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/house/{id}").
		WithMethod("GET").
		WithHandler(h.TestO.FindById()).
		WithVar("id", "id")

	body := testRequest.Verify(h.T(), http.StatusBadRequest)

	assert.Equal(h.T(), "the id is not valid id\n", string(body))
}

func (h *HouseHandlerTestSuite) Test_FindById_WithMissingId() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/house/{id}").
		WithMethod("GET").
		WithHandler(h.TestO.FindById())

	body := testRequest.Verify(h.T(), http.StatusBadRequest)

	assert.Equal(h.T(), "parameter 'id' not found\n", string(body))
}

func (h *HouseHandlerTestSuite) Test_FindByUserId() {
	houseResponse := mocks.GenerateHouseResponse()

	houseResponses := []model.HouseDto{houseResponse}
	h.houses.On("FindByUserId", houseResponse.Id).Return(houseResponses)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/house/user/{id}").
		WithMethod("GET").
		WithHandler(h.TestO.FindByUserId()).
		WithVar("id", houseResponse.Id.String())

	body := testRequest.Verify(h.T(), http.StatusOK)

	var responses []model.HouseDto
	json.Unmarshal(body, &responses)

	assert.Equal(h.T(), houseResponses, responses)
}

func (h *HouseHandlerTestSuite) Test_FindByUserId_WithEmptyResponse() {
	id := uuid.New()

	h.houses.On("FindByUserId", id).Return([]model.HouseDto{})

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/house/user/{id}").
		WithMethod("GET").
		WithHandler(h.TestO.FindByUserId()).
		WithVar("id", id.String())

	body := testRequest.Verify(h.T(), http.StatusOK)

	var responses []model.HouseDto
	json.Unmarshal(body, &responses)

	assert.Equal(h.T(), []model.HouseDto{}, responses)
}

func (h *HouseHandlerTestSuite) Test_FindByUserId_WithInvalidId() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/house/user/{id}").
		WithMethod("GET").
		WithHandler(h.TestO.FindByUserId()).
		WithVar("id", "id")

	body := testRequest.Verify(h.T(), http.StatusBadRequest)

	assert.Equal(h.T(), "the id is not valid id\n", string(body))
}

func (h *HouseHandlerTestSuite) Test_FindByUserId_WithMissingId() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/house/user/{id}").
		WithMethod("GET").
		WithHandler(h.TestO.FindByUserId())

	body := testRequest.Verify(h.T(), http.StatusBadRequest)

	assert.Equal(h.T(), "parameter 'id' not found\n", string(body))
}
