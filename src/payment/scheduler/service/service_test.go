package service

import (
	"errors"
	"fmt"
	houseMocks "github.com/VlasovArtem/hob/src/house/mocks"
	paymentMocks "github.com/VlasovArtem/hob/src/payment/mocks"
	paymentModel "github.com/VlasovArtem/hob/src/payment/model"
	"github.com/VlasovArtem/hob/src/payment/scheduler/mocks"
	paymentScheduler "github.com/VlasovArtem/hob/src/payment/scheduler/model"
	schedulerMocks "github.com/VlasovArtem/hob/src/scheduler/mocks"
	userMocks "github.com/VlasovArtem/hob/src/user/mocks"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"testing"
)

var (
	users                      *userMocks.UserService
	houses                     *houseMocks.HouseService
	payments                   *paymentMocks.PaymentService
	schedulers                 *schedulerMocks.ServiceScheduler
	paymentSchedulerRepository *mocks.PaymentSchedulerRepository
)

func serviceGenerator() PaymentSchedulerService {
	users = new(userMocks.UserService)
	houses = new(houseMocks.HouseService)
	payments = new(paymentMocks.PaymentService)
	schedulers = new(schedulerMocks.ServiceScheduler)
	paymentSchedulerRepository = new(mocks.PaymentSchedulerRepository)

	return NewPaymentSchedulerService(users, houses, payments, schedulers, paymentSchedulerRepository)
}

func Test_Add(t *testing.T) {
	service := serviceGenerator()

	request := mocks.GenerateCreatePaymentSchedulerRequest(mocks.HouseId, mocks.UserId)

	users.On("ExistsById", mocks.UserId).
		Return(true)
	houses.On("ExistsById", mocks.HouseId).
		Return(true)
	paymentSchedulerRepository.On("Create", mock.Anything).
		Return(
			func(model paymentScheduler.PaymentScheduler) paymentScheduler.PaymentScheduler {
				return model
			}, nil)
	schedulers.On("Add", mock.AnythingOfType("uuid.UUID"), "@daily", mock.Anything).
		Return(cron.EntryID(0), nil)

	payment, err := service.Add(request)

	expectedEntity := request.ToEntity()
	expectedEntity.Id = payment.Id
	expectedResponse := expectedEntity.ToDto()

	assert.Nil(t, err)
	assert.Equal(t, expectedResponse, payment)
	schedulers.AssertCalled(t, "Add", expectedEntity.Id, "@daily", mock.Anything)

	payments.On("Add", mock.Anything).Return(paymentModel.PaymentDto{}, nil)

	function := schedulers.Calls[0].Arguments.Get(2).(func())
	function()

	createPaymentRequest := payments.Calls[0].Arguments.Get(0).(paymentModel.CreatePaymentRequest)

	assert.Equal(t, paymentModel.CreatePaymentRequest{
		Name:        "Test Payment",
		Description: "Test Payment Description",
		HouseId:     mocks.HouseId,
		UserId:      mocks.UserId,
		Date:        createPaymentRequest.Date,
		Sum:         1000,
	}, createPaymentRequest)
}

func Test_Add_WithNegativeSum(t *testing.T) {
	service := serviceGenerator()

	request := mocks.GenerateCreatePaymentSchedulerRequest(mocks.HouseId, mocks.UserId)
	request.Sum = -1000

	payment, err := service.Add(request)

	assert.Equal(t, errors.New("sum should not be zero of negative"), err)
	assert.Equal(t, paymentScheduler.PaymentSchedulerDto{}, payment)
	schedulers.AssertNotCalled(t, "Create", mock.Anything, mock.Anything, mock.Anything)
}

func Test_Add_WithZeroSum(t *testing.T) {
	service := serviceGenerator()

	request := mocks.GenerateCreatePaymentSchedulerRequest(mocks.HouseId, mocks.UserId)
	request.Sum = 0

	payment, err := service.Add(request)

	assert.Equal(t, errors.New("sum should not be zero of negative"), err)
	assert.Equal(t, paymentScheduler.PaymentSchedulerDto{}, payment)
	schedulers.AssertNotCalled(t, "Create", mock.Anything, mock.Anything, mock.Anything)
}

