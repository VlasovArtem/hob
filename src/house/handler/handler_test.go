package handler

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"helper"
	"helper/testhelper"
	"house/model"
	"house/service"
	"net/http"
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
		WithHandler(handler.AddHouseHandler())

	testRequest.VerifyWithStatus(t, http.StatusBadRequest)
}

func TestAddHouseHandler(t *testing.T) {
	createRequest := helper.GenerateCreateHouseRequest()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/house").
		WithMethod("POST").
		WithBody(createRequest).
		WithHandler(handler.AddHouseHandler())

	body := testRequest.VerifyWithStatus(t, http.StatusCreated)

	responseBody := model.HouseResponse{}

	json.Unmarshal(body, &responseBody)

	expectedResponse := helper.GenerateHouseResponse(responseBody.Id, responseBody.Name)

	assert.Equal(t, expectedResponse, responseBody)
}

func TestFindAllHousesHandler(t *testing.T) {
	houses := houseService.FindAllHouses()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/house").
		WithMethod("GET").
		WithHandler(handler.FindAllHousesHandler())

	body := testRequest.VerifyWithStatus(t, http.StatusOK)

	var responses []model.HouseResponse
	json.Unmarshal(body, &responses)

	assert.Equal(t, houses, responses)
}

func TestFindHouseByIdHandler(t *testing.T) {
	houseRequest := helper.GenerateCreateHouseRequest()
	err, houseResponse := houseService.AddHouse(houseRequest)

	assert.Nil(t, err)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/house/{id}").
		WithMethod("GET").
		WithHandler(handler.FindHouseByIdHandler()).
		WithVar("id", houseResponse.Id.String())

	body := testRequest.VerifyWithStatus(t, http.StatusOK)

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
				WithHandler(handler.FindHouseByIdHandler()).
				WithVar("id", tt.args.id)

			testRequest.VerifyWithStatus(t, tt.args.code)
		})
	}
}
