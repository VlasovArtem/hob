package service

import (
	"errors"
	"github.com/VlasovArtem/hob/src/common/int-errors"
	houseMocks "github.com/VlasovArtem/hob/src/house/mocks"
	paymentMocks "github.com/VlasovArtem/hob/src/payment/mocks"
	paymentModel "github.com/VlasovArtem/hob/src/payment/model"
	"github.com/VlasovArtem/hob/src/payment/scheduler/mocks"
	paymentScheduler "github.com/VlasovArtem/hob/src/payment/scheduler/model"
	providerMocks "github.com/VlasovArtem/hob/src/provider/mocks"
	schedulerMocks "github.com/VlasovArtem/hob/src/scheduler/mocks"
	"github.com/VlasovArtem/hob/src/test/testhelper"
	userMocks "github.com/VlasovArtem/hob/src/user/mocks"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type PaymentSchedulerServiceTestSuite struct {
	testhelper.MockTestSuite[PaymentSchedulerService]
	userService                *userMocks.UserService
	houseService               *houseMocks.HouseService
	paymentService             *paymentMocks.PaymentService
	serviceScheduler           *schedulerMocks.ServiceScheduler
	providerService            *providerMocks.ProviderService
	paymentSchedulerRepository *mocks.PaymentSchedulerRepository
}

func TestPaymentSchedulerServiceTestSuite(t *testing.T) {
	ts := &PaymentSchedulerServiceTestSuite{}
	ts.TestObjectGenerator = func() PaymentSchedulerService {
		ts.userService = new(userMocks.UserService)
		ts.houseService = new(houseMocks.HouseService)
		ts.paymentService = new(paymentMocks.PaymentService)
		ts.serviceScheduler = new(schedulerMocks.ServiceScheduler)
		ts.providerService = new(providerMocks.ProviderService)
		ts.paymentSchedulerRepository = new(mocks.PaymentSchedulerRepository)

		return NewPaymentSchedulerService(ts.userService, ts.houseService, ts.paymentService, ts.providerService, ts.serviceScheduler, ts.paymentSchedulerRepository)
	}

	suite.Run(t, ts)
}

func (p *PaymentSchedulerServiceTestSuite) Test_Add() {
	request := mocks.GenerateCreatePaymentSchedulerRequest()

	p.userService.On("ExistsById", mocks.UserId).
		Return(true)
	p.houseService.On("ExistsById", mocks.HouseId).
		Return(true)
	p.providerService.On("ExistsById", mocks.ProviderId).
		Return(true)
	p.paymentSchedulerRepository.On("Create", mock.Anything).
		Return(
			func(model paymentScheduler.PaymentScheduler) paymentScheduler.PaymentScheduler {
				return model
			}, nil)
	p.serviceScheduler.On("Add", mock.AnythingOfType("uuid.UUID"), "@daily", mock.Anything).
		Return(cron.EntryID(0), nil)

	payment, err := p.TestO.Add(request)

	expectedEntity := request.ToEntity()
	expectedEntity.Id = payment.Id
	expectedResponse := expectedEntity.ToDto()

	assert.Nil(p.T(), err)
	assert.Equal(p.T(), expectedResponse, payment)
	p.serviceScheduler.AssertCalled(p.T(), "Add", expectedEntity.Id, "@daily", mock.Anything)

	p.paymentService.On("Add", mock.Anything).Return(paymentModel.PaymentDto{}, nil)

	function := p.serviceScheduler.Calls[0].Arguments.Get(2).(func())
	function()

	createPaymentRequest := p.paymentService.Calls[0].Arguments.Get(0).(paymentModel.CreatePaymentRequest)

	assert.Equal(p.T(), paymentModel.CreatePaymentRequest{
		Name:        "Test Payment",
		Description: "Test Payment Description",
		HouseId:     mocks.HouseId,
		UserId:      mocks.UserId,
		ProviderId:  &mocks.ProviderId,
		Date:        createPaymentRequest.Date,
		Sum:         1000,
	}, createPaymentRequest)
}

func (p *PaymentSchedulerServiceTestSuite) Test_Add_WithNegativeSum() {
	request := mocks.GenerateCreatePaymentSchedulerRequest()
	request.Sum = -1000

	payment, err := p.TestO.Add(request)

	assert.Equal(p.T(), errors.New("sum should not be zero of negative"), err)
	assert.Equal(p.T(), paymentScheduler.PaymentSchedulerDto{}, payment)
	p.serviceScheduler.AssertNotCalled(p.T(), "Create", mock.Anything, mock.Anything, mock.Anything)
}

func (p *PaymentSchedulerServiceTestSuite) Test_Add_WithZeroSum() {
	request := mocks.GenerateCreatePaymentSchedulerRequest()
	request.Sum = 0

	payment, err := p.TestO.Add(request)

	assert.Equal(p.T(), errors.New("sum should not be zero of negative"), err)
	assert.Equal(p.T(), paymentScheduler.PaymentSchedulerDto{}, payment)
	p.serviceScheduler.AssertNotCalled(p.T(), "Create", mock.Anything, mock.Anything, mock.Anything)
}

func (p *PaymentSchedulerServiceTestSuite) Test_Add_WithUserNotExists() {
	p.userService.On("ExistsById", mocks.UserId).Return(false)

	request := mocks.GenerateCreatePaymentSchedulerRequest()

	payment, err := p.TestO.Add(request)

	assert.Equal(p.T(), int_errors.NewErrNotFound("user with id %s in not exists", request.UserId), err)
	assert.Equal(p.T(), paymentScheduler.PaymentSchedulerDto{}, payment)
}

func (p *PaymentSchedulerServiceTestSuite) Test_Add_WithHouseNotExists() {
	p.userService.On("ExistsById", mocks.UserId).Return(true)
	p.houseService.On("ExistsById", mocks.HouseId).Return(false)

	request := mocks.GenerateCreatePaymentSchedulerRequest()

	payment, err := p.TestO.Add(request)

	assert.Equal(p.T(), int_errors.NewErrNotFound("house with id %s in not exists", request.HouseId), err)
	assert.Equal(p.T(), paymentScheduler.PaymentSchedulerDto{}, payment)
}

func (p *PaymentSchedulerServiceTestSuite) Test_Add_WithProviderNotExists() {
	p.userService.On("ExistsById", mocks.UserId).Return(true)
	p.houseService.On("ExistsById", mocks.HouseId).Return(true)
	p.providerService.On("ExistsById", mocks.ProviderId).Return(false)

	request := mocks.GenerateCreatePaymentSchedulerRequest()

	payment, err := p.TestO.Add(request)

	assert.Equal(p.T(), int_errors.NewErrNotFound("provider with id %s in not exists", request.ProviderId), err)
	assert.Equal(p.T(), paymentScheduler.PaymentSchedulerDto{}, payment)
}

func (p *PaymentSchedulerServiceTestSuite) Test_Add_WithInvalidSpec() {
	p.userService.On("ExistsById", mocks.UserId).
		Return(true)
	p.houseService.On("ExistsById", mocks.HouseId).
		Return(true)
	p.providerService.On("ExistsById", mocks.ProviderId).Return(true)
	p.serviceScheduler.On("Create", mock.AnythingOfType("uuid.UUID"), "@daily", mock.Anything).
		Return(cron.EntryID(0), nil)

	request := mocks.GenerateCreatePaymentSchedulerRequest()
	request.Spec = ""

	payment, err := p.TestO.Add(request)

	assert.Equal(p.T(), errors.New("scheduler configuration not provided"), err)
	assert.Equal(p.T(), paymentScheduler.PaymentSchedulerDto{}, payment)
}

func (p *PaymentSchedulerServiceTestSuite) Test_Add_WithErrorDuringScheduling() {
	p.userService.On("ExistsById", mocks.UserId).
		Return(true)
	p.houseService.On("ExistsById", mocks.HouseId).
		Return(true)
	p.providerService.On("ExistsById", mocks.ProviderId).Return(true)
	p.paymentSchedulerRepository.On("Create", mock.Anything).
		Return(
			func(model paymentScheduler.PaymentScheduler) paymentScheduler.PaymentScheduler {
				return model
			}, nil)
	p.serviceScheduler.On("Add", mock.AnythingOfType("uuid.UUID"), "@daily", mock.Anything).
		Return(cron.EntryID(0), errors.New("error"))
	p.paymentSchedulerRepository.On("DeleteById", mock.AnythingOfType("uuid.UUID")).Return()

	request := mocks.GenerateCreatePaymentSchedulerRequest()

	payment, err := p.TestO.Add(request)

	p.paymentSchedulerRepository.AssertCalled(p.T(), "DeleteById", mock.AnythingOfType("uuid.UUID"))

	assert.Equal(p.T(), errors.New("error"), err)
	assert.Equal(p.T(), paymentScheduler.PaymentSchedulerDto{}, payment)
}

func (p *PaymentSchedulerServiceTestSuite) Test_Add_WithErrorDuringCreateScheduleEntity() {
	p.userService.On("ExistsById", mocks.UserId).
		Return(true)
	p.houseService.On("ExistsById", mocks.HouseId).
		Return(true)
	p.providerService.On("ExistsById", mocks.ProviderId).Return(true)
	p.paymentSchedulerRepository.On("Create", mock.Anything).
		Return(
			func(model paymentScheduler.PaymentScheduler) paymentScheduler.PaymentScheduler {
				return model
			}, errors.New("error"))

	request := mocks.GenerateCreatePaymentSchedulerRequest()

	payment, err := p.TestO.Add(request)

	p.paymentSchedulerRepository.AssertNotCalled(p.T(), "DeleteById", mock.AnythingOfType("uuid.UUID"))
	p.serviceScheduler.AssertNotCalled(p.T(), "Add", mock.AnythingOfType("uuid.UUID"), "@daily", mock.Anything)

	assert.Equal(p.T(), errors.New("error"), err)
	assert.Equal(p.T(), paymentScheduler.PaymentSchedulerDto{}, payment)
}

func (p *PaymentSchedulerServiceTestSuite) Test_Remove() {
	id := uuid.New()

	p.paymentSchedulerRepository.On("ExistsById", id).Return(true)
	p.paymentSchedulerRepository.On("DeleteById", id).Return()
	p.serviceScheduler.On("Remove", id).Return(nil)

	err := p.TestO.Remove(id)

	p.paymentSchedulerRepository.AssertCalled(p.T(), "DeleteById", id)

	assert.Nil(p.T(), err)
}

func (p *PaymentSchedulerServiceTestSuite) Test_Remove_WithErrorFromScheduler() {
	id := uuid.New()

	p.paymentSchedulerRepository.On("ExistsById", id).Return(true)
	p.paymentSchedulerRepository.On("DeleteById", id).Return()

	p.serviceScheduler.On("Remove", id).Return(errors.New("error"))

	err := p.TestO.Remove(id)

	p.paymentSchedulerRepository.AssertNotCalled(p.T(), "DeleteById", id)

	assert.Nil(p.T(), err)
}

func (p *PaymentSchedulerServiceTestSuite) Test_Remove_WithMissingRecord() {
	id := uuid.New()

	p.paymentSchedulerRepository.On("ExistsById", id).Return(false)

	err := p.TestO.Remove(id)

	p.paymentSchedulerRepository.AssertNotCalled(p.T(), "DeleteById", id)

	assert.Equal(p.T(), int_errors.NewErrNotFound("payment scheduler with id %s not found", id), err)
}

func (p *PaymentSchedulerServiceTestSuite) Test_FindById() {
	scheduler := mocks.GeneratePaymentScheduler(mocks.HouseId, mocks.UserId, mocks.ProviderId)

	p.paymentSchedulerRepository.On("FindById", scheduler.Id).Return(scheduler, nil)

	actual, err := p.TestO.FindById(scheduler.Id)

	assert.Nil(p.T(), err)
	assert.Equal(p.T(), scheduler.ToDto(), actual)
}

func (p *PaymentSchedulerServiceTestSuite) Test_FindById_WithNotExistingId() {
	id := uuid.New()

	p.paymentSchedulerRepository.On("FindById", id).Return(paymentScheduler.PaymentScheduler{}, gorm.ErrRecordNotFound)

	actual, err := p.TestO.FindById(id)

	assert.Equal(p.T(), int_errors.NewErrNotFound("payment scheduler with id %s not found", id), err)
	assert.Equal(p.T(), paymentScheduler.PaymentSchedulerDto{}, actual)
}

func (p *PaymentSchedulerServiceTestSuite) Test_FindById_WithError() {
	id := uuid.New()
	expectedError := errors.New("error")

	p.paymentSchedulerRepository.On("FindById", id).Return(paymentScheduler.PaymentScheduler{}, expectedError)

	actual, err := p.TestO.FindById(id)

	assert.Equal(p.T(), expectedError, err)
	assert.Equal(p.T(), paymentScheduler.PaymentSchedulerDto{}, actual)
}

func (p *PaymentSchedulerServiceTestSuite) Test_FindByHouseId() {
	scheduler := mocks.GeneratePaymentScheduler(mocks.HouseId, mocks.UserId, mocks.ProviderId)
	dto := scheduler.ToDto()

	p.paymentSchedulerRepository.On("FindByHouseId", scheduler.HouseId).Return([]paymentScheduler.PaymentSchedulerDto{dto})

	actual := p.TestO.FindByHouseId(scheduler.HouseId)

	assert.Equal(p.T(), []paymentScheduler.PaymentSchedulerDto{dto}, actual)
}

func (p *PaymentSchedulerServiceTestSuite) Test_FindByHouseId_WithNotExistingRecords() {
	id := uuid.New()

	p.paymentSchedulerRepository.On("FindByHouseId", id).Return([]paymentScheduler.PaymentSchedulerDto{})

	actual := p.TestO.FindByHouseId(id)

	assert.Equal(p.T(), []paymentScheduler.PaymentSchedulerDto{}, actual)
}

func (p *PaymentSchedulerServiceTestSuite) Test_FindByUserId() {
	scheduler := mocks.GeneratePaymentScheduler(mocks.HouseId, mocks.UserId, mocks.ProviderId)
	dto := scheduler.ToDto()

	p.paymentSchedulerRepository.On("FindByUserId", scheduler.UserId).Return([]paymentScheduler.PaymentSchedulerDto{dto})

	actual := p.TestO.FindByUserId(scheduler.UserId)

	assert.Equal(p.T(), []paymentScheduler.PaymentSchedulerDto{dto}, actual)
}

func (p *PaymentSchedulerServiceTestSuite) Test_FindByUserId_WithNotExistingRecords() {
	id := uuid.New()

	p.paymentSchedulerRepository.On("FindByUserId", id).Return([]paymentScheduler.PaymentSchedulerDto{})

	actual := p.TestO.FindByUserId(id)

	assert.Equal(p.T(), []paymentScheduler.PaymentSchedulerDto{}, actual)
}

func (p *PaymentSchedulerServiceTestSuite) Test_FindByProviderId() {
	scheduler := mocks.GeneratePaymentScheduler(mocks.HouseId, mocks.UserId, mocks.ProviderId)
	dto := scheduler.ToDto()

	p.paymentSchedulerRepository.On("FindByProviderId", scheduler.ProviderId).Return([]paymentScheduler.PaymentSchedulerDto{dto})

	actual := p.TestO.FindByProviderId(scheduler.ProviderId)

	assert.Equal(p.T(), []paymentScheduler.PaymentSchedulerDto{dto}, actual)
}

func (p *PaymentSchedulerServiceTestSuite) Test_FindByProviderId_WithNotExistingRecords() {
	id := uuid.New()

	p.paymentSchedulerRepository.On("FindByProviderId", id).Return([]paymentScheduler.PaymentSchedulerDto{})

	actual := p.TestO.FindByProviderId(id)

	assert.Equal(p.T(), []paymentScheduler.PaymentSchedulerDto{}, actual)
}

func (p *PaymentSchedulerServiceTestSuite) Test_Update() {
	id, request := mocks.GenerateUpdatePaymentSchedulerRequest()
	scheduler := mocks.GeneratePaymentScheduler(uuid.New(), uuid.New(), uuid.New())
	scheduler.Id = id

	p.paymentSchedulerRepository.On("ExistsById", id).Return(true)
	p.providerService.On("ExistsById", request.ProviderId).Return(true)
	p.paymentSchedulerRepository.On("Update", id, request).Return(scheduler, nil)
	p.serviceScheduler.On("Update", mock.AnythingOfType("uuid.UUID"), string(request.Spec), mock.Anything).
		Return(cron.EntryID(0), nil)

	entity := request.ToEntity(id)
	entity.HouseId = mocks.HouseId
	entity.UserId = mocks.UserId

	p.paymentSchedulerRepository.On("FindById", id).Return(entity, nil)

	err := p.TestO.Update(id, request)

	assert.Nil(p.T(), err)
	p.serviceScheduler.AssertCalled(p.T(), "Update", id, string(request.Spec), mock.Anything)

	p.paymentSchedulerRepository.AssertNotCalled(p.T(), "DeleteById", mock.Anything)
}

func (p *PaymentSchedulerServiceTestSuite) Test_Update_WithMissingScheduler() {
	id, request := mocks.GenerateUpdatePaymentSchedulerRequest()

	p.paymentSchedulerRepository.On("ExistsById", id).Return(false)

	err := p.TestO.Update(id, request)

	assert.Equal(p.T(), int_errors.NewErrNotFound("payment schedule with id %s not found", id), err)

	p.paymentSchedulerRepository.AssertNotCalled(p.T(), "Update", mock.Anything)
	p.paymentSchedulerRepository.AssertNotCalled(p.T(), "DeleteById", mock.Anything)
}

func (p *PaymentSchedulerServiceTestSuite) Test_Update_WithInvalidSum() {
	id, request := mocks.GenerateUpdatePaymentSchedulerRequest()
	request.Sum = 0

	p.paymentSchedulerRepository.On("ExistsById", mock.Anything).Return(true)

	err := p.TestO.Update(id, request)

	assert.Equal(p.T(), errors.New("sum should not be zero of negative"), err)

	p.paymentSchedulerRepository.AssertNotCalled(p.T(), "Update", mock.Anything)
	p.paymentSchedulerRepository.AssertNotCalled(p.T(), "DeleteById", mock.Anything)
}

func (p *PaymentSchedulerServiceTestSuite) Test_Update_WithNotExistsProvider() {
	id, request := mocks.GenerateUpdatePaymentSchedulerRequest()

	p.paymentSchedulerRepository.On("ExistsById", mock.Anything).Return(true)
	p.providerService.On("ExistsById", request.ProviderId).Return(false)

	err := p.TestO.Update(id, request)

	assert.Equal(p.T(), int_errors.NewErrNotFound("provider with id %s not found", request.ProviderId), err)

	p.paymentSchedulerRepository.AssertNotCalled(p.T(), "Update", mock.Anything)
	p.paymentSchedulerRepository.AssertNotCalled(p.T(), "DeleteById", mock.Anything)
}

func (p *PaymentSchedulerServiceTestSuite) Test_Update_WithNotSchedulerSpec() {
	id, request := mocks.GenerateUpdatePaymentSchedulerRequest()
	request.Spec = ""

	p.paymentSchedulerRepository.On("ExistsById", mock.Anything).Return(true)
	p.providerService.On("ExistsById", request.ProviderId).Return(true)

	err := p.TestO.Update(id, request)

	assert.Equal(p.T(), errors.New("scheduler configuration not provided"), err)

	p.paymentSchedulerRepository.AssertNotCalled(p.T(), "Update", mock.Anything)
	p.paymentSchedulerRepository.AssertNotCalled(p.T(), "DeleteById", mock.Anything)
}

func (p *PaymentSchedulerServiceTestSuite) Test_Update_WithErrorFromUpdate() {
	id, request := mocks.GenerateUpdatePaymentSchedulerRequest()

	p.paymentSchedulerRepository.On("ExistsById", mock.Anything).Return(true)
	p.providerService.On("ExistsById", request.ProviderId).Return(true)
	p.paymentSchedulerRepository.On("Update", id, request).Return(paymentScheduler.PaymentScheduler{}, errors.New("error"))

	err := p.TestO.Update(id, request)

	assert.Equal(p.T(), errors.New("error"), err)

	p.paymentSchedulerRepository.AssertNotCalled(p.T(), "DeleteById", mock.Anything)
}

func (p *PaymentSchedulerServiceTestSuite) Test_Update_WithErrorFromScheduler() {
	id, request := mocks.GenerateUpdatePaymentSchedulerRequest()
	scheduler := mocks.GeneratePaymentScheduler(uuid.New(), uuid.New(), uuid.New())

	p.paymentSchedulerRepository.On("ExistsById", mock.Anything).Return(true)
	p.providerService.On("ExistsById", request.ProviderId).Return(true)
	p.paymentSchedulerRepository.On("Update", id, request).Return(scheduler, nil)
	p.serviceScheduler.On("Update", mock.AnythingOfType("uuid.UUID"), string(request.Spec), mock.Anything).
		Return(cron.EntryID(0), errors.New("error"))
	p.paymentSchedulerRepository.On("DeleteById", id).Return()

	entity := request.ToEntity(id)
	entity.HouseId = mocks.HouseId
	entity.UserId = mocks.UserId

	p.paymentSchedulerRepository.On("FindById", id).Return(entity, nil)

	err := p.TestO.Update(id, request)

	assert.Equal(p.T(), errors.New("error"), err)
}
