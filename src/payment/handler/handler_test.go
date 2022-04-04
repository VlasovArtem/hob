package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/common"
	int_errors "github.com/VlasovArtem/hob/src/common/int-errors"
	"github.com/VlasovArtem/hob/src/payment/mocks"
	"github.com/VlasovArtem/hob/src/payment/model"
	"github.com/VlasovArtem/hob/src/test/testhelper"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type PaymentHandlerTestSuite struct {
	testhelper.MockTestSuite[PaymentHandler]
	payments *mocks.PaymentService
}

func TestPaymentHandlerTestSuite(t *testing.T) {
	testingSuite := &PaymentHandlerTestSuite{}
	testingSuite.TestObjectGenerator = func() PaymentHandler {
		testingSuite.payments = new(mocks.PaymentService)
		return NewPaymentHandler(testingSuite.payments)
	}

	suite.Run(t, testingSuite)
}

func (p *PaymentHandlerTestSuite) Test_Add() {
	request := mocks.GenerateCreatePaymentRequest()

	p.payments.On("Add", request).Return(request.ToEntity().ToDto(), nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment").
		WithMethod("POST").
		WithHandler(p.TestO.Add()).
		WithBody(request)

	responseByteArray := testRequest.Verify(p.T(), http.StatusCreated)

	actual := model.PaymentDto{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(p.T(), model.PaymentDto{
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

func (p *PaymentHandlerTestSuite) Test_Add_WithInvalidRequest() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment").
		WithMethod("POST").
		WithHandler(p.TestO.Add())

	testRequest.Verify(p.T(), http.StatusBadRequest)
}

func (p *PaymentHandlerTestSuite) Test_Add_WithErrorFromService() {
	request := mocks.GenerateCreatePaymentRequest()

	expected := errors.New("error")

	p.payments.On("Add", request).Return(model.PaymentDto{}, expected)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment").
		WithMethod("POST").
		WithHandler(p.TestO.Add()).
		WithBody(request)

	responseByteArray := testRequest.Verify(p.T(), http.StatusBadRequest)

	assert.Equal(p.T(), fmt.Sprintf("%s\n", expected.Error()), string(responseByteArray))
}

func (p *PaymentHandlerTestSuite) Test_AddBatch() {
	request := mocks.GenerateCreatePaymentBatchRequest(1)

	p.payments.On("AddBatch", request).Return(common.MapSlice(request.Payments, func(request model.CreatePaymentRequest) model.PaymentDto {
		return request.ToEntity().ToDto()
	}), nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/batch").
		WithMethod("POST").
		WithHandler(p.TestO.AddBatch()).
		WithBody(request)

	responseByteArray := testRequest.Verify(p.T(), http.StatusCreated)

	var actual []model.PaymentDto

	err := json.Unmarshal(responseByteArray, &actual)

	assert.Nil(p.T(), err)
	assert.Equal(p.T(), []model.PaymentDto{
		{
			Id:          actual[0].Id,
			Name:        "Payment Name #0",
			Description: "Test Payment Description",
			HouseId:     mocks.HouseId,
			UserId:      mocks.UserId,
			ProviderId:  mocks.ProviderId,
			Date:        mocks.Date,
			Sum:         1000,
		},
	}, actual)
}

func (p *PaymentHandlerTestSuite) Test_AddBatch_WithInvalidRequest() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/batch").
		WithMethod("POST").
		WithHandler(p.TestO.AddBatch())

	testRequest.Verify(p.T(), http.StatusBadRequest)
}

func (p *PaymentHandlerTestSuite) Test_AddBatch_WithErrorFromService() {
	request := mocks.GenerateCreatePaymentBatchRequest(1)

	errorResponse := int_errors.NewBuilder().
		WithDetail("message").
		WithMessage("error")
	err := int_errors.NewErrResponse(errorResponse)

	p.payments.On("AddBatch", request).Return([]model.PaymentDto{}, err)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/batch").
		WithMethod("POST").
		WithHandler(p.TestO.AddBatch()).
		WithBody(request)

	actual := testRequest.Verify(p.T(), http.StatusBadRequest)

	expected, err := json.Marshal(errorResponse.Build())

	assert.Nil(p.T(), err)
	assert.Equal(p.T(), append(expected, []byte("\n")...), actual)
}

func (p *PaymentHandlerTestSuite) Test_Update() {
	request := mocks.GenerateUpdatePaymentRequest()
	id := uuid.New()

	p.payments.On("Update", id, request).Return(nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/{id}").
		WithMethod("PUT").
		WithHandler(p.TestO.Update()).
		WithBody(request).
		WithVar("id", id.String())

	testRequest.Verify(p.T(), http.StatusOK)
}

func (p *PaymentHandlerTestSuite) Test_Update_WithInvalidId() {
	request := mocks.GenerateUpdatePaymentRequest()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/{id}").
		WithMethod("PUT").
		WithHandler(p.TestO.Update()).
		WithBody(request).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(p.T(), http.StatusBadRequest)

	assert.Equal(p.T(), "the id is not valid id\n", string(responseByteArray))

	p.payments.AssertNotCalled(p.T(), "Update", mock.Anything, mock.Anything)
}

func (p *PaymentHandlerTestSuite) Test_Update_WithInvalidRequest() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/{id}").
		WithMethod("PUT").
		WithHandler(p.TestO.Update()).
		WithVar("id", uuid.New().String())

	testRequest.Verify(p.T(), http.StatusBadRequest)
}

func (p *PaymentHandlerTestSuite) Test_Update_WithErrorFromService() {
	request := mocks.GenerateUpdatePaymentRequest()
	id := uuid.New()

	expected := errors.New("error")

	p.payments.On("Update", id, request).Return(expected)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/{id}").
		WithMethod("PUT").
		WithHandler(p.TestO.Update()).
		WithVar("id", id.String()).
		WithBody(request)

	responseByteArray := testRequest.Verify(p.T(), http.StatusBadRequest)

	assert.Equal(p.T(), fmt.Sprintf("%s\n", expected.Error()), string(responseByteArray))
}

func (p *PaymentHandlerTestSuite) Test_Delete() {
	id := uuid.New()

	p.payments.On("DeleteById", id).Return(nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/{id}").
		WithMethod("DELETE").
		WithHandler(p.TestO.Delete()).
		WithVar("id", id.String())

	testRequest.Verify(p.T(), http.StatusNoContent)
}

func (p *PaymentHandlerTestSuite) Test_Delete_WithMissingParameter() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/{id}").
		WithMethod("DELETE").
		WithHandler(p.TestO.Delete())

	responseByteArray := testRequest.Verify(p.T(), http.StatusBadRequest)

	assert.Equal(p.T(), "parameter 'id' not found\n", string(responseByteArray))
}

func (p *PaymentHandlerTestSuite) Test_Delete_WithInvalidParameter() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/{id}").
		WithMethod("DELETE").
		WithHandler(p.TestO.Delete()).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(p.T(), http.StatusBadRequest)

	assert.Equal(p.T(), "the id is not valid id\n", string(responseByteArray))
}

func (p *PaymentHandlerTestSuite) Test_FindById() {
	paymentResponse := mocks.GeneratePaymentResponse()

	p.payments.On("FindById", paymentResponse.Id).
		Return(paymentResponse, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/{id}").
		WithMethod("GET").
		WithHandler(p.TestO.FindById()).
		WithVar("id", paymentResponse.Id.String())

	responseByteArray := testRequest.Verify(p.T(), http.StatusOK)

	actual := model.PaymentDto{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(p.T(), paymentResponse, actual)
}

func (p *PaymentHandlerTestSuite) Test_FindById_WithError() {
	id := uuid.New()

	expected := errors.New("error")

	p.payments.On("FindById", id).
		Return(model.PaymentDto{}, expected)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/{id}").
		WithMethod("GET").
		WithHandler(p.TestO.FindById()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(p.T(), http.StatusBadRequest)

	assert.Equal(p.T(), fmt.Sprintf("%s\n", expected.Error()), string(responseByteArray))
}

func (p *PaymentHandlerTestSuite) Test_FindById_WithInvalidParameter() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/{id}").
		WithMethod("GET").
		WithHandler(p.TestO.FindById()).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(p.T(), http.StatusBadRequest)

	assert.Equal(p.T(), "the id is not valid id\n", string(responseByteArray))
}

func (p *PaymentHandlerTestSuite) Test_FindByHouseId() {
	response := mocks.GeneratePaymentResponse()

	paymentResponses := []model.PaymentDto{response}

	p.payments.On("FindByHouseId", response.Id).
		Return(paymentResponses, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/house/{id}").
		WithMethod("GET").
		WithHandler(p.TestO.FindByHouseId()).
		WithVar("id", response.Id.String())

	responseByteArray := testRequest.Verify(p.T(), http.StatusOK)

	actual := []model.PaymentDto{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(p.T(), paymentResponses, actual)
}

func (p *PaymentHandlerTestSuite) Test_FindByHouseId_WithEmptyResponse() {
	id := uuid.New()

	paymentResponses := []model.PaymentDto{}

	p.payments.On("FindByHouseId", id).
		Return(paymentResponses, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/house/{id}").
		WithMethod("GET").
		WithHandler(p.TestO.FindByHouseId()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(p.T(), http.StatusOK)

	actual := []model.PaymentDto{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(p.T(), paymentResponses, actual)
}

func (p *PaymentHandlerTestSuite) Test_FindByHouseId_WithInvalidParameter() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/house/{id}").
		WithMethod("GET").
		WithHandler(p.TestO.FindByHouseId()).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(p.T(), http.StatusBadRequest)

	assert.Equal(p.T(), "the id is not valid id\n", string(responseByteArray))
}

func (p *PaymentHandlerTestSuite) Test_FindByUserId() {
	response := mocks.GeneratePaymentResponse()

	paymentResponses := []model.PaymentDto{response}

	p.payments.On("FindByUserId", response.Id).
		Return(paymentResponses, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/user/{id}").
		WithMethod("GET").
		WithHandler(p.TestO.FindByUserId()).
		WithVar("id", response.Id.String())

	responseByteArray := testRequest.Verify(p.T(), http.StatusOK)

	var actual []model.PaymentDto

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(p.T(), paymentResponses, actual)
}

func (p *PaymentHandlerTestSuite) Test_FindByUserId_WithEmptyResponse() {
	id := uuid.New()

	var paymentResponses []model.PaymentDto

	p.payments.On("FindByUserId", id).
		Return(paymentResponses, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/user/{id}").
		WithMethod("GET").
		WithHandler(p.TestO.FindByUserId()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(p.T(), http.StatusOK)

	var actual []model.PaymentDto

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(p.T(), paymentResponses, actual)
}

func (p *PaymentHandlerTestSuite) Test_FindByUserId_WithInvalidParameter() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/user/{id}").
		WithMethod("GET").
		WithHandler(p.TestO.FindByUserId()).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(p.T(), http.StatusBadRequest)

	assert.Equal(p.T(), "the id is not valid id\n", string(responseByteArray))
}

func (p *PaymentHandlerTestSuite) Test_FindByProviderId() {
	response := mocks.GeneratePaymentResponse()

	paymentResponses := []model.PaymentDto{response}

	p.payments.On("FindByProviderId", response.Id).
		Return(paymentResponses, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/provider/{id}").
		WithMethod("GET").
		WithHandler(p.TestO.FindByProviderId()).
		WithVar("id", response.Id.String())

	responseByteArray := testRequest.Verify(p.T(), http.StatusOK)

	var actual []model.PaymentDto

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(p.T(), paymentResponses, actual)
}

func (p *PaymentHandlerTestSuite) Test_FindByProviderId_WithEmptyResponse() {
	id := uuid.New()

	var paymentResponses []model.PaymentDto

	p.payments.On("FindByProviderId", id).
		Return(paymentResponses, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/provider/{id}").
		WithMethod("GET").
		WithHandler(p.TestO.FindByProviderId()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(p.T(), http.StatusOK)

	var actual []model.PaymentDto

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(p.T(), paymentResponses, actual)
}

func (p *PaymentHandlerTestSuite) Test_FindByProviderId_WithInvalidParameter() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/payment/provider/{id}").
		WithMethod("GET").
		WithHandler(p.TestO.FindByUserId()).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(p.T(), http.StatusBadRequest)

	assert.Equal(p.T(), "the id is not valid id\n", string(responseByteArray))
}
