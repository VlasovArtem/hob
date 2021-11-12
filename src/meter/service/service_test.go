package service

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"meter/model"
	"test/mock"
	"test/testhelper"
	"testing"
)

var (
	payments  *mock.PaymentServiceMock
	paymentId = testhelper.ParseUUID("8db24b37-6978-4c7f-ae8d-516eaabb323b")
)

func serviceGenerator() MeterService {
	payments = new(mock.PaymentServiceMock)

	return NewMeterService(payments)
}

func Test_AddMeter(t *testing.T) {
	service := serviceGenerator()

	payments.On("ExistsById", paymentId).Return(true)

	request := model.CreateMeterRequest{
		Name: "Name",
		Details: map[string]float64{
			"first":  1.1,
			"second": 2.2,
		},
		Description: "Description",
		PaymentId:   paymentId,
	}

	meter, err := service.AddMeter(request)

	expectedResponse := request.ToEntity().ToResponse()
	expectedResponse.Id = meter.Id

	assert.Nil(t, err)
	assert.Equal(t, expectedResponse, meter)
}

func Test_AddMeter_WithNotExistingPayment(t *testing.T) {
	service := serviceGenerator()

	payments.On("ExistsById", paymentId).Return(false)

	request := model.CreateMeterRequest{
		Name: "Name",
		Details: map[string]float64{
			"first":  1.1,
			"second": 2.2,
		},
		Description: "Description",
		PaymentId:   paymentId,
	}

	meter, err := service.AddMeter(request)

	assert.Equal(t, fmt.Sprintf("payment with id %s in not exists", paymentId.String()), err.Error())
	assert.Equal(t, model.MeterResponse{}, meter)
}

func Test_FindById(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()

	expected := model.Meter{Name: "test"}

	service.(*meterServiceObject).meters[id] = expected

	actual, err := service.FindById(id)

	assert.Nil(t, err)
	assert.Equal(t, expected.ToResponse(), actual)
}

func Test_FindById_WithMissingId(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()

	actual, err := service.FindById(id)

	assert.Equal(t, fmt.Sprintf("meter with id %s in not exists", id), err.Error())
	assert.Equal(t, model.MeterResponse{}, actual)
}

func Test_FindByPaymentId(t *testing.T) {
	service := serviceGenerator()

	expected := model.Meter{Name: "test2"}

	service.(*meterServiceObject).paymentMeter[paymentId] = expected

	actual, err := service.FindByPaymentId(paymentId)

	assert.Nil(t, err)
	assert.Equal(t, expected.ToResponse(), actual)
}

func Test_FindByPaymentId_WithMissing(t *testing.T) {
	service := serviceGenerator()

	actual, err := service.FindByPaymentId(paymentId)

	assert.Equal(t, fmt.Sprintf("meters with payment id %s not found", paymentId), err.Error())
	assert.Equal(t, model.MeterResponse{}, actual)
}
