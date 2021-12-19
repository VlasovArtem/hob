package service

import (
	"errors"
	"fmt"
	houseMocks "github.com/VlasovArtem/hob/src/house/mocks"
	paymentMocks "github.com/VlasovArtem/hob/src/payment/mocks"
	paymentModel "github.com/VlasovArtem/hob/src/payment/model"
	"github.com/VlasovArtem/hob/src/payment/scheduler/mocks"
	paymentScheduler "github.com/VlasovArtem/hob/src/payment/scheduler/model"
	providerMocks "github.com/VlasovArtem/hob/src/provider/mocks"
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
	userService                *userMocks.UserService
	houseService               *houseMocks.HouseService
	paymentService             *paymentMocks.PaymentService
	serviceScheduler           *schedulerMocks.ServiceScheduler
	providerService            *providerMocks.ProviderService
	paymentSchedulerRepository *mocks.PaymentSchedulerRepository
)

func serviceGenerator() PaymentSchedulerService {
	userService = new(userMocks.UserService)
	houseService = new(houseMocks.HouseService)
	paymentService = new(paymentMocks.PaymentService)
	serviceScheduler = new(schedulerMocks.ServiceScheduler)
	providerService = new(providerMocks.ProviderService)
	paymentSchedulerRepository = new(mocks.PaymentSchedulerRepository)

	return NewPaymentSchedulerService(userService, houseService, paymentService, providerService, serviceScheduler, paymentSchedulerRepository)
}

func Test_Add(t *testing.T) {
	service := serviceGenerator()

	request := mocks.GenerateCreatePaymentSchedulerRequest()

	userService.On("ExistsById", mocks.UserId).
		Return(true)
	houseService.On("ExistsById", mocks.HouseId).
		Return(true)
	providerService.On("ExistsById", mocks.ProviderId).
		Return(true)
	paymentSchedulerRepository.On("Create", mock.Anything).
		Return(
			func(model paymentScheduler.PaymentScheduler) paymentScheduler.PaymentScheduler {
				return model
			}, nil)
	serviceScheduler.On("Add", mock.AnythingOfType("uuid.UUID"), "@daily", mock.Anything).
		Return(cron.EntryID(0), nil)

	payment, err := service.Add(request)

	expectedEntity := request.ToEntity()
	expectedEntity.Id = payment.Id
	expectedResponse := expectedEntity.ToDto()

	assert.Nil(t, err)
	assert.Equal(t, expectedResponse, payment)
	serviceScheduler.AssertCalled(t, "Add", expectedEntity.Id, "@daily", mock.Anything)

	paymentService.On("Add", mock.Anything).Return(paymentModel.PaymentDto{}, nil)

	function := serviceScheduler.Calls[0].Arguments.Get(2).(func())
	function()

	createPaymentRequest := paymentService.Calls[0].Arguments.Get(0).(paymentModel.CreatePaymentRequest)

	assert.Equal(t, paymentModel.CreatePaymentRequest{
		Name:        "Test Payment",
		Description: "Test Payment Description",
		HouseId:     mocks.HouseId,
		UserId:      mocks.UserId,
		ProviderId:  mocks.ProviderId,
		Date:        createPaymentRequest.Date,
		Sum:         1000,
	}, createPaymentRequest)
}

func Test_Add_WithNegativeSum(t *testing.T) {
	service := serviceGenerator()

	request := mocks.GenerateCreatePaymentSchedulerRequest()
	request.Sum = -1000

	payment, err := service.Add(request)

	assert.Equal(t, errors.New("sum should not be zero of negative"), err)
	assert.Equal(t, paymentScheduler.PaymentSchedulerDto{}, payment)
	serviceScheduler.AssertNotCalled(t, "Create", mock.Anything, mock.Anything, mock.Anything)
}

func Test_Add_WithZeroSum(t *testing.T) {
	service := serviceGenerator()

	request := mocks.GenerateCreatePaymentSchedulerRequest()
	request.Sum = 0

	payment, err := service.Add(request)

	assert.Equal(t, errors.New("sum should not be zero of negative"), err)
	assert.Equal(t, paymentScheduler.PaymentSchedulerDto{}, payment)
	serviceScheduler.AssertNotCalled(t, "Create", mock.Anything, mock.Anything, mock.Anything)
}

func Test_Add_WithUserNotExists(t *testing.T) {
	service := serviceGenerator()

	userService.On("ExistsById", mocks.UserId).Return(false)

	request := mocks.GenerateCreatePaymentSchedulerRequest()

	payment, err := service.Add(request)

	assert.Equal(t, errors.New(fmt.Sprintf("user with id %s in not exists", request.UserId)), err)
	assert.Equal(t, paymentScheduler.PaymentSchedulerDto{}, payment)
}

