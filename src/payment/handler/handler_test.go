package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"payment/mocks"
	"payment/model"
	"test/testhelper"
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

	payments.On("Add", request).Return(request.ToEntity().ToDto(), nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment").
		WithMethod("POST").
		WithHandler(handler.Add()).
		WithBody(request)

	responseByteArray := testRequest.Verify(t, http.StatusCreated)

	actual := model.PaymentDto{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, model.PaymentDto{
		Payment: model.Payment{
			Id:          actual.Id,
			Name:        "Test Payment",
			Description: "Test Payment Description",
			HouseId:     mocks.HouseId,
			UserId:      mocks.UserId,
			Date:        mocks.Date,
			Sum:         1000,
		},
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

	actual := []model.PaymentDto{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, paymentResponses, actual)
}

func Test_FindByUserId_WithEmptyResponse(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	paymentResponses := []model.PaymentDto{}

	payments.On("FindByUserId", id).
		Return(paymentResponses, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/user/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByUserId()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	actual := []model.PaymentDto{}

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
