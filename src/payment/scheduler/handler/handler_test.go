package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	paymentModel "payment/model"
	paymentScheduler "payment/scheduler/model"
	scheduler2 "scheduler"
	"test/mock"
	"test/testhelper"
	"testing"
)

var (
	paymentsScheduler *mock.PaymentSchedulerServiceMock
	houseId           = testhelper.ParseUUID("84d801b0-75fd-4304-aab4-97c8c46356bb")
	userId            = testhelper.ParseUUID("eddbd5ca-cc87-4cbd-9753-7bbb11bdef83")
)

func handlerGenerator() PaymentSchedulerHandler {
	paymentsScheduler = new(mock.PaymentSchedulerServiceMock)

	return NewPaymentSchedulerHandler(paymentsScheduler)
}

func Test_Add(t *testing.T) {
	handler := handlerGenerator()

	request := generateCreatePaymentSchedulerRequest()

	paymentsScheduler.On("Add", request).Return(request.ToEntity().ToResponse(), nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/scheduler/scheduler").
		WithMethod("POST").
		WithHandler(handler.Add()).
		WithBody(request)

	responseByteArray := testRequest.Verify(t, http.StatusCreated)

	actual := paymentScheduler.PaymentSchedulerResponse{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, generatePaymentSchedulerResponse(actual.Id), actual)
}

func Test_Add_WithInvalidRequest(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/scheduler/scheduler").
		WithMethod("POST").
		WithHandler(handler.Add())

	testRequest.Verify(t, http.StatusBadRequest)
}

func Test_Add_WithErrorFromService(t *testing.T) {
	handler := handlerGenerator()

	request := generateCreatePaymentSchedulerRequest()

	expected := errors.New("error")

	paymentsScheduler.On("Add", request).Return(paymentScheduler.PaymentSchedulerResponse{}, expected)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/scheduler").
		WithMethod("POST").
		WithHandler(handler.Add()).
		WithBody(request)

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, fmt.Sprintf("%s\n", expected.Error()), string(responseByteArray))
}

func Test_Remove(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	paymentsScheduler.On("Remove", id).Return(nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/scheduler/{id}").
		WithMethod("DELETE").
		WithHandler(handler.Remove()).
		WithVar("id", id.String())

	testRequest.Verify(t, http.StatusNoContent)
}

func Test_Remove_WithMissingParameter(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/scheduler/{id}").
		WithMethod("DELETE").
		WithHandler(handler.Remove())

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "parameter 'id' not found\n", string(responseByteArray))
}

func Test_Remove_WithInvalidParameter(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/scheduler/{id}").
		WithMethod("DELETE").
		WithHandler(handler.Remove()).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "the id is not valid id\n", string(responseByteArray))
}

func Test_FindById(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	paymentResponse := generatePaymentSchedulerResponse(id)

	paymentsScheduler.On("FindById", id).
		Return(paymentResponse, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/scheduler/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	actual := paymentScheduler.PaymentSchedulerResponse{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, paymentResponse, actual)
}

func Test_FindById_WithError(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	expected := errors.New("error")

	paymentsScheduler.On("FindById", id).
		Return(paymentScheduler.PaymentSchedulerResponse{}, expected)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/scheduler/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, fmt.Sprintf("%s\n", expected.Error()), string(responseByteArray))
}

func Test_FindById_WithInvalidParameter(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/scheduler/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "the id is not valid id\n", string(responseByteArray))
}

func Test_FindByHouseId(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	paymentResponses := []paymentScheduler.PaymentSchedulerResponse{generatePaymentSchedulerResponse(id)}

	paymentsScheduler.On("FindByHouseId", id).
		Return(paymentResponses, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/scheduler/house/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByHouseId()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	actual := []paymentScheduler.PaymentSchedulerResponse{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, paymentResponses, actual)
}

func Test_FindByHouseId_WithEmptyResponse(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	paymentResponses := []paymentScheduler.PaymentSchedulerResponse{}

	paymentsScheduler.On("FindByHouseId", id).
		Return(paymentResponses, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/scheduler/house/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByHouseId()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	actual := []paymentScheduler.PaymentSchedulerResponse{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, paymentResponses, actual)
}

func Test_FindByHouseId_WithInvalidParameter(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/scheduler/house/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByHouseId()).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "the id is not valid id\n", string(responseByteArray))
}

func Test_FindByUserId(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	paymentResponses := []paymentScheduler.PaymentSchedulerResponse{generatePaymentSchedulerResponse(id)}

	paymentsScheduler.On("FindByUserId", id).
		Return(paymentResponses, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/scheduler/user/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByUserId()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	actual := []paymentScheduler.PaymentSchedulerResponse{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, paymentResponses, actual)
}

func Test_FindByUserId_WithEmptyResponse(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	paymentResponses := []paymentScheduler.PaymentSchedulerResponse{}

	paymentsScheduler.On("FindByUserId", id).
		Return(paymentResponses, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/scheduler/user/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByUserId()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	actual := []paymentScheduler.PaymentSchedulerResponse{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, paymentResponses, actual)
}

func Test_FindByUserId_WithInvalidParameter(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/scheduler/user/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByUserId()).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "the id is not valid id\n", string(responseByteArray))
}

func generateCreatePaymentSchedulerRequest() paymentScheduler.CreatePaymentSchedulerRequest {
	return paymentScheduler.CreatePaymentSchedulerRequest{
		Name:        "Test Payment",
		Description: "Test Payment Description",
		HouseId:     houseId,
		UserId:      userId,
		Sum:         1000,
		Spec:        scheduler2.DAILY,
	}
}

func generatePaymentSchedulerResponse(id uuid.UUID) paymentScheduler.PaymentSchedulerResponse {
	return paymentScheduler.PaymentSchedulerResponse{
		PaymentScheduler: paymentScheduler.PaymentScheduler{
			Payment: paymentModel.Payment{
				Id:          id,
				Name:        "Test Payment",
				Description: "Test Payment Description",
				HouseId:     houseId,
				UserId:      userId,
				Sum:         1000,
			},
			Spec: scheduler2.DAILY,
		},
	}
}