func Test_Add_WithUserNotExists(t *testing.T) {
	service := serviceGenerator()

	users.On("ExistsById", mocks.UserId).Return(false)

	request := mocks.GenerateCreatePaymentSchedulerRequest(mocks.HouseId, mocks.UserId)

	payment, err := service.Add(request)

	assert.Equal(t, errors.New(fmt.Sprintf("user with id %s in not exists", request.UserId)), err)
	assert.Equal(t, paymentScheduler.PaymentSchedulerDto{}, payment)
}

func Test_Add_WithHouseNotExists(t *testing.T) {
	service := serviceGenerator()

	users.On("ExistsById", mocks.UserId).Return(true)
	houses.On("ExistsById", mocks.HouseId).Return(false)

	request := mocks.GenerateCreatePaymentSchedulerRequest(mocks.HouseId, mocks.UserId)

	payment, err := service.Add(request)

	assert.Equal(t, errors.New(fmt.Sprintf("house with id %s in not exists", request.HouseId)), err)
	assert.Equal(t, paymentScheduler.PaymentSchedulerDto{}, payment)
}

func Test_Add_WithInvalidSpec(t *testing.T) {
	service := serviceGenerator()

	users.On("ExistsById", mocks.UserId).
		Return(true)
	houses.On("ExistsById", mocks.HouseId).
		Return(true)
	schedulers.On("Create", mock.AnythingOfType("uuid.UUID"), "@daily", mock.Anything).
		Return(cron.EntryID(0), nil)

	request := mocks.GenerateCreatePaymentSchedulerRequest(mocks.HouseId, mocks.UserId)
	request.Spec = ""

	payment, err := service.Add(request)

	assert.Equal(t, errors.New("scheduler configuration not provided"), err)
	assert.Equal(t, paymentScheduler.PaymentSchedulerDto{}, payment)
}

func Test_Add_WithErrorDuringScheduling(t *testing.T) {
	service := serviceGenerator()

	users.On("ExistsById", mocks.UserId).
		Return(true)
	houses.On("ExistsById", mocks.HouseId).
		Return(true)
	paymentSchedulerRepository.On("Create", mock.Anything).
		Return(
			func(model paymentScheduler.PaymentScheduler) paymentScheduler.PaymentScheduler {
				return model
			}, nil)
	schedulers.On("Add", mock.AnythingOfType("uuid.UUID"), "@daily", mock.Anything).
		Return(cron.EntryID(0), errors.New("error"))
	paymentSchedulerRepository.On("DeleteById", mock.AnythingOfType("uuid.UUID")).Return()

	request := mocks.GenerateCreatePaymentSchedulerRequest(mocks.HouseId, mocks.UserId)

	payment, err := service.Add(request)

	paymentSchedulerRepository.AssertCalled(t, "DeleteById", mock.AnythingOfType("uuid.UUID"))

	assert.Equal(t, errors.New("error"), err)
	assert.Equal(t, paymentScheduler.PaymentSchedulerDto{}, payment)
}

func Test_Add_WithErrorDuringCreateScheduleEntity(t *testing.T) {
	service := serviceGenerator()

	users.On("ExistsById", mocks.UserId).
		Return(true)
	houses.On("ExistsById", mocks.HouseId).
		Return(true)
	paymentSchedulerRepository.On("Create", mock.Anything).
		Return(
			func(model paymentScheduler.PaymentScheduler) paymentScheduler.PaymentScheduler {
				return model
			}, errors.New("error"))

	request := mocks.GenerateCreatePaymentSchedulerRequest(mocks.HouseId, mocks.UserId)

	payment, err := service.Add(request)

	paymentSchedulerRepository.AssertNotCalled(t, "DeleteById", mock.AnythingOfType("uuid.UUID"))
	schedulers.AssertNotCalled(t, "Add", mock.AnythingOfType("uuid.UUID"), "@daily", mock.Anything)

	assert.Equal(t, errors.New("error"), err)
	assert.Equal(t, paymentScheduler.PaymentSchedulerDto{}, payment)
}

