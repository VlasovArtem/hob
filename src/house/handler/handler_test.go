package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	int_errors "github.com/VlasovArtem/hob/src/common/int-errors"
	"github.com/VlasovArtem/hob/src/house/mocks"
	"github.com/VlasovArtem/hob/src/house/model"
	"github.com/VlasovArtem/hob/src/test"
	"github.com/VlasovArtem/hob/src/test/testhelper"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var (
	houses *mocks.HouseService
)

func handlerGenerator() HouseHandler {
	houses = new(mocks.HouseService)

	return NewHouseHandler(houses)
}

func Test_Add_WithNotValidRequest(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/house").
		WithMethod("POST").
		WithHandler(handler.Add())

	testRequest.Verify(t, http.StatusBadRequest)
}

func Test_Add(t *testing.T) {
	handler := handlerGenerator()

	request := mocks.GenerateCreateHouseRequest()

	houses.On("Add", request).Return(request.ToEntity(test.CountryObject).ToDto(), nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/house").
		WithMethod("POST").
		WithBody(request).
		WithHandler(handler.Add())

	body := testRequest.Verify(t, http.StatusCreated)

	actual := model.HouseDto{}

	json.Unmarshal(body, &actual)

	assert.Equal(t,
		model.HouseDto{
			Id:          actual.Id,
			Name:        "Test House",
			CountryCode: "UA",
			City:        "City",
			StreetLine1: "StreetLine1",
			StreetLine2: "StreetLine2",
			UserId:      actual.UserId,
		}, actual)
}

func Test_Add_WithErrorFromService(t *testing.T) {
	handler := handlerGenerator()

	request := mocks.GenerateCreateHouseRequest()

	houses.On("Add", request).Return(model.HouseDto{}, errors.New("error"))

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/house").
		WithMethod("POST").
		WithBody(request).
		WithHandler(handler.Add())

	actual := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "error\n", string(actual))
}

func Test_FindById(t *testing.T) {
	handler := handlerGenerator()

	houseResponse := mocks.GenerateHouseResponse()

	houses.On("FindById", houseResponse.Id).Return(houseResponse, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/house/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", houseResponse.Id.String())

	body := testRequest.Verify(t, http.StatusOK)

	var responses model.HouseDto
	json.Unmarshal(body, &responses)

	assert.Equal(t, houseResponse, responses)
}

func Test(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: test cases
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

		})
	}
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

		houses.On("FindById", id).Return(model.HouseDto{}, test.err)

		testRequest := testhelper.NewTestRequest().
			WithURL("https://test.com/api/v1/house/{id}").
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
		WithURL("https://test.com/api/v1/house/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", "id")

	body := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "the id is not valid id\n", string(body))
}

func Test_FindById_WithMissingId(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/house/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById())

	body := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "parameter 'id' not found\n", string(body))
}

func Test_FindByUserId(t *testing.T) {
	handler := handlerGenerator()

	houseResponse := mocks.GenerateHouseResponse()

	houseResponses := []model.HouseDto{houseResponse}
	houses.On("FindByUserId", houseResponse.Id).Return(houseResponses)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/house/user/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByUserId()).
		WithVar("id", houseResponse.Id.String())

	body := testRequest.Verify(t, http.StatusOK)

	var responses []model.HouseDto
	json.Unmarshal(body, &responses)

	assert.Equal(t, houseResponses, responses)
}

func Test_FindByUserId_WithEmptyResponse(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	houses.On("FindByUserId", id).Return([]model.HouseDto{})

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/house/user/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByUserId()).
		WithVar("id", id.String())

	body := testRequest.Verify(t, http.StatusOK)

	var responses []model.HouseDto
	json.Unmarshal(body, &responses)

	assert.Equal(t, []model.HouseDto{}, responses)
}

func Test_FindByUserId_WithInvalidId(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/house/user/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByUserId()).
		WithVar("id", "id")

	body := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "the id is not valid id\n", string(body))
}

func Test_FindByUserId_WithMissingId(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/house/user/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByUserId())

	body := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "parameter 'id' not found\n", string(body))
}