func Test_Add_WithHouseNotExists(t *testing.T) {
	service := serviceGenerator()

	userService.On("ExistsById", mocks.UserId).Return(true)
	houseService.On("ExistsById", mocks.HouseId).Return(false)

	request := mocks.GenerateCreatePaymentSchedulerRequest()

	payment, err := service.Add(request)

	assert.Equal(t, errors.New(fmt.Sprintf("house with id %s in not exists", request.HouseId)), err)
	assert.Equal(t, paymentScheduler.PaymentSchedulerDto{}, payment)
}

func Test_Add_WithProviderNotExists(t *testing.T) {
	service := serviceGenerator()

	userService.On("ExistsById", mocks.UserId).Return(true)
	houseService.On("ExistsById", mocks.HouseId).Return(true)
	providerService.On("ExistsById", mocks.ProviderId).Return(false)

	request := mocks.GenerateCreatePaymentSchedulerRequest()

	payment, err := service.Add(request)

	assert.Equal(t, errors.New(fmt.Sprintf("provider with id %s in not exists", request.ProviderId)), err)
	assert.Equal(t, paymentScheduler.PaymentSchedulerDto{}, payment)
}

func Test_Add_WithInvalidSpec(t *testing.T) {
	service := serviceGenerator()

	userService.On("ExistsById", mocks.UserId).
		Return(true)
	houseService.On("ExistsById", mocks.HouseId).
		Return(true)
	providerService.On("ExistsById", mocks.ProviderId).Return(true)
	serviceScheduler.On("Create", mock.AnythingOfType("uuid.UUID"), "@daily", mock.Anything).
		Return(cron.EntryID(0), nil)

	request := mocks.GenerateCreatePaymentSchedulerRequest()
	request.Spec = ""

	payment, err := service.Add(request)

	assert.Equal(t, errors.New("scheduler configuration not provided"), err)
	assert.Equal(t, paymentScheduler.PaymentSchedulerDto{}, payment)
}

func Test_Add_WithErrorDuringScheduling(t *testing.T) {
	service := serviceGenerator()

	userService.On("ExistsById", mocks.UserId).
		Return(true)
	houseService.On("ExistsById", mocks.HouseId).
		Return(true)
	providerService.On("ExistsById", mocks.ProviderId).Return(true)
	paymentSchedulerRepository.On("Create", mock.Anything).
		Return(
			func(model paymentScheduler.PaymentScheduler) paymentScheduler.PaymentScheduler {
				return model
			}, nil)
	serviceScheduler.On("Add", mock.AnythingOfType("uuid.UUID"), "@daily", mock.Anything).
		Return(cron.EntryID(0), errors.New("error"))
	paymentSchedulerRepository.On("DeleteById", mock.AnythingOfType("uuid.UUID")).Return()

	request := mocks.GenerateCreatePaymentSchedulerRequest()

	payment, err := service.Add(request)

	paymentSchedulerRepository.AssertCalled(t, "DeleteById", mock.AnythingOfType("uuid.UUID"))

	assert.Equal(t, errors.New("error"), err)
	assert.Equal(t, paymentScheduler.PaymentSchedulerDto{}, payment)
}

func Test_Add_WithErrorDuringCreateScheduleEntity(t *testing.T) {
	service := serviceGenerator()

	userService.On("ExistsById", mocks.UserId).
		Return(true)
	houseService.On("ExistsById", mocks.HouseId).
		Return(true)
	providerService.On("ExistsById", mocks.ProviderId).Return(true)
	paymentSchedulerRepository.On("Create", mock.Anything).
		Return(
			func(model paymentScheduler.PaymentScheduler) paymentScheduler.PaymentScheduler {
				return model
			}, errors.New("error"))

	request := mocks.GenerateCreatePaymentSchedulerRequest()

	payment, err := service.Add(request)

	paymentSchedulerRepository.AssertNotCalled(t, "DeleteById", mock.AnythingOfType("uuid.UUID"))
	serviceScheduler.AssertNotCalled(t, "Add", mock.AnythingOfType("uuid.UUID"), "@daily", mock.Anything)

	assert.Equal(t, errors.New("error"), err)
	assert.Equal(t, paymentScheduler.PaymentSchedulerDto{}, payment)
}

func Test_Remove(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()

	paymentSchedulerRepository.On("ExistsById", id).Return(true)
	paymentSchedulerRepository.On("DeleteById", id).Return()
	serviceScheduler.On("Remove", id).Return(nil)

	err := service.Remove(id)

	paymentSchedulerRepository.AssertCalled(t, "DeleteById", id)

	assert.Nil(t, err)
}