func Test_Remove(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()

	paymentSchedulerRepository.On("ExistsById", id).Return(true)
	paymentSchedulerRepository.On("DeleteById", id).Return()
	schedulers.On("Remove", id).Return(nil)

	err := service.Remove(id)

	paymentSchedulerRepository.AssertCalled(t, "DeleteById", id)

	assert.Nil(t, err)
}

func Test_Remove_WithErrorFromScheduler(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()

	paymentSchedulerRepository.On("ExistsById", id).Return(true)
	paymentSchedulerRepository.On("DeleteById", id).Return()

	schedulers.On("Remove", id).Return(errors.New("error"))

	err := service.Remove(id)

	paymentSchedulerRepository.AssertCalled(t, "DeleteById", id)

	assert.Nil(t, err)
}

func Test_Remove_WithMissingRecord(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()

	paymentSchedulerRepository.On("ExistsById", id).Return(false)

	err := service.Remove(id)

	paymentSchedulerRepository.AssertNotCalled(t, "DeleteById", id)

	assert.Equal(t, errors.New(fmt.Sprintf("payment scheduler with id %s not found", id)), err)
}

func Test_FindById(t *testing.T) {
	service := serviceGenerator()

	scheduler := mocks.GeneratePaymentScheduler(mocks.HouseId, mocks.UserId)

	paymentSchedulerRepository.On("FindById", scheduler.Id).Return(scheduler, nil)

	actual, err := service.FindById(scheduler.Id)

	assert.Nil(t, err)
	assert.Equal(t, scheduler.ToDto(), actual)
}

func Test_FindById_WithNotExistingId(t *testing.T) {
	paymentService := serviceGenerator()

	id := uuid.New()

	paymentSchedulerRepository.On("FindById", id).Return(paymentScheduler.PaymentScheduler{}, gorm.ErrRecordNotFound)

	actual, err := paymentService.FindById(id)

	assert.Equal(t, errors.New(fmt.Sprintf("payment scheduler with id %s not found", id)), err)
	assert.Equal(t, paymentScheduler.PaymentSchedulerDto{}, actual)
}

func Test_FindById_WithError(t *testing.T) {
	paymentService := serviceGenerator()

	id := uuid.New()
	expectedError := errors.New("error")

	paymentSchedulerRepository.On("FindById", id).Return(paymentScheduler.PaymentScheduler{}, expectedError)

	actual, err := paymentService.FindById(id)

	assert.Equal(t, expectedError, err)
	assert.Equal(t, paymentScheduler.PaymentSchedulerDto{}, actual)
}

func Test_FindByHouseId(t *testing.T) {
	service := serviceGenerator()

	scheduler := mocks.GeneratePaymentScheduler(mocks.HouseId, mocks.UserId)

	paymentSchedulerRepository.On("FindByHouseId", scheduler.HouseId).Return([]paymentScheduler.PaymentScheduler{scheduler})

	actual := service.FindByHouseId(scheduler.HouseId)

	assert.Equal(t, []paymentScheduler.PaymentSchedulerDto{scheduler.ToDto()}, actual)
}

func Test_FindByHouseId_WithNotExistingRecords(t *testing.T) {
	paymentService := serviceGenerator()

	id := uuid.New()

	paymentSchedulerRepository.On("FindByHouseId", id).Return([]paymentScheduler.PaymentScheduler{})

	actual := paymentService.FindByHouseId(id)

	assert.Equal(t, []paymentScheduler.PaymentSchedulerDto{}, actual)
}

func Test_FindByUserId(t *testing.T) {
	service := serviceGenerator()

	scheduler := mocks.GeneratePaymentScheduler(mocks.HouseId, mocks.UserId)

	paymentSchedulerRepository.On("FindByUserId", scheduler.UserId).Return([]paymentScheduler.PaymentScheduler{scheduler})

	actual := service.FindByUserId(scheduler.UserId)

	assert.Equal(t, []paymentScheduler.PaymentSchedulerDto{scheduler.ToDto()}, actual)
}

func Test_FindPaymentByUserId_WithNotExistingRecords(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()

	paymentSchedulerRepository.On("FindByUserId", id).Return([]paymentScheduler.PaymentScheduler{})

	actual := service.FindByUserId(id)

	assert.Equal(t, []paymentScheduler.PaymentSchedulerDto{}, actual)
}
