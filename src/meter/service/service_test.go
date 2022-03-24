package service

import (
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/common/int-errors"
	meterMocks "github.com/VlasovArtem/hob/src/meter/mocks"
	"github.com/VlasovArtem/hob/src/meter/model"
	paymentMocks "github.com/VlasovArtem/hob/src/payment/mocks"
	"github.com/VlasovArtem/hob/src/test/testhelper"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type MeterServiceTestSuite struct {
	testhelper.MockTestSuite[MeterService]
	payments        *paymentMocks.PaymentService
	meterRepository *meterMocks.MeterRepository
}

func TestMeterServiceTestSuite(t *testing.T) {
	ts := &MeterServiceTestSuite{}
	ts.TestObjectGenerator = func() MeterService {
		ts.payments = new(paymentMocks.PaymentService)
		ts.meterRepository = new(meterMocks.MeterRepository)
		return NewMeterService(ts.payments, ts.meterRepository)
	}

	suite.Run(t, ts)
}

func (m *MeterServiceTestSuite) Test_Add() {
	var savedMeter model.Meter

	request := meterMocks.GenerateCreateMeterRequest()

	m.payments.On("ExistsById", request.PaymentId).Return(true)
	m.meterRepository.On("Create", mock.Anything).Return(
		func(meter model.Meter) model.Meter {
			savedMeter = meter

			return meter
		},
		nil,
	)

	meter, err := m.TestO.Add(request)

	assert.Nil(m.T(), err)
	assert.Equal(m.T(), savedMeter.ToDto(), meter)
}

func (m *MeterServiceTestSuite) Test_Add_WithNotExistingPayment() {
	request := meterMocks.GenerateCreateMeterRequest()

	m.payments.On("ExistsById", request.PaymentId).Return(false)

	meter, err := m.TestO.Add(request)

	assert.Equal(m.T(), fmt.Sprintf("payment with id %s not found", request.PaymentId.String()), err.Error())
	assert.Equal(m.T(), model.MeterDto{}, meter)
}

func (m *MeterServiceTestSuite) Test_Add_WithErrorFromRepository() {
	expectedError := errors.New("error")
	request := meterMocks.GenerateCreateMeterRequest()

	m.payments.On("ExistsById", request.PaymentId).Return(true)
	m.meterRepository.On("Create", mock.Anything).Return(model.Meter{}, expectedError)

	meter, err := m.TestO.Add(request)

	assert.Equal(m.T(), expectedError, err)
	assert.Equal(m.T(), model.MeterDto{}, meter)
}

func (m *MeterServiceTestSuite) Test_Update() {
	id, request := meterMocks.GenerateUpdateMeterRequest()

	m.meterRepository.On("ExistsById", id).Return(true)
	m.meterRepository.On("Update", id, request.ToEntity()).Return(nil)

	err := m.TestO.Update(id, request)

	assert.Nil(m.T(), err)
}

func (m *MeterServiceTestSuite) Test_Update_WithMissingId() {
	id, request := meterMocks.GenerateUpdateMeterRequest()

	m.meterRepository.On("ExistsById", id).Return(false)

	err := m.TestO.Update(id, request)

	assert.Equal(m.T(), int_errors.NewErrNotFound("meter with id %s not found", id), err)

	m.meterRepository.AssertNotCalled(m.T(), "Update", id, request)
}

func (m *MeterServiceTestSuite) Test_Update_WithErrorFromRepository() {
	id, request := meterMocks.GenerateUpdateMeterRequest()

	m.meterRepository.On("ExistsById", id).Return(true)
	m.meterRepository.On("Update", id, request.ToEntity()).Return(errors.New("test"))

	err := m.TestO.Update(id, request)

	assert.Equal(m.T(), errors.New("test"), err)
}

func (m *MeterServiceTestSuite) Test_DeleteById() {
	id := uuid.New()

	m.meterRepository.On("ExistsById", id).Return(true)
	m.meterRepository.On("DeleteById", id).Return(nil)

	err := m.TestO.DeleteById(id)

	assert.Nil(m.T(), err)
}

func (m *MeterServiceTestSuite) Test_DeleteById_WithMissingId() {
	id := uuid.New()

	m.meterRepository.On("ExistsById", id).Return(false)

	err := m.TestO.DeleteById(id)

	assert.Equal(m.T(), int_errors.NewErrNotFound("meter with id %s not found", id), err)

	m.meterRepository.AssertNotCalled(m.T(), "DeleteById", id)
}

func (m *MeterServiceTestSuite) Test_DeleteById_WithErrorFromRepository() {
	id := uuid.New()

	m.meterRepository.On("ExistsById", id).Return(true)
	m.meterRepository.On("DeleteById", id).Return(errors.New("test"))

	err := m.TestO.DeleteById(id)

	assert.Equal(m.T(), errors.New("test"), err)
}

func (m *MeterServiceTestSuite) Test_FindById() {
	id := uuid.New()

	expected := meterMocks.GenerateMeter(uuid.New())

	m.meterRepository.On("FindById", id).Return(expected, nil)

	actual, err := m.TestO.FindById(id)

	assert.Nil(m.T(), err)
	assert.Equal(m.T(), expected.ToDto(), actual)
}

func (m *MeterServiceTestSuite) Test_FindById_WithMissingId() {
	id := uuid.New()

	m.meterRepository.On("FindById", id).Return(model.Meter{}, gorm.ErrRecordNotFound)

	actual, err := m.TestO.FindById(id)

	assert.Equal(m.T(), int_errors.NewErrNotFound("meter with id %s in not exists", id), err)
	assert.Equal(m.T(), model.MeterDto{}, actual)
}

func (m *MeterServiceTestSuite) Test_FindById_WithError() {
	id := uuid.New()
	expectedError := errors.New("error")

	m.meterRepository.On("FindById", id).Return(model.Meter{}, expectedError)

	actual, err := m.TestO.FindById(id)

	assert.Equal(m.T(), expectedError, err)
	assert.Equal(m.T(), model.MeterDto{}, actual)
}

func (m *MeterServiceTestSuite) Test_FindByPaymentId() {
	id := uuid.New()

	expected := meterMocks.GenerateMeter(uuid.New())

	m.meterRepository.On("FindByPaymentId", id).Return(expected, nil)

	actual, err := m.TestO.FindByPaymentId(id)

	assert.Nil(m.T(), err)
	assert.Equal(m.T(), expected.ToDto(), actual)
}

func (m *MeterServiceTestSuite) Test_FindByPaymentId_WithMissingId() {
	id := uuid.New()

	m.meterRepository.On("FindByPaymentId", id).Return(model.Meter{}, gorm.ErrRecordNotFound)

	actual, err := m.TestO.FindByPaymentId(id)

	assert.Equal(m.T(), int_errors.NewErrNotFound("meter with payment id %s in not exists", id), err)
	assert.Equal(m.T(), model.MeterDto{}, actual)
}

func (m *MeterServiceTestSuite) Test_FindByPaymentId_WithError() {
	id := uuid.New()
	expectedError := errors.New("error")

	m.meterRepository.On("FindByPaymentId", id).Return(model.Meter{}, expectedError)

	actual, err := m.TestO.FindByPaymentId(id)

	assert.Equal(m.T(), expectedError, err)
	assert.Equal(m.T(), model.MeterDto{}, actual)
}
