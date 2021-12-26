package service

import (
	"errors"
	"fmt"
	int_errors "github.com/VlasovArtem/hob/src/common/int-errors"
	houseMocks "github.com/VlasovArtem/hob/src/house/mocks"
	meterMocks "github.com/VlasovArtem/hob/src/meter/mocks"
	"github.com/VlasovArtem/hob/src/meter/model"
	paymentMocks "github.com/VlasovArtem/hob/src/payment/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
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
	assert.Equal(t, savedMeter.ToDto(), meter)
}

func Test_Add_WithNotExistingPayment(t *testing.T) {
	service := serviceGenerator()
	request := meterMocks.GenerateCreateMeterRequest()

	payments.On("ExistsById", request.PaymentId).Return(false)

	meter, err := service.Add(request)

	assert.Equal(t, fmt.Sprintf("payment with id %s not found", request.PaymentId.String()), err.Error())
	assert.Equal(t, model.MeterDto{}, meter)
}

func Test_Add_WithNotExistingHouse(t *testing.T) {
	service := serviceGenerator()
	request := meterMocks.GenerateCreateMeterRequest()

	payments.On("ExistsById", request.PaymentId).Return(true)
	houses.On("ExistsById", request.HouseId).Return(false)

	meter, err := service.Add(request)

	assert.Equal(t, fmt.Sprintf("house with id %s not found", request.HouseId.String()), err.Error())
	assert.Equal(t, model.MeterDto{}, meter)
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
	assert.Equal(t, model.MeterDto{}, meter)
}

func Test_Update(t *testing.T) {
	service := serviceGenerator()

	id, request := meterMocks.GenerateUpdateMeterRequest()

	meterRepository.On("ExistsById", id).Return(true)
	meterRepository.On("Update", id, request.ToEntity()).Return(nil)

	err := service.Update(id, request)

	assert.Nil(t, err)
}

func Test_Update_WithMissingId(t *testing.T) {
	service := serviceGenerator()

	id, request := meterMocks.GenerateUpdateMeterRequest()

	meterRepository.On("ExistsById", id).Return(false)

	err := service.Update(id, request)

	assert.Equal(t, int_errors.NewErrNotFound("meter with id %s not found", id), err)

	meterRepository.AssertNotCalled(t, "Update", id, request)
}

func Test_Update_WithErrorFromRepository(t *testing.T) {
	service := serviceGenerator()

	id, request := meterMocks.GenerateUpdateMeterRequest()

	meterRepository.On("ExistsById", id).Return(true)
	meterRepository.On("Update", id, request.ToEntity()).Return(errors.New("test"))

	err := service.Update(id, request)

	assert.Equal(t, errors.New("test"), err)
}

func Test_DeleteById(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()

	meterRepository.On("ExistsById", id).Return(true)
	meterRepository.On("DeleteById", id).Return(nil)

	err := service.DeleteById(id)

	assert.Nil(t, err)
}

func Test_DeleteById_WithMissingId(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()

	meterRepository.On("ExistsById", id).Return(false)

	err := service.DeleteById(id)

	assert.Equal(t, int_errors.NewErrNotFound("meter with id %s not found", id), err)

	meterRepository.AssertNotCalled(t, "DeleteById", id)
}

func Test_DeleteById_WithErrorFromRepository(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()

	meterRepository.On("ExistsById", id).Return(true)
	meterRepository.On("DeleteById", id).Return(errors.New("test"))

	err := service.DeleteById(id)

	assert.Equal(t, errors.New("test"), err)
}

func Test_FindById(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()

	expected := meterMocks.GenerateMeter(uuid.New(), uuid.New())

	meterRepository.On("FindById", id).Return(expected, nil)

	actual, err := service.FindById(id)

	assert.Nil(t, err)
	assert.Equal(t, expected.ToDto(), actual)
}

func Test_FindById_WithMissingId(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()

	meterRepository.On("FindById", id).Return(model.Meter{}, gorm.ErrRecordNotFound)

	actual, err := service.FindById(id)

	assert.Equal(t, fmt.Errorf("meter with id %s in not exists", id), err)
	assert.Equal(t, model.MeterDto{}, actual)
}

func Test_FindById_WithError(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()
	expectedError := errors.New("error")

	meterRepository.On("FindById", id).Return(model.Meter{}, expectedError)

	actual, err := service.FindById(id)

	assert.Equal(t, expectedError, err)
	assert.Equal(t, model.MeterDto{}, actual)
}

func Test_FindByPaymentId(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()

	expected := meterMocks.GenerateMeter(uuid.New(), uuid.New())

	meterRepository.On("FindByPaymentId", id).Return(expected, nil)

	actual, err := service.FindByPaymentId(id)

	assert.Nil(t, err)
	assert.Equal(t, expected.ToDto(), actual)
}

func Test_FindByPaymentId_WithMissingId(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()

	meterRepository.On("FindByPaymentId", id).Return(model.Meter{}, gorm.ErrRecordNotFound)

	actual, err := service.FindByPaymentId(id)

	assert.Equal(t, fmt.Errorf("meter with payment id %s in not exists", id), err)
	assert.Equal(t, model.MeterDto{}, actual)
}

func Test_FindByPaymentId_WithError(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()
	expectedError := errors.New("error")

	meterRepository.On("FindByPaymentId", id).Return(model.Meter{}, expectedError)

	actual, err := service.FindByPaymentId(id)

	assert.Equal(t, expectedError, err)
	assert.Equal(t, model.MeterDto{}, actual)
}

func Test_FindByHouseId(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()

	expected := meterMocks.GenerateMeter(uuid.New(), uuid.New())

	meterRepository.On("FindByHouseId", id).Return([]model.Meter{expected})

	actual := service.FindByHouseId(id)

	assert.Equal(t, []model.MeterDto{expected.ToDto()}, actual)
}

func Test_FindByHouseId_WithEmptyResponse(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()

	meterRepository.On("FindByHouseId", id).Return([]model.Meter{})

	actual := service.FindByHouseId(id)

	assert.Equal(t, []model.MeterDto{}, actual)
}
