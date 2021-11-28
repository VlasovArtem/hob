package service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	houseMocks "house/mocks"
	meterMocks "meter/mocks"
	"meter/model"
	paymentMocks "payment/mocks"
	"testing"
)

var (
	payments        *paymentMocks.PaymentService
	houses          *houseMocks.HouseService
	meterRepository *meterMocks.MeterRepository
)

func serviceGenerator() MeterService {
	payments = new(paymentMocks.PaymentService)
	houses = new(houseMocks.HouseService)
	meterRepository = new(meterMocks.MeterRepository)

	return NewMeterService(payments, houses, meterRepository)
}

func Test_Add(t *testing.T) {
	service := serviceGenerator()

	var savedMeter model.Meter

	request := meterMocks.GenerateCreateMeterRequest()

	payments.On("ExistsById", request.PaymentId).Return(true)
	houses.On("ExistsById", request.HouseId).Return(true)
	meterRepository.On("Create", mock.Anything).Return(
		func(meter model.Meter) model.Meter {
			savedMeter = meter

			return meter
		},
		nil,
	)

	meter, err := service.Add(request)

	assert.Nil(t, err)
	assert.Equal(t, savedMeter.ToResponse(), meter)
}

func Test_Add_WithNotExistingPayment(t *testing.T) {
	service := serviceGenerator()
	request := meterMocks.GenerateCreateMeterRequest()

	payments.On("ExistsById", request.PaymentId).Return(false)

	meter, err := service.Add(request)

	assert.Equal(t, fmt.Sprintf("payment with id %s in not exists", request.PaymentId.String()), err.Error())
	assert.Equal(t, model.MeterResponse{}, meter)
}

func Test_Add_WithNotExistingHouse(t *testing.T) {
	service := serviceGenerator()
	request := meterMocks.GenerateCreateMeterRequest()

	payments.On("ExistsById", request.PaymentId).Return(true)
	houses.On("ExistsById", request.HouseId).Return(false)

	meter, err := service.Add(request)

	assert.Equal(t, fmt.Sprintf("house with id %s in not exists", request.HouseId.String()), err.Error())
	assert.Equal(t, model.MeterResponse{}, meter)
}

func Test_Add_WithErrorFromRepository(t *testing.T) {
	service := serviceGenerator()

	expectedError := errors.New("error")
	request := meterMocks.GenerateCreateMeterRequest()

	payments.On("ExistsById", request.PaymentId).Return(true)
	houses.On("ExistsById", request.HouseId).Return(true)
	meterRepository.On("Create", mock.Anything).Return(model.Meter{}, expectedError)

	meter, err := service.Add(request)

	assert.Equal(t, expectedError, err)
	assert.Equal(t, model.MeterResponse{}, meter)
}

func Test_FindById(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()

	expected := meterMocks.GenerateMeter(uuid.New(), uuid.New())

	meterRepository.On("FindById", id).Return(expected, nil)

	actual, err := service.FindById(id)

	assert.Nil(t, err)
	assert.Equal(t, expected.ToResponse(), actual)
}

func Test_FindById_WithMissingId(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()

	meterRepository.On("FindById", id).Return(model.Meter{}, gorm.ErrRecordNotFound)

	actual, err := service.FindById(id)

	assert.Equal(t, errors.New(fmt.Sprintf("meter with id %s in not exists", id)), err)
	assert.Equal(t, model.MeterResponse{}, actual)
}

func Test_FindById_WithError(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()
	expectedError := errors.New("error")

	meterRepository.On("FindById", id).Return(model.Meter{}, expectedError)

	actual, err := service.FindById(id)

	assert.Equal(t, expectedError, err)
	assert.Equal(t, model.MeterResponse{}, actual)
}

func Test_FindByPaymentId(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()

	expected := meterMocks.GenerateMeter(uuid.New(), uuid.New())

	meterRepository.On("FindByPaymentId", id).Return(expected, nil)

	actual, err := service.FindByPaymentId(id)

	assert.Nil(t, err)
	assert.Equal(t, expected.ToResponse(), actual)
}

func Test_FindByPaymentId_WithMissingId(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()

	meterRepository.On("FindByPaymentId", id).Return(model.Meter{}, gorm.ErrRecordNotFound)

	actual, err := service.FindByPaymentId(id)

	assert.Equal(t, errors.New(fmt.Sprintf("meter with payment id %s in not exists", id)), err)
	assert.Equal(t, model.MeterResponse{}, actual)
}

func Test_FindByPaymentId_WithError(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()
	expectedError := errors.New("error")

	meterRepository.On("FindByPaymentId", id).Return(model.Meter{}, expectedError)

	actual, err := service.FindByPaymentId(id)

	assert.Equal(t, expectedError, err)
	assert.Equal(t, model.MeterResponse{}, actual)
}

func Test_FindByHouseId(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()

	expected := meterMocks.GenerateMeter(uuid.New(), uuid.New())

	meterRepository.On("FindByHouseId", id).Return([]model.Meter{expected})

	actual := service.FindByHouseId(id)

	assert.Equal(t, []model.MeterResponse{expected.ToResponse()}, actual)
}

func Test_FindByHouseId_WithEmptyResponse(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()

	meterRepository.On("FindByHouseId", id).Return([]model.Meter{})

	actual := service.FindByHouseId(id)

	assert.Equal(t, []model.MeterResponse{}, actual)
}