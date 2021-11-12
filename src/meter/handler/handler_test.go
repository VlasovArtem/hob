package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"meter/model"
	"net/http"
	"test/mock"
	"test/testhelper"
	"testing"
)

var (
	meters    *mock.MeterServiceMock
	paymentId = testhelper.ParseUUID("8db24b37-6978-4c7f-ae8d-516eaabb323b")
)

func handlerGenerator() MeterHandler {
	meters = new(mock.MeterServiceMock)

	return NewMeterHandler(meters)
}

func Test_AddMeter(t *testing.T) {
	handler := handlerGenerator()

	request := generateCreateMeterRequest()

	meters.On("AddMeter", request).Return(request.ToEntity().ToResponse(), nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meter").
		WithMethod("POST").
		WithHandler(handler.AddMeter()).
		WithBody(request)

	responseByteArray := testRequest.Verify(t, http.StatusCreated)

	actual := model.MeterResponse{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, model.MeterResponse{
		Meter: model.Meter{
			Id:   actual.Id,
			Name: "Name",
			Details: map[string]float64{
				"first":  1.1,
				"second": 2.2,
			},
			Description: "Description",
			PaymentId:   paymentId,
		},
	}, actual)
}

func Test_AddMeter_WithInvalidRequest(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meter").
		WithMethod("POST").
		WithHandler(handler.AddMeter())

	testRequest.Verify(t, http.StatusBadRequest)
}

func Test_AddMeter_WithErrorFromService(t *testing.T) {
	handler := handlerGenerator()

	request := generateCreateMeterRequest()

	err := errors.New("error")
	meters.On("AddMeter", request).Return(model.MeterResponse{}, err)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meter").
		WithMethod("POST").
		WithHandler(handler.AddMeter()).
		WithBody(request)

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "error\n", string(responseByteArray))
}

func Test_FindById(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	meterResponse := generateMeterResponse(id)

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
	handler := handlerGenerator()

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
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meter/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "the id is not valid id\n", string(responseByteArray))
}

func Test_FindByPaymentId(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	meterResponse := generateMeterResponse(id)

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
	handler := handlerGenerator()

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
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meter/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByPaymentId()).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "the id is not valid id\n", string(responseByteArray))
}

func generateCreateMeterRequest() model.CreateMeterRequest {
	return model.CreateMeterRequest{
		Name: "Name",
		Details: map[string]float64{
			"first":  1.1,
			"second": 2.2,
		},
		Description: "Description",
		PaymentId:   paymentId,
	}
}

func generateMeterResponse(id uuid.UUID) model.MeterResponse {
	return model.MeterResponse{
		Meter: model.Meter{
			Id:   id,
			Name: "Name",
			Details: map[string]float64{
				"first":  1.1,
				"second": 2.2,
			},
			Description: "Description",
			PaymentId:   paymentId,
		},
	}
}
