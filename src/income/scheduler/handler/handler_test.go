package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	incomeModel "income/model"
	"income/scheduler/mocks"
	incomeSchedulerModel "income/scheduler/model"
	"net/http"
	scheduler2 "scheduler"
	"test/testhelper"
	"testing"
)

var (
	incomesScheduler *mocks.IncomeSchedulerService
	houseId          = testhelper.ParseUUID("d0495341-d8fe-4b2b-af9d-73a516cde342")
)

func handlerGenerator() IncomeSchedulerHandler {
	incomesScheduler = new(mocks.IncomeSchedulerService)

	return NewIncomeSchedulerHandler(incomesScheduler)
}

func Test_Add(t *testing.T) {
	handler := handlerGenerator()

	request := generateCreateIncomeSchedulerRequest()

	incomesScheduler.On("Add", request).Return(request.ToEntity().ToDto(), nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/income/scheduler").
		WithMethod("POST").
		WithHandler(handler.Add()).
		WithBody(request)

	responseByteArray := testRequest.Verify(t, http.StatusCreated)

	actual := incomeSchedulerModel.IncomeSchedulerDto{}

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

	incomesScheduler.On("Add", request).Return(incomeSchedulerModel.IncomeSchedulerDto{}, expected)

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

	incomeSchedulerResponse := generateIncomeSchedulerResponse(id)

	incomesScheduler.On("FindById", id).
		Return(incomeSchedulerResponse, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/income/scheduler/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	actual := incomeSchedulerModel.IncomeSchedulerDto{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, incomeSchedulerResponse, actual)
}

func Test_FindById_WithError(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	expected := errors.New("error")

	incomesScheduler.On("FindById", id).
		Return(incomeSchedulerModel.IncomeSchedulerDto{}, expected)

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

	incomeSchedulerResponse := generateIncomeSchedulerResponse(id)

	incomesScheduler.On("FindByHouseId", id).
		Return([]incomeSchedulerModel.IncomeSchedulerDto{incomeSchedulerResponse})

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/income/scheduler/house/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByHouseId()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	var actual []incomeSchedulerModel.IncomeSchedulerDto

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, []incomeSchedulerModel.IncomeSchedulerDto{incomeSchedulerResponse}, actual)
}

func Test_FindByHouseId_WithEmptyResponse(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	incomesScheduler.On("FindByHouseId", id).
		Return([]incomeSchedulerModel.IncomeSchedulerDto{})

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/income/scheduler/house/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByHouseId()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	var actual []incomeSchedulerModel.IncomeSchedulerDto

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, []incomeSchedulerModel.IncomeSchedulerDto{}, actual)
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

func Test_FindByHouseId_WithMissingParameter(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/income/scheduler/house/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByHouseId())

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "parameter 'id' not found\n", string(responseByteArray))
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

func generateIncomeSchedulerResponse(id uuid.UUID) incomeSchedulerModel.IncomeSchedulerDto {
	return incomeSchedulerModel.IncomeSchedulerDto{
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
