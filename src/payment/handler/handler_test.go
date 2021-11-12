package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"payment/model"
	"test/mock"
	"test/testhelper"
	"testing"
	"time"
)

var (
	payments *mock.PaymentServiceMock
	houseId  = testhelper.ParseUUID("d077adaa-00d7-4e80-ac86-57512267505d")
	userId   = testhelper.ParseUUID("ad2c5035-6745-48d0-9eee-fd22f5dae8e0")
	date     = time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
)

func handlerGenerator() PaymentHandler {
	payments = new(mock.PaymentServiceMock)

	return NewPaymentHandler(payments)
}

func Test_AddPayment(t *testing.T) {
	handler := handlerGenerator()

	request := generateCreatePaymentRequest()

	payments.On("AddPayment", request).Return(request.ToEntity().ToResponse(), nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment").
		WithMethod("POST").
		WithHandler(handler.AddPayment()).
		WithBody(request)

	responseByteArray := testRequest.Verify(t, http.StatusCreated)

	actual := model.PaymentResponse{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, model.PaymentResponse{
		Payment: model.Payment{
			Id:          actual.Id,
			Name:        "Test Payment",
			Description: "Test Payment Description",
			HouseId:     houseId,
			UserId:      userId,
			Date:        date,
			Sum:         1000,
		},
	}, actual)
}

func Test_AddPayment_WithInvalidRequest(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment").
		WithMethod("POST").
		WithHandler(handler.AddPayment())

	testRequest.Verify(t, http.StatusBadRequest)
}

func Test_AddPayment_WithErrorFromService(t *testing.T) {
	handler := handlerGenerator()

	request := generateCreatePaymentRequest()

	expected := errors.New("error")

	payments.On("AddPayment", request).Return(model.PaymentResponse{}, expected)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment").
		WithMethod("POST").
		WithHandler(handler.AddPayment()).
		WithBody(request)

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, fmt.Sprintf("%s\n", expected.Error()), string(responseByteArray))
}

func Test_FindPaymentById(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	paymentResponse := generatePaymentResponse(id)

	payments.On("FindPaymentById", id).
		Return(paymentResponse, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/{id}").
		WithMethod("GET").
		WithHandler(handler.FindPaymentById()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	actual := model.PaymentResponse{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, paymentResponse, actual)
}

func Test_FindPaymentById_WithError(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	expected := errors.New("error")

	payments.On("FindPaymentById", id).
		Return(model.PaymentResponse{}, expected)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/{id}").
		WithMethod("GET").
		WithHandler(handler.FindPaymentById()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, fmt.Sprintf("%s\n", expected.Error()), string(responseByteArray))
}

func Test_FindPaymentById_WithInvalidParameter(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/{id}").
		WithMethod("GET").
		WithHandler(handler.FindPaymentById()).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "the id is not valid id\n", string(responseByteArray))
}

func Test_FindPaymentByHouseId(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	paymentResponses := []model.PaymentResponse{generatePaymentResponse(id)}

	payments.On("FindPaymentByHouseId", id).
		Return(paymentResponses, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/{id}").
		WithMethod("GET").
		WithHandler(handler.FindPaymentByHouseId()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	actual := []model.PaymentResponse{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, paymentResponses, actual)
}

func Test_FindPaymentByHouseId_WithEmptyResponse(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	paymentResponses := []model.PaymentResponse{}

	payments.On("FindPaymentByHouseId", id).
		Return(paymentResponses, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/{id}").
		WithMethod("GET").
		WithHandler(handler.FindPaymentByHouseId()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	actual := []model.PaymentResponse{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, paymentResponses, actual)
}

func Test_FindPaymentByHouseId_WithInvalidParameter(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/{id}").
		WithMethod("GET").
		WithHandler(handler.FindPaymentByHouseId()).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "the id is not valid id\n", string(responseByteArray))
}

func Test_FindPaymentByUserId(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	paymentResponses := []model.PaymentResponse{generatePaymentResponse(id)}

	payments.On("FindPaymentByUserId", id).
		Return(paymentResponses, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/{id}").
		WithMethod("GET").
		WithHandler(handler.FindPaymentByUserId()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	actual := []model.PaymentResponse{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, paymentResponses, actual)
}

func Test_FindPaymentByUserId_WithEmptyResponse(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	paymentResponses := []model.PaymentResponse{}

	payments.On("FindPaymentByUserId", id).
		Return(paymentResponses, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/{id}").
		WithMethod("GET").
		WithHandler(handler.FindPaymentByUserId()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	actual := []model.PaymentResponse{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, paymentResponses, actual)
}

func Test_FindPaymentByUserId_WithInvalidParameter(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/{id}").
		WithMethod("GET").
		WithHandler(handler.FindPaymentByUserId()).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "the id is not valid id\n", string(responseByteArray))
}

func generateCreatePaymentRequest() model.CreatePaymentRequest {
	return model.CreatePaymentRequest{
		Name:        "Test Payment",
		Description: "Test Payment Description",
		HouseId:     houseId,
		UserId:      userId,
		Date:        date,
		Sum:         1000,
	}
}

func generatePaymentResponse(id uuid.UUID) model.PaymentResponse {
	return model.PaymentResponse{
		Payment: model.Payment{
			Id:          id,
			Name:        "Test Payment",
			Description: "Test Payment Description",
			HouseId:     houseId,
			UserId:      userId,
			Date:        date,
			Sum:         1000,
		},
	}
}
