package handler

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"house/model"
	"net/http"
	"test"
	"test/mock"
	"test/testhelper"
	"testing"
)

var (
	houses *mock.HouseServiceMock
	userId = testhelper.ParseUUID("98c8cab4-18be-42cf-83e9-e6369dbb2689")
)

func handlerGenerator() HouseHandler {
	houses = new(mock.HouseServiceMock)

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

	request := generateCreateHouseRequest()

	houses.On("Add", request).Return(request.ToEntity(test.CountryObject).ToResponse(), nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/house").
		WithMethod("POST").
		WithBody(request).
		WithHandler(handler.Add())

	body := testRequest.Verify(t, http.StatusCreated)

	actual := model.HouseResponse{}

	json.Unmarshal(body, &actual)

	assert.Equal(t,
		model.HouseResponse{
			Id:          actual.Id,
			Name:        "Test House",
			Country:     "Ukraine",
			City:        "City",
			StreetLine1: "StreetLine1",
			StreetLine2: "StreetLine2",
			UserId:      userId,
		}, actual)
}

func Test_Add_WithErrorFromService(t *testing.T) {
	handler := handlerGenerator()

	request := generateCreateHouseRequest()

	houses.On("Add", request).Return(model.HouseResponse{}, errors.New("error"))

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

	houseResponse := generateHouseResponse()

	houses.On("FindById", houseResponse.Id).Return(houseResponse, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/house/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", houseResponse.Id.String())

	body := testRequest.Verify(t, http.StatusOK)

	var responses model.HouseResponse
	json.Unmarshal(body, &responses)

	assert.Equal(t, houseResponse, responses)
}

func Test_FindById_WithErrorFromService(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	houses.On("FindById", id).Return(model.HouseResponse{}, errors.New("error"))

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/house/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", id.String())

	body := testRequest.Verify(t, http.StatusNotFound)

	assert.Equal(t, "error\n", string(body))
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

	houseResponse := generateHouseResponse()

	houseResponses := []model.HouseResponse{houseResponse}
	houses.On("FindByUserId", houseResponse.Id).Return(houseResponses)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/house/user/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByUserId()).
		WithVar("id", houseResponse.Id.String())

	body := testRequest.Verify(t, http.StatusOK)

	var responses []model.HouseResponse
	json.Unmarshal(body, &responses)

	assert.Equal(t, houseResponses, responses)
}

func Test_FindByUserId_WithEmptyResponse(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	houses.On("FindByUserId", id).Return([]model.HouseResponse{})

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/house/user/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByUserId()).
		WithVar("id", id.String())

	body := testRequest.Verify(t, http.StatusOK)

	var responses []model.HouseResponse
	json.Unmarshal(body, &responses)

	assert.Equal(t, []model.HouseResponse{}, responses)
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

func generateCreateHouseRequest() model.CreateHouseRequest {
	return model.CreateHouseRequest{
		Name:        "Test House",
		Country:     "UA",
		City:        "City",
		StreetLine1: "StreetLine1",
		StreetLine2: "StreetLine2",
		UserId:      userId,
	}
}

func generateHouseResponse() model.HouseResponse {
	return model.HouseResponse{
		Id:          uuid.New(),
		Name:        "Test Name",
		Country:     "Ukraine",
		City:        "City",
		StreetLine1: "StreetLine1",
		StreetLine2: "StreetLine2",
		UserId:      userId,
	}
}