func Test_Remove_WithErrorFromScheduler(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()

	paymentSchedulerRepository.On("ExistsById", id).Return(true)
	paymentSchedulerRepository.On("DeleteById", id).Return()

	serviceScheduler.On("Remove", id).Return(errors.New("error"))

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

	scheduler := mocks.GeneratePaymentScheduler(mocks.HouseId, mocks.UserId, mocks.ProviderId)

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

	scheduler := mocks.GeneratePaymentScheduler(mocks.HouseId, mocks.UserId, mocks.ProviderId)
	dto := scheduler.ToDto()

	paymentSchedulerRepository.On("FindByHouseId", scheduler.HouseId).Return([]paymentScheduler.PaymentSchedulerDto{dto})

	actual := service.FindByHouseId(scheduler.HouseId)

	assert.Equal(t, []paymentScheduler.PaymentSchedulerDto{dto}, actual)
}

func Test_FindByHouseId_WithNotExistingRecords(t *testing.T) {
	paymentService := serviceGenerator()

	id := uuid.New()

	paymentSchedulerRepository.On("FindByHouseId", id).Return([]paymentScheduler.PaymentSchedulerDto{})

	actual := paymentService.FindByHouseId(id)

	assert.Equal(t, []paymentScheduler.PaymentSchedulerDto{}, actual)
}

func Test_FindByUserId(t *testing.T) {
	service := serviceGenerator()

	scheduler := mocks.GeneratePaymentScheduler(mocks.HouseId, mocks.UserId, mocks.ProviderId)
	dto := scheduler.ToDto()

	paymentSchedulerRepository.On("FindByUserId", scheduler.UserId).Return([]paymentScheduler.PaymentSchedulerDto{dto})

	actual := service.FindByUserId(scheduler.UserId)

	assert.Equal(t, []paymentScheduler.PaymentSchedulerDto{dto}, actual)
}

func Test_FindByUserId_WithNotExistingRecords(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()

	paymentSchedulerRepository.On("FindByUserId", id).Return([]paymentScheduler.PaymentSchedulerDto{})

	actual := service.FindByUserId(id)

	assert.Equal(t, []paymentScheduler.PaymentSchedulerDto{}, actual)
}

func Test_FindByProviderId(t *testing.T) {
	service := serviceGenerator()

	scheduler := mocks.GeneratePaymentScheduler(mocks.HouseId, mocks.UserId, mocks.ProviderId)
	dto := scheduler.ToDto()

	paymentSchedulerRepository.On("FindByProviderId", scheduler.ProviderId).Return([]paymentScheduler.PaymentSchedulerDto{dto})

	actual := service.FindByProviderId(scheduler.ProviderId)

	assert.Equal(t, []paymentScheduler.PaymentSchedulerDto{dto}, actual)
}

func Test_FindByProviderId_WithNotExistingRecords(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()

	paymentSchedulerRepository.On("FindByProviderId", id).Return([]paymentScheduler.PaymentSchedulerDto{})

	actual := service.FindByProviderId(id)

	assert.Equal(t, []paymentScheduler.PaymentSchedulerDto{}, actual)
}

func Test_Update(t *testing.T) {
	service := serviceGenerator()

	request := mocks.GenerateUpdatePaymentSchedulerRequest()
	schedulerId := uuid.New()

	paymentSchedulerRepository.On("ExistsById", schedulerId).Return(true)
	providerService.On("ExistsById", request.ProviderId).Return(true)
	paymentSchedulerRepository.On("Update", mock.Anything).Return(nil)
	serviceScheduler.On("Update", mock.AnythingOfType("uuid.UUID"), string(request.Spec), mock.Anything).
		Return(cron.EntryID(0), nil)

	entity := request.ToEntity(schedulerId)
	entity.HouseId = mocks.HouseId
	entity.UserId = mocks.UserId

	paymentSchedulerRepository.On("FindById", schedulerId).Return(entity, nil)

	err := service.Update(schedulerId, request)

	assert.Nil(t, err)
	paymentSchedulerRepository.AssertCalled(t, "Update", request.ToEntity(schedulerId))
	serviceScheduler.AssertCalled(t, "Update", schedulerId, string(request.Spec), mock.Anything)

	paymentService.On("Add", mock.Anything).Return(paymentModel.PaymentDto{}, nil)

	function := serviceScheduler.Calls[0].Arguments.Get(2).(func())
	function()

	createPaymentRequest := paymentService.Calls[0].Arguments.Get(0).(paymentModel.CreatePaymentRequest)

	assert.Equal(t, paymentModel.CreatePaymentRequest{
		Name:        "Test Payment Updated",
		Description: "Test Payment Description Updated",
		HouseId:     mocks.HouseId,
		UserId:      mocks.UserId,
		ProviderId:  request.ProviderId,
		Date:        createPaymentRequest.Date,
		Sum:         1000,
	}, createPaymentRequest)

	paymentSchedulerRepository.AssertNotCalled(t, "DeleteById", mock.Anything)
}

