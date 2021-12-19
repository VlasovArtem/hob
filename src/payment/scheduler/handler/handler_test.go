package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/payment/scheduler/mocks"
	paymentScheduler "github.com/VlasovArtem/hob/src/payment/scheduler/model"
	"github.com/VlasovArtem/hob/src/test/testhelper"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

var paymentsScheduler *mocks.PaymentSchedulerService

func handlerGenerator() PaymentSchedulerHandler {
	paymentsScheduler = new(mocks.PaymentSchedulerService)

	return NewPaymentSchedulerHandler(paymentsScheduler)
}

func Test_Add(t *testing.T) {
	handler := handlerGenerator()

	request := mocks.GenerateCreatePaymentSchedulerRequest()

	paymentsScheduler.On("Add", request).Return(request.ToEntity().ToDto(), nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/scheduler").
		WithMethod("POST").
		WithHandler(handler.Add()).
		WithBody(request)

	responseByteArray := testRequest.Verify(t, http.StatusCreated)

	actual := paymentScheduler.PaymentSchedulerDto{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, mocks.GeneratePaymentSchedulerResponse(actual.Id), actual)
}

func Test_Add_WithInvalidRequest(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/scheduler").
		WithMethod("POST").
		WithHandler(handler.Add())

	testRequest.Verify(t, http.StatusBadRequest)
}

func Test_Add_WithErrorFromService(t *testing.T) {
	handler := handlerGenerator()

	request := mocks.GenerateCreatePaymentSchedulerRequest()

	expected := errors.New("error")

	paymentsScheduler.On("Add", request).Return(paymentScheduler.PaymentSchedulerDto{}, expected)

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

	paymentResponse := mocks.GeneratePaymentSchedulerResponse(id)

	paymentsScheduler.On("FindById", id).
		Return(paymentResponse, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/scheduler/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	actual := paymentScheduler.PaymentSchedulerDto{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, paymentResponse, actual)
}

func Test_FindById_WithError(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	expected := errors.New("error")

	paymentsScheduler.On("FindById", id).
		Return(paymentScheduler.PaymentSchedulerDto{}, expected)

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

	response := mocks.GeneratePaymentSchedulerResponse(id)

	paymentResponses := []paymentScheduler.PaymentSchedulerDto{response}

	paymentsScheduler.On("FindByHouseId", id).
		Return(paymentResponses, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/scheduler/house/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByHouseId()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	actual := []paymentScheduler.PaymentSchedulerDto{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, paymentResponses, actual)
}

func Test_FindByHouseId_WithEmptyResponse(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	paymentResponses := []paymentScheduler.PaymentSchedulerDto{}

	paymentsScheduler.On("FindByHouseId", id).
		Return(paymentResponses, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/scheduler/house/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByHouseId()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	actual := []paymentScheduler.PaymentSchedulerDto{}

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

	response := mocks.GeneratePaymentSchedulerResponse(id)

	paymentResponses := []paymentScheduler.PaymentSchedulerDto{response}

	paymentsScheduler.On("FindByUserId", id).
		Return(paymentResponses, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/scheduler/user/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByUserId()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	var actual []paymentScheduler.PaymentSchedulerDto

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, paymentResponses, actual)
}

func Test_FindByUserId_WithEmptyResponse(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	var paymentResponses []paymentScheduler.PaymentSchedulerDto

	paymentsScheduler.On("FindByUserId", id).
		Return(paymentResponses, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/scheduler/user/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByUserId()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	var actual []paymentScheduler.PaymentSchedulerDto

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

func Test_FindByProviderId(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	response := mocks.GeneratePaymentSchedulerResponse(id)

	paymentResponses := []paymentScheduler.PaymentSchedulerDto{response}

	paymentsScheduler.On("FindByProviderId", id).
		Return(paymentResponses, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/scheduler/provider/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByProviderId()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	var actual []paymentScheduler.PaymentSchedulerDto

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, paymentResponses, actual)
}

func Test_FindByProviderId_WithEmptyResponse(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	var paymentResponses []paymentScheduler.PaymentSchedulerDto

	paymentsScheduler.On("FindByProviderId", id).
		Return(paymentResponses, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/scheduler/provider/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByProviderId()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	var actual []paymentScheduler.PaymentSchedulerDto

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, paymentResponses, actual)
}

func Test_FindByProviderId_WithInvalidParameter(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/scheduler/provider/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByProviderId()).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "the id is not valid id\n", string(responseByteArray))
}

func Test_Update(t *testing.T) {
	handler := handlerGenerator()

	request := mocks.GenerateUpdatePaymentSchedulerRequest()
	id := uuid.New()

	paymentsScheduler.On("Update", id, request).Return(nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/scheduler/{id}").
		WithMethod("PUT").
		WithHandler(handler.Update()).
		WithBody(request).
		WithVar("id", id.String())

	testRequest.Verify(t, http.StatusOK)
}

func Test_Update_WithInvalidId(t *testing.T) {
	handler := handlerGenerator()

	request := mocks.GenerateUpdatePaymentSchedulerRequest()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/scheduler/{id}").
		WithMethod("PUT").
		WithHandler(handler.Update()).
		WithBody(request).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "the id is not valid id\n", string(responseByteArray))

	paymentsScheduler.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func Test_Update_WithInvalidRequest(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/scheduler/{id}").
		WithMethod("PUT").
		WithHandler(handler.Update()).
		WithVar("id", uuid.New().String())

	testRequest.Verify(t, http.StatusBadRequest)
}

func Test_Update_WithErrorFromService(t *testing.T) {
	handler := handlerGenerator()

	request := mocks.GenerateUpdatePaymentSchedulerRequest()
	id := uuid.New()

	expected := errors.New("error")

	paymentsScheduler.On("Update", id, request).Return(expected)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/scheduler/{id}").
		WithMethod("PUT").
		WithHandler(handler.Update()).
		WithVar("id", id.String()).
		WithBody(request)

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, fmt.Sprintf("%s\n", expected.Error()), string(responseByteArray))
}
