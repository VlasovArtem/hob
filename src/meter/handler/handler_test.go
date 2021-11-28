package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/meter/mocks"
	"github.com/VlasovArtem/hob/src/meter/model"
	"github.com/VlasovArtem/hob/src/test/testhelper"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var (
	meters *mocks.MeterService
)

func generateHandler() MeterHandler {
	meters = new(mocks.MeterService)

	return NewMeterHandler(meters)
}

func Test_AddMeter(t *testing.T) {
	handler := generateHandler()

	request := mocks.GenerateCreateMeterRequest()

	meters.On("Add", request).Return(request.ToEntity().ToResponse(), nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meter").
		WithMethod("POST").
		WithHandler(handler.Add()).
		WithBody(request)

	responseByteArray := testRequest.Verify(t, http.StatusCreated)

	actual := model.MeterResponse{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, model.MeterResponse{
		Id:   actual.Id,
		Name: "Name",
		Details: map[string]float64{
			"first":  1.1,
			"second": 2.2,
		},
		Description: "Description",
		PaymentId:   request.PaymentId,
		HouseId:     request.HouseId,
	}, actual)
}

func Test_AddMeter_WithInvalidRequest(t *testing.T) {
	handler := generateHandler()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meter").
		WithMethod("POST").
		WithHandler(handler.Add())

	testRequest.Verify(t, http.StatusBadRequest)
}

func Test_AddMeter_WithErrorFromService(t *testing.T) {
	handler := generateHandler()

	request := mocks.GenerateCreateMeterRequest()

	err := errors.New("error")
	meters.On("Add", request).Return(model.MeterResponse{}, err)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meter").
		WithMethod("POST").
		WithHandler(handler.Add()).
		WithBody(request)

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "error\n", string(responseByteArray))
}

func Test_FindById(t *testing.T) {
	handler := generateHandler()

	id := uuid.New()

	meterResponse := mocks.GenerateMeterResponse(id)

	meters.On("FindById", id).
		Return(meterResponse, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meter/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	actual := model.MeterResponse{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, meterResponse, actual)
}

func Test_FindById_WithError(t *testing.T) {
	handler := generateHandler()

	id := uuid.New()

	expected := errors.New("error")

	meters.On("FindById", id).
		Return(model.MeterResponse{}, expected)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meter/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, fmt.Sprintf("%s\n", expected.Error()), string(responseByteArray))
}

func Test_FindById_WithInvalidParameter(t *testing.T) {
	handler := generateHandler()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meter/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "the id is not valid id\n", string(responseByteArray))
}

func Test_FindByPaymentId(t *testing.T) {
	handler := generateHandler()

	id := uuid.New()

	meterResponse := mocks.GenerateMeterResponse(id)

	meters.On("FindByPaymentId", id).
		Return(meterResponse, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meter/payment/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByPaymentId()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	actual := model.MeterResponse{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, meterResponse, actual)
}

func Test_FindByPaymentId_WithError(t *testing.T) {
	handler := generateHandler()

	id := uuid.New()

	expected := errors.New("error")

	meters.On("FindByPaymentId", id).
		Return(model.MeterResponse{}, expected)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meter/payment/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByPaymentId()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	assert.Empty(t, responseByteArray)
}

func Test_FindByPaymentId_WithInvalidParameter(t *testing.T) {
	handler := generateHandler()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meter/payment/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByPaymentId()).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "the id is not valid id\n", string(responseByteArray))
}

func Test_FindByHouseId(t *testing.T) {
	handler := generateHandler()

	id := uuid.New()

	meterResponse := []model.MeterResponse{mocks.GenerateMeterResponse(id)}

	meters.On("FindByHouseId", meterResponse[0].HouseId).
		Return(meterResponse, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meter/house/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByHouseId()).
		WithVar("id", meterResponse[0].HouseId.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	var actual []model.MeterResponse

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, meterResponse, actual)
}

func Test_FindByHouseId_WithEmptyResponse(t *testing.T) {
	handler := generateHandler()

	id := uuid.New()

	var meterResponse []model.MeterResponse

	meters.On("FindByHouseId", id).
		Return(meterResponse, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meter/house/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByHouseId()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	var actual []model.MeterResponse

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, meterResponse, actual)
}

func Test_FindByHouseId_WithInvalidParameter(t *testing.T) {
	handler := generateHandler()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meter/house/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByPaymentId()).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "the id is not valid id\n", string(responseByteArray))
}

func Test_FindByHouseId_WithMissingParameter(t *testing.T) {
	handler := generateHandler()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meter/house/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByPaymentId())

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "parameter 'id' not found\n", string(responseByteArray))
}