func Test_Update_WithMissingScheduler(t *testing.T) {
	service := serviceGenerator()

	request := mocks.GenerateUpdatePaymentSchedulerRequest()
	request.Sum = 0
	schedulerId := uuid.New()

	paymentSchedulerRepository.On("ExistsById", schedulerId).Return(false)

	err := service.Update(schedulerId, request)

	assert.Equal(t, errors.New(fmt.Sprintf("payment schedule with id %s not found", schedulerId)), err)

	paymentSchedulerRepository.AssertNotCalled(t, "Update", mock.Anything)
	paymentSchedulerRepository.AssertNotCalled(t, "DeleteById", mock.Anything)
}

func Test_Update_WithInvalidSum(t *testing.T) {
	service := serviceGenerator()

	request := mocks.GenerateUpdatePaymentSchedulerRequest()
	request.Sum = 0

	paymentSchedulerRepository.On("ExistsById", mock.Anything).Return(true)

	err := service.Update(uuid.New(), request)

	assert.Equal(t, errors.New("sum should not be zero of negative"), err)

	paymentSchedulerRepository.AssertNotCalled(t, "Update", mock.Anything)
	paymentSchedulerRepository.AssertNotCalled(t, "DeleteById", mock.Anything)
}

func Test_Update_WithNotExistsProvider(t *testing.T) {
	service := serviceGenerator()

	request := mocks.GenerateUpdatePaymentSchedulerRequest()

	paymentSchedulerRepository.On("ExistsById", mock.Anything).Return(true)
	providerService.On("ExistsById", request.ProviderId).Return(false)

	err := service.Update(uuid.New(), request)

	assert.Equal(t, errors.New(fmt.Sprintf("provider with id %s not found", request.ProviderId)), err)

	paymentSchedulerRepository.AssertNotCalled(t, "Update", mock.Anything)
	paymentSchedulerRepository.AssertNotCalled(t, "DeleteById", mock.Anything)
}

func Test_Update_WithNotSchedulerSpec(t *testing.T) {
	service := serviceGenerator()

	request := mocks.GenerateUpdatePaymentSchedulerRequest()
	request.Spec = ""

	paymentSchedulerRepository.On("ExistsById", mock.Anything).Return(true)
	providerService.On("ExistsById", request.ProviderId).Return(true)

	err := service.Update(uuid.New(), request)

	assert.Equal(t, errors.New("scheduler configuration not provided"), err)

	paymentSchedulerRepository.AssertNotCalled(t, "Update", mock.Anything)
	paymentSchedulerRepository.AssertNotCalled(t, "DeleteById", mock.Anything)
}

func Test_Update_WithErrorFromUpdate(t *testing.T) {
	service := serviceGenerator()

	request := mocks.GenerateUpdatePaymentSchedulerRequest()

	paymentSchedulerRepository.On("ExistsById", mock.Anything).Return(true)
	providerService.On("ExistsById", request.ProviderId).Return(true)
	paymentSchedulerRepository.On("Update", mock.Anything).Return(errors.New("error"))

	err := service.Update(uuid.New(), request)

	assert.Equal(t, errors.New("error"), err)

	paymentSchedulerRepository.AssertNotCalled(t, "DeleteById", mock.Anything)
}

func Test_Update_WithErrorFromScheduler(t *testing.T) {
	service := serviceGenerator()

	request := mocks.GenerateUpdatePaymentSchedulerRequest()
	providerId := uuid.New()

	paymentSchedulerRepository.On("ExistsById", mock.Anything).Return(true)
	providerService.On("ExistsById", request.ProviderId).Return(true)
	paymentSchedulerRepository.On("Update", mock.Anything).Return(nil)
	serviceScheduler.On("Update", mock.AnythingOfType("uuid.UUID"), string(request.Spec), mock.Anything).
		Return(cron.EntryID(0), errors.New("error"))
	paymentSchedulerRepository.On("DeleteById", providerId).Return()

	entity := request.ToEntity(providerId)
	entity.HouseId = mocks.HouseId
	entity.UserId = mocks.UserId

	paymentSchedulerRepository.On("FindById", providerId).Return(entity, nil)

	err := service.Update(providerId, request)

	assert.Equal(t, errors.New("error"), err)
}
