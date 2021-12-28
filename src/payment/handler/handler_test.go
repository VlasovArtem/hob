package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/payment/mocks"
	"github.com/VlasovArtem/hob/src/payment/model"
	"github.com/VlasovArtem/hob/src/test/testhelper"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

var payments *mocks.PaymentService

func handlerGenerator() PaymentHandler {
	payments = new(mocks.PaymentService)

	return NewPaymentHandler(payments)
}

func Test_Add(t *testing.T) {
	handler := handlerGenerator()

	request := mocks.GenerateCreatePaymentRequest()

	payments.On("Add", request).Return(request.CreateToEntity().ToDto(), nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment").
		WithMethod("POST").
		WithHandler(handler.Add()).
		WithBody(request)

	responseByteArray := testRequest.Verify(t, http.StatusCreated)

	actual := model.PaymentDto{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, model.PaymentDto{
		Id:          actual.Id,
		Name:        "Test Payment",
		Description: "Test Payment Description",
		HouseId:     mocks.HouseId,
		UserId:      mocks.UserId,
		ProviderId:  mocks.ProviderId,
		Date:        mocks.Date,
		Sum:         1000,
	}, actual)
}

func Test_Add_WithInvalidRequest(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment").
		WithMethod("POST").
		WithHandler(handler.Add())

	testRequest.Verify(t, http.StatusBadRequest)
}

func Test_Add_WithErrorFromService(t *testing.T) {
	handler := handlerGenerator()

	request := mocks.GenerateCreatePaymentRequest()

	expected := errors.New("error")

	payments.On("Add", request).Return(model.PaymentDto{}, expected)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment").
		WithMethod("POST").
		WithHandler(handler.Add()).
		WithBody(request)

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, fmt.Sprintf("%s\n", expected.Error()), string(responseByteArray))
}

func Test_Update(t *testing.T) {
	handler := handlerGenerator()

	request := mocks.GenerateUpdatePaymentRequest()
	id := uuid.New()

	payments.On("Update", id, request).Return(nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/{id}").
		WithMethod("PUT").
		WithHandler(handler.Update()).
		WithBody(request).
		WithVar("id", id.String())

	testRequest.Verify(t, http.StatusOK)
}

func Test_Update_WithInvalidId(t *testing.T) {
	handler := handlerGenerator()

	request := mocks.GenerateUpdatePaymentRequest()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/{id}").
		WithMethod("PUT").
		WithHandler(handler.Update()).
		WithBody(request).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "the id is not valid id\n", string(responseByteArray))

	payments.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func Test_Update_WithInvalidRequest(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/{id}").
		WithMethod("PUT").
		WithHandler(handler.Update()).
		WithVar("id", uuid.New().String())

	testRequest.Verify(t, http.StatusBadRequest)
}

func Test_Update_WithErrorFromService(t *testing.T) {
	handler := handlerGenerator()

	request := mocks.GenerateUpdatePaymentRequest()
	id := uuid.New()

	expected := errors.New("error")

	payments.On("Update", id, request).Return(expected)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/{id}").
		WithMethod("PUT").
		WithHandler(handler.Update()).
		WithVar("id", id.String()).
		WithBody(request)

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, fmt.Sprintf("%s\n", expected.Error()), string(responseByteArray))
}

func Test_Delete(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	payments.On("DeleteById", id).Return(nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/{id}").
		WithMethod("DELETE").
		WithHandler(handler.Delete()).
		WithVar("id", id.String())

	testRequest.Verify(t, http.StatusNoContent)
}

func Test_Delete_WithMissingParameter(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/{id}").
		WithMethod("DELETE").
		WithHandler(handler.Delete())

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "parameter 'id' not found\n", string(responseByteArray))
}

func Test_Delete_WithInvalidParameter(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/{id}").
		WithMethod("DELETE").
		WithHandler(handler.Delete()).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "the id is not valid id\n", string(responseByteArray))
}

func Test_FindById(t *testing.T) {
	handler := handlerGenerator()

	paymentResponse := mocks.GeneratePaymentResponse()

	payments.On("FindById", paymentResponse.Id).
		Return(paymentResponse, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", paymentResponse.Id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	actual := model.PaymentDto{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, paymentResponse, actual)
}

func Test_FindById_WithError(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	expected := errors.New("error")

	payments.On("FindById", id).
		Return(model.PaymentDto{}, expected)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, fmt.Sprintf("%s\n", expected.Error()), string(responseByteArray))
}

func Test_FindById_WithInvalidParameter(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "the id is not valid id\n", string(responseByteArray))
}

func Test_FindByHouseId(t *testing.T) {
	handler := handlerGenerator()

	response := mocks.GeneratePaymentResponse()

	paymentResponses := []model.PaymentDto{response}

	payments.On("FindByHouseId", response.Id).
		Return(paymentResponses, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/house/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByHouseId()).
		WithVar("id", response.Id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	actual := []model.PaymentDto{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, paymentResponses, actual)
}

func Test_FindByHouseId_WithEmptyResponse(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	paymentResponses := []model.PaymentDto{}

	payments.On("FindByHouseId", id).
		Return(paymentResponses, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/house/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByHouseId()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	actual := []model.PaymentDto{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, paymentResponses, actual)
}

func Test_FindByHouseId_WithInvalidParameter(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/house/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByHouseId()).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "the id is not valid id\n", string(responseByteArray))
}

func Test_FindByUserId(t *testing.T) {
	handler := handlerGenerator()

	response := mocks.GeneratePaymentResponse()

	paymentResponses := []model.PaymentDto{response}

	payments.On("FindByUserId", response.Id).
		Return(paymentResponses, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/user/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByUserId()).
		WithVar("id", response.Id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	var actual []model.PaymentDto

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, paymentResponses, actual)
}

func Test_FindByUserId_WithEmptyResponse(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	var paymentResponses []model.PaymentDto

	payments.On("FindByUserId", id).
		Return(paymentResponses, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/user/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByUserId()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	var actual []model.PaymentDto

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, paymentResponses, actual)
}

func Test_FindByUserId_WithInvalidParameter(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/user/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByUserId()).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "the id is not valid id\n", string(responseByteArray))
}

func Test_FindByProviderId(t *testing.T) {
	handler := handlerGenerator()

	response := mocks.GeneratePaymentResponse()

	paymentResponses := []model.PaymentDto{response}

	payments.On("FindByProviderId", response.Id).
		Return(paymentResponses, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/provider/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByProviderId()).
		WithVar("id", response.Id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	var actual []model.PaymentDto

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, paymentResponses, actual)
}

func Test_FindByProviderId_WithEmptyResponse(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	var paymentResponses []model.PaymentDto

	payments.On("FindByProviderId", id).
		Return(paymentResponses, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/provider/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByProviderId()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	var actual []model.PaymentDto

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, paymentResponses, actual)
}

func Test_FindByProviderId_WithInvalidParameter(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/provider/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByUserId()).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "the id is not valid id\n", string(responseByteArray))
}
