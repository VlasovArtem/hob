package service

import (
	"errors"
	"fmt"
	interrors "github.com/VlasovArtem/hob/src/common/int-errors"
	houseMocks "github.com/VlasovArtem/hob/src/house/mocks"
	"github.com/VlasovArtem/hob/src/payment/mocks"
	"github.com/VlasovArtem/hob/src/payment/model"
	providerMocks "github.com/VlasovArtem/hob/src/provider/mocks"
	"github.com/VlasovArtem/hob/src/test/testhelper"
	userMocks "github.com/VlasovArtem/hob/src/user/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
	"time"
)

type PaymentServiceTestSuite struct {
	testhelper.MockTestSuite[PaymentService]
	userService       *userMocks.UserService
	houseService      *houseMocks.HouseService
	providerService   *providerMocks.ProviderService
	paymentRepository *mocks.PaymentRepository
}

func TestIncomeServiceTestSuite(t *testing.T) {
	ts := &PaymentServiceTestSuite{}
	ts.TestObjectGenerator = func() PaymentService {
		ts.userService = new(userMocks.UserService)
		ts.houseService = new(houseMocks.HouseService)
		ts.providerService = new(providerMocks.ProviderService)
		ts.paymentRepository = new(mocks.PaymentRepository)

		return NewPaymentService(ts.userService, ts.houseService, ts.providerService, ts.paymentRepository)
	}

	suite.Run(t, ts)
}

func (p *PaymentServiceTestSuite) Test_Add() {
	p.userService.On("ExistsById", mocks.UserId).Return(true)
	p.houseService.On("ExistsById", mocks.HouseId).Return(true)
	p.providerService.On("ExistsById", mocks.ProviderId).Return(true)
	p.paymentRepository.On("Create", mock.Anything).Return(
		func(payment model.Payment) model.Payment { return payment },
		nil,
	)

	request := mocks.GenerateCreatePaymentRequest()

	payment, err := p.TestO.Add(request)

	expectedEntity := request.CreateToEntity()
	expectedEntity.Id = payment.Id
	expectedResponse := expectedEntity.ToDto()

	assert.Nil(p.T(), err)
	assert.Equal(p.T(), expectedResponse, payment)
}

func (p *PaymentServiceTestSuite) Test_Add_WithUserNotExists() {
	p.userService.On("ExistsById", mocks.UserId).Return(false)

	request := mocks.GenerateCreatePaymentRequest()

	payment, err := p.TestO.Add(request)

	assert.Equal(p.T(), fmt.Errorf("user with id %s in not exists", request.UserId), err)
	assert.Equal(p.T(), model.PaymentDto{}, payment)
}

func (p *PaymentServiceTestSuite) Test_Add_WithHouseNotExists() {
	p.userService.On("ExistsById", mocks.UserId).Return(true)
	p.houseService.On("ExistsById", mocks.HouseId).Return(false)

	request := mocks.GenerateCreatePaymentRequest()

	payment, err := p.TestO.Add(request)

	assert.Equal(p.T(), fmt.Errorf("house with id %s in not exists", request.HouseId), err)
	assert.Equal(p.T(), model.PaymentDto{}, payment)
}

func (p *PaymentServiceTestSuite) Test_Add_WithProviderNotExists() {
	p.userService.On("ExistsById", mocks.UserId).Return(true)
	p.houseService.On("ExistsById", mocks.HouseId).Return(true)
	p.providerService.On("ExistsById", mocks.ProviderId).Return(false)

	request := mocks.GenerateCreatePaymentRequest()

	payment, err := p.TestO.Add(request)

	assert.Equal(p.T(), fmt.Errorf("provider with id %s in not exists", request.ProviderId), err)
	assert.Equal(p.T(), model.PaymentDto{}, payment)

	p.paymentRepository.AssertNotCalled(p.T(), "Create", mock.Anything)
}

