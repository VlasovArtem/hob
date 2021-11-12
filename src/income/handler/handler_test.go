package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"income/model"
	"net/http"
	"test/mock"
	"test/testhelper"
	"testing"
	"time"
)

var (
	incomes *mock.IncomeServiceMock
	houseId = testhelper.ParseUUID("73998efa-afa9-4923-ba4b-3d1d60241823")
	date    = time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
)

func handlerGenerator() IncomeHandler {
	incomes = new(mock.IncomeServiceMock)

	return NewIncomeHandler(incomes)
}

func Test_AddIncome(t *testing.T) {
	handler := handlerGenerator()

	request := generateCreateIncomeRequest()

	incomes.On("Add", request).Return(request.ToEntity().ToResponse(), nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/income").
		WithMethod("POST").
		WithHandler(handler.Add()).
		WithBody(request)

	responseByteArray := testRequest.Verify(t, http.StatusCreated)

	actual := model.IncomeResponse{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, model.IncomeResponse{
		Income: model.Income{
			Id:          actual.Id,
			Name:        "Name",
			Date:        date,
			Description: "Description",
			Sum:         100.1,
			HouseId:     houseId,
		},
	}, actual)
}

func Test_AddIncome_WithInvalidRequest(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/income").
		WithMethod("POST").
		WithHandler(handler.Add())

	testRequest.Verify(t, http.StatusBadRequest)
}

func Test_AddIncome_WithErrorFromService(t *testing.T) {
	handler := handlerGenerator()

	request := generateCreateIncomeRequest()

	err := errors.New("error")
	incomes.On("Add", request).Return(model.IncomeResponse{}, err)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/income").
		WithMethod("POST").
		WithHandler(handler.Add()).
		WithBody(request)

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "error\n", string(responseByteArray))
}

func Test_FindById(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	response := generateIncomeResponse(id)

	incomes.On("FindById", id).
		Return(response, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/income/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	actual := model.IncomeResponse{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, response, actual)
}

func Test_FindById_WithError(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	expected := errors.New("error")

	incomes.On("FindById", id).
		Return(model.IncomeResponse{}, expected)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/income/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, fmt.Sprintf("%s\n", expected.Error()), string(responseByteArray))
}

func Test_FindById_WithInvalidParameter(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/income/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "the id is not valid id\n", string(responseByteArray))
}

func Test_FindByHouseId(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	meterResponse := generateIncomeResponse(id)

	incomes.On("FindByHouseId", id).
		Return(meterResponse, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meter/house/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByHouseId()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	actual := model.IncomeResponse{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, meterResponse, actual)
}

func Test_FindByHouseId_WithError(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	expected := errors.New("error")

	incomes.On("FindByHouseId", id).
		Return(model.IncomeResponse{}, expected)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meter/house/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByHouseId()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	assert.Empty(t, responseByteArray)
}

func Test_FindByHouseId_WithInvalidParameter(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meter/house/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByHouseId()).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "the id is not valid id\n", string(responseByteArray))
}

func generateCreateIncomeRequest() model.CreateIncomeRequest {
	return model.CreateIncomeRequest{
		Name:        "Name",
		Date:        date,
		Description: "Description",
		Sum:         100.1,
		HouseId:     houseId,
	}
}

func generateIncomeResponse(id uuid.UUID) model.IncomeResponse {
	return model.IncomeResponse{
		Income: model.Income{
			Id:          id,
			Name:        "Name",
			Date:        date,
			Description: "Description",
			Sum:         100.1,
			HouseId:     houseId,
		},
	}
}
