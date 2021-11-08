package handler

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"house/model"
	"house/service"
	"net/http"
	"net/http/httptest"
	"testing"
)

func generateCreateHouseRequest() model.CreateHouseRequest {
	return model.CreateHouseRequest{
		Name:        "Test House",
		Country:     "Country",
		City:        "City",
		StreetLine1: "StreetLine1",
		StreetLine2: "StreetLine2",
	}
}

func generateHouse(id uuid.UUID) model.House {
	return model.House{
		Id:          id,
		Name:        "Test House",
		Country:     "Country",
		City:        "City",
		StreetLine1: "StreetLine1",
		StreetLine2: "StreetLine2",
		Deleted:     false,
	}
}

func TestAddHouseHandlerWithNotValidRequest(t *testing.T) {
	handler := AddHouseHandler()

	request := httptest.NewRequest("POST", "https://test.com/api/v1/house", nil)
	recorder := httptest.NewRecorder()

	handler(recorder, request)

	response := recorder.Result()

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestAddHouseHandler(t *testing.T) {
	oldHousesMap := make(map[uuid.UUID]model.House)

	for _, house := range service.FindAllHouses() {
		oldHousesMap[house.Id] = house
	}

	house := generateCreateHouseRequest()

	houseJson, _ := json.Marshal(house)

	buffer := bytes.Buffer{}

	buffer.Write(houseJson)

	handler := AddHouseHandler()

	request := httptest.NewRequest("POST", "https://test.com/api/v1/house", &buffer)

	responseByteArray := verifyWithResponse(t, request, handler, http.StatusCreated)

	responseBody := model.CreateHouseResponse{}

	json.Unmarshal(responseByteArray, &responseBody)

	for _, house := range service.FindAllHouses() {
		if _, found := oldHousesMap[house.Id]; !found {
			assert.Equal(t, model.CreateHouseResponse{Id: house.Id}, responseBody)

			return
		}
	}
}

func TestFindAllHousesHandler(t *testing.T) {
	houses := service.FindAllHouses()

	handler := FindAllHousesHandler()

	request := httptest.NewRequest("GET", "https://test.com/api/v1/house", nil)

	response := verifyWithResponse(t, request, handler, 200)

	var responseContent []model.House

	json.Unmarshal(response, &responseContent)

	assert.Equal(t, houses, responseContent)
}

func TestFindHouseByIdHandler(t *testing.T) {
	house := generateHouse(uuid.New())
	service.AddHouse(house)

	handlerFunc := FindHouseByIdHandler()

	request := httptest.NewRequest("GET", "https://test.com/api/v1/house/{id}", nil)

	request = mux.SetURLVars(request, map[string]string{
		"id": house.Id.String(),
	})

	response := verifyWithResponse(t, request, handlerFunc, 200)

	var responseObject model.House

	json.Unmarshal(response, &responseObject)

	assert.Equal(t, house, responseObject)
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
			handlerFunc := FindHouseByIdHandler()

			request := httptest.NewRequest("GET", "https://test.com/api/v1/house/{id}", nil)

			request = mux.SetURLVars(request, map[string]string{
				"id": tt.args.id,
			})

			recorder := httptest.NewRecorder()

			handlerFunc(recorder, request)

			response := recorder.Result()

			assert.Equal(t, tt.args.code, response.StatusCode, "Status code should be the same")
		})
	}
}

func verifyWithResponse(t *testing.T, request *http.Request, handler http.HandlerFunc, expectedStatus int) []byte {
	recorder := httptest.NewRecorder()

	handler(recorder, request)

	response := recorder.Result()

	assert.Equal(t, expectedStatus, response.StatusCode)

	return readBytes(response)
}

func readBytes(response *http.Response) []byte {
	buffer := bytes.Buffer{}

	buffer.ReadFrom(response.Body)

	return buffer.Bytes()
}
