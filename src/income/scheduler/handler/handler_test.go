package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	incomeModel "income/model"
	incomeSchedulerModel "income/scheduler/model"
	"net/http"
	scheduler2 "scheduler"
	"test/mock"
	"test/testhelper"
	"testing"
)

var (
	incomesScheduler *mock.IncomeSchedulerServiceMock
	houseId          = testhelper.ParseUUID("d0495341-d8fe-4b2b-af9d-73a516cde342")
)

func handlerGenerator() IncomeSchedulerHandler {
	incomesScheduler = new(mock.IncomeSchedulerServiceMock)

	return NewIncomeSchedulerHandler(incomesScheduler)
}

func Test_Add(t *testing.T) {
	handler := handlerGenerator()

	request := generateCreateIncomeSchedulerRequest()

	incomesScheduler.On("Add", request).Return(request.ToEntity().ToResponse(), nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/income/scheduler").
		WithMethod("POST").
		WithHandler(handler.Add()).
		WithBody(request)

	responseByteArray := testRequest.Verify(t, http.StatusCreated)

	actual := incomeSchedulerModel.IncomeSchedulerResponse{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, generateIncomeSchedulerResponse(actual.Id), actual)
}

func Test_Add_WithInvalidRequest(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/income/scheduler").
		WithMethod("POST").
		WithHandler(handler.Add())

	testRequest.Verify(t, http.StatusBadRequest)
}

func Test_Add_WithErrorFromService(t *testing.T) {
	handler := handlerGenerator()

	request := generateCreateIncomeSchedulerRequest()

	expected := errors.New("error")

	incomesScheduler.On("Add", request).Return(incomeSchedulerModel.IncomeSchedulerResponse{}, expected)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/income/scheduler").
		WithMethod("POST").
		WithHandler(handler.Add()).
		WithBody(request)

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, fmt.Sprintf("%s\n", expected.Error()), string(responseByteArray))
}

func Test_Remove(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	incomesScheduler.On("Remove", id).Return(nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/income/scheduler/{id}").
		WithMethod("DELETE").
		WithHandler(handler.Remove()).
		WithVar("id", id.String())

	testRequest.Verify(t, http.StatusNoContent)
}

func Test_Remove_WithMissingParameter(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/income/scheduler/{id}").
		WithMethod("DELETE").
		WithHandler(handler.Remove())

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "parameter 'id' not found\n", string(responseByteArray))
}

func Test_Remove_WithInvalidParameter(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/income/scheduler/{id}").
		WithMethod("DELETE").
		WithHandler(handler.Remove()).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "the id is not valid id\n", string(responseByteArray))
}

func Test_FindById(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	paymentResponse := generateIncomeSchedulerResponse(id)

	incomesScheduler.On("FindById", id).
		Return(paymentResponse, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/income/scheduler/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	actual := incomeSchedulerModel.IncomeSchedulerResponse{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, paymentResponse, actual)
}

func Test_FindById_WithError(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	expected := errors.New("error")

	incomesScheduler.On("FindById", id).
		Return(incomeSchedulerModel.IncomeSchedulerResponse{}, expected)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/income/scheduler/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, fmt.Sprintf("%s\n", expected.Error()), string(responseByteArray))
}

func Test_FindById_WithInvalidParameter(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/income/scheduler/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "the id is not valid id\n", string(responseByteArray))
}

func Test_FindByHouseId(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	paymentResponse := generateIncomeSchedulerResponse(id)

	incomesScheduler.On("FindByHouseId", id).
		Return(paymentResponse, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/income/scheduler/house/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByHouseId()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	actual := incomeSchedulerModel.IncomeSchedulerResponse{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, paymentResponse, actual)
}

func Test_FindByHouseId_WithInvalidParameter(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/income/scheduler/house/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByHouseId()).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "the id is not valid id\n", string(responseByteArray))
}

func generateCreateIncomeSchedulerRequest() incomeSchedulerModel.CreateIncomeSchedulerRequest {
	return incomeSchedulerModel.CreateIncomeSchedulerRequest{
		Name:        "Test Income",
		Description: "Test Income Description",
		HouseId:     houseId,
		Sum:         1000,
		Spec:        scheduler2.DAILY,
	}
}

func generateIncomeSchedulerResponse(id uuid.UUID) incomeSchedulerModel.IncomeSchedulerResponse {
	return incomeSchedulerModel.IncomeSchedulerResponse{
		IncomeScheduler: incomeSchedulerModel.IncomeScheduler{
			Income: incomeModel.Income{
				Id:          id,
				Name:        "Test Income",
				Description: "Test Income Description",
				HouseId:     houseId,
				Sum:         1000,
			},
			Spec: scheduler2.DAILY,
		},
	}
}
