package handler

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"house/model"
	"house/service"
	"net/http"
	"test"
	"test/testhelper"
	"testing"
)

var countryService = testhelper.InitCountryService()

var houseService, handler = func() (service.HouseService, HouseHandler) {
	houseService := service.NewHouseService(countryService)

	return houseService, NewHouseHandler(houseService)
}()

func TestAddHouseHandlerWithNotValidRequest(t *testing.T) {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/house").
		WithMethod("POST").
		WithHandler(handler.Add())

	testRequest.Verify(t, http.StatusBadRequest)
}

func TestAddHouseHandler(t *testing.T) {
	createRequest := test.GenerateCreateHouseRequest()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/house").
		WithMethod("POST").
		WithBody(createRequest).
		WithHandler(handler.Add())

	body := testRequest.Verify(t, http.StatusCreated)

	responseBody := model.HouseResponse{}

	json.Unmarshal(body, &responseBody)

	expectedResponse := test.GenerateHouseResponse(responseBody.Id, responseBody.Name)

	assert.Equal(t, expectedResponse, responseBody)
}

func TestFindAllHousesHandler(t *testing.T) {
	houses := houseService.FindAll()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/house").
		WithMethod("GET").
		WithHandler(handler.FindAll())

	body := testRequest.Verify(t, http.StatusOK)

	var responses []model.HouseResponse
	json.Unmarshal(body, &responses)

	assert.Equal(t, houses, responses)
}

func TestFindHouseByIdHandler(t *testing.T) {
	houseRequest := test.GenerateCreateHouseRequest()
	houseResponse, err := houseService.Add(houseRequest)

	assert.Nil(t, err)

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

func TestFindHouseByIdHandlerWithError(t *testing.T) {
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
				WithURL("https://test.com/api/v1/house/{id}").
				WithMethod("GET").
				WithHandler(handler.FindById()).
				WithVar("id", tt.args.id)

			testRequest.Verify(t, tt.args.code)
		})
	}
}