func (p *PaymentServiceTestSuite) Test_FindById() {
	payment := mocks.GeneratePayment(mocks.HouseId, mocks.UserId, mocks.ProviderId)

	p.paymentRepository.On("FindById", payment.Id).Return(payment, nil)

	actual, err := p.TestO.FindById(payment.Id)

	assert.Nil(p.T(), err)
	assert.Equal(p.T(), payment.ToDto(), actual)
}

func (p *PaymentServiceTestSuite) Test_FindById_WithNotExistingId() {
	id := uuid.New()

	p.paymentRepository.On("FindById", id).Return(model.Payment{}, gorm.ErrRecordNotFound)

	actual, err := p.TestO.FindById(id)

	assert.Equal(p.T(), interrors.NewErrNotFound("payment with id %s not found", id), err)
	assert.Equal(p.T(), model.PaymentDto{}, actual)
}

func (p *PaymentServiceTestSuite) Test_FindById_WithError() {
	id := uuid.New()
	err2 := errors.New("error")

	p.paymentRepository.On("FindById", id).Return(model.Payment{}, err2)

	actual, err := p.TestO.FindById(id)

	assert.Equal(p.T(), err2, err)
	assert.Equal(p.T(), model.PaymentDto{}, actual)
}

func (p *PaymentServiceTestSuite) Test_FindByHouseId() {
	payment := mocks.GeneratePayment(mocks.HouseId, mocks.UserId, mocks.ProviderId)

	houseId := uuid.New()

	dto := payment.ToDto()
	p.paymentRepository.On("FindByHouseId", houseId).Return([]model.PaymentDto{dto})

	payments := p.TestO.FindByHouseId(houseId)

	assert.Equal(p.T(), []model.PaymentDto{dto}, payments)
}

func (p *PaymentServiceTestSuite) Test_FindByHouseId_WithNotExistingRecords() {
	houseId := uuid.New()

	p.paymentRepository.On("FindByHouseId", houseId).Return([]model.PaymentDto{})

	payments := p.TestO.FindByHouseId(houseId)

	assert.Equal(p.T(), []model.PaymentDto{}, payments)
}

func (p *PaymentServiceTestSuite) Test_FindByUserId() {
	payment := mocks.GeneratePayment(mocks.HouseId, mocks.UserId, mocks.ProviderId)

	userId := uuid.New()

	dto := payment.ToDto()
	p.paymentRepository.On("FindByUserId", userId).Return([]model.PaymentDto{dto})

	payments := p.TestO.FindByUserId(userId)

	assert.Equal(p.T(), []model.PaymentDto{dto}, payments)
}

func (p *PaymentServiceTestSuite) Test_FindByUserId_WithNotExistingRecords() {
	userId := uuid.New()

	p.paymentRepository.On("FindByUserId", userId).Return([]model.PaymentDto{})

	payments := p.TestO.FindByUserId(userId)

	assert.Equal(p.T(), []model.PaymentDto{}, payments)
}

func (p *PaymentServiceTestSuite) Test_FindByProviderId() {
	payment := mocks.GeneratePayment(mocks.HouseId, mocks.UserId, mocks.ProviderId)

	userId := uuid.New()

	dto := payment.ToDto()
	p.paymentRepository.On("FindByProviderId", userId).Return([]model.PaymentDto{dto})

	payments := p.TestO.FindByProviderId(userId)

	assert.Equal(p.T(), []model.PaymentDto{dto}, payments)
}

func (p *PaymentServiceTestSuite) Test_FindByProviderId_WithNotExistingRecords() {
	userId := uuid.New()

	p.paymentRepository.On("FindByProviderId", userId).Return([]model.PaymentDto{})

	payments := p.TestO.FindByProviderId(userId)

	assert.Equal(p.T(), []model.PaymentDto{}, payments)
}

func (p *PaymentServiceTestSuite) Test_ExistsById() {
	id := uuid.New()

	p.paymentRepository.On("ExistsById", id).Return(true)

	assert.True(p.T(), p.TestO.ExistsById(id))
}

func (p *PaymentServiceTestSuite) Test_ExistsById_WithNotExists() {
	id := uuid.New()

	p.paymentRepository.On("ExistsById", id).Return(false)

	assert.False(p.T(), p.TestO.ExistsById(id))
}

func (p *PaymentServiceTestSuite) Test_DeleteById() {
	id := uuid.New()

	p.paymentRepository.On("ExistsById", id).Return(true)
	p.paymentRepository.On("DeleteById", id).Return(nil)

	assert.Nil(p.T(), p.TestO.DeleteById(id))
}

func (p *PaymentServiceTestSuite) Test_DeleteById_WithNotExists() {
	id := uuid.New()

	p.paymentRepository.On("ExistsById", id).Return(false)

	assert.Equal(p.T(), fmt.Errorf("payment with id %s not found", id), p.TestO.DeleteById(id))

	p.paymentRepository.AssertNotCalled(p.T(), "DeleteById", id)
}

func (p *PaymentServiceTestSuite) Test_Update() {
	request := mocks.GenerateUpdatePaymentRequest()
	id := uuid.New()

	p.paymentRepository.On("ExistsById", id).Return(true)
	p.providerService.On("ExistsById", request.ProviderId).Return(true)
	p.paymentRepository.On("Update", mock.Anything).Return(nil)

	assert.Nil(p.T(), p.TestO.Update(id, request))

	p.paymentRepository.AssertCalled(p.T(), "Update", model.Payment{
		Id:          id,
		Name:        request.Name,
		Description: request.Description,
		ProviderId:  request.ProviderId,
		Date:        request.Date,
		Sum:         request.Sum,
	})
}

func (p *PaymentServiceTestSuite) Test_Update_WithErrorFromDatabase() {
	request := mocks.GenerateUpdatePaymentRequest()
	id := uuid.New()

	p.paymentRepository.On("ExistsById", id).Return(true)
	p.providerService.On("ExistsById", request.ProviderId).Return(true)
	p.paymentRepository.On("Update", mock.Anything).Return(errors.New("test"))

	err := p.TestO.Update(id, request)
	assert.Equal(p.T(), errors.New("test"), err)
}

func (p *PaymentServiceTestSuite) Test_Update_WithNotExists() {
	request := mocks.GenerateUpdatePaymentRequest()
	id := uuid.New()

	p.paymentRepository.On("ExistsById", id).Return(false)

	err := p.TestO.Update(id, request)
	assert.Equal(p.T(), fmt.Errorf("payment with id %s not found", id), err)

	p.paymentRepository.AssertNotCalled(p.T(), "Update", mock.Anything)
}

func (p *PaymentServiceTestSuite) Test_Update_WithDateAfterCurrentDate() {
	request := mocks.GenerateUpdatePaymentRequest()
	request.Date = time.Now().Add(time.Hour)
	id := uuid.New()

	p.paymentRepository.On("ExistsById", id).Return(true)
	p.providerService.On("ExistsById", request.ProviderId).Return(true)

	err := p.TestO.Update(id, request)
	assert.Equal(p.T(), errors.New("date should not be after current date"), err)

	p.paymentRepository.AssertNotCalled(p.T(), "Update", mock.Anything)
}

func (p *PaymentServiceTestSuite) Test_Update_WithProviderNotExists() {
	request := mocks.GenerateUpdatePaymentRequest()
	id := uuid.New()

	p.paymentRepository.On("ExistsById", id).Return(true)
	p.providerService.On("ExistsById", request.ProviderId).Return(false)

	err := p.TestO.Update(id, request)
	assert.Equal(p.T(), fmt.Errorf("provider with id %s not found", request.ProviderId), err)

	p.paymentRepository.AssertNotCalled(p.T(), "Update", mock.Anything)
}
