package service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/robfig/cron"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	paymentModel "payment/model"
	paymentScheduler "payment/scheduler/model"
	scheduler2 "scheduler"
	innerMock "test/mock"
	"test/testhelper"
	"testing"
)

var (
	users      *innerMock.UserServiceMock
	houses     *innerMock.HouseServiceMock
	payments   *innerMock.PaymentServiceMock
	schedulers *innerMock.SchedulerServiceMock
	houseId    = testhelper.ParseUUID("d077adaa-00d7-4e80-ac86-57512267505d")
	userId     = testhelper.ParseUUID("ad2c5035-6745-48d0-9eee-fd22f5dae8e0")
)

func serviceGenerator() PaymentSchedulerService {
	users = new(innerMock.UserServiceMock)
	houses = new(innerMock.HouseServiceMock)
	payments = new(innerMock.PaymentServiceMock)
	schedulers = new(innerMock.SchedulerServiceMock)

	return NewPaymentSchedulerService(users, houses, payments, schedulers)
}

func Test_Add(t *testing.T) {
	service := serviceGenerator()

	users.On("ExistsById", userId).
		Return(true)
	houses.On("ExistsById", houseId).
		Return(true)
	schedulers.On("Add", mock.AnythingOfType("uuid.UUID"), "@daily", mock.Anything).
		Return(cron.EntryID(0), nil)

	request := generateCreatePaymentSchedulerRequest()

	payment, err := service.Add(request)

	expectedEntity := request.ToEntity()
	expectedEntity.Id = payment.Id
	expectedResponse := expectedEntity.ToResponse()

	assert.Nil(t, err)
	assert.Equal(t, expectedResponse, payment)
	schedulers.AssertCalled(t, "Add", expectedEntity.Id, "@daily", mock.Anything)

	payments.On("Add", mock.Anything).Return(paymentModel.PaymentResponse{}, nil)

	function := schedulers.Calls[0].Arguments.Get(2).(func())
	function()

	createPaymentRequest := payments.Calls[0].Arguments.Get(0).(paymentModel.CreatePaymentRequest)

	assert.Equal(t, paymentModel.CreatePaymentRequest{
		Name:        "Test Payment",
		Description: "Test Payment Description",
		HouseId:     houseId,
		UserId:      userId,
		Date:        createPaymentRequest.Date,
		Sum:         1000,
	}, createPaymentRequest)

	serviceObject := service.(*paymentSchedulerServiceObject)

	_, paymentExists := serviceObject.payments[payment.Id]
	assert.True(t, paymentExists)

	_, housePaymentExists := serviceObject.housePayments[payment.HouseId]
	assert.True(t, housePaymentExists)

	_, userPaymentExists := serviceObject.userPayments[payment.UserId]
	assert.True(t, userPaymentExists)
}

func Test_Add_WithNegativeSum(t *testing.T) {
	service := serviceGenerator()

	request := generateCreatePaymentSchedulerRequest()
	request.Sum = -1000

	payment, err := service.Add(request)

	assert.Equal(t, errors.New("sum should not be zero of negative"), err)
	assert.Equal(t, paymentScheduler.PaymentSchedulerResponse{}, payment)
	assert.Len(t, service.(*paymentSchedulerServiceObject).payments, 0)
	assert.Len(t, service.(*paymentSchedulerServiceObject).housePayments, 0)
	assert.Len(t, service.(*paymentSchedulerServiceObject).userPayments, 0)
}

func Test_Add_WithZeroSum(t *testing.T) {
	service := serviceGenerator()

	request := generateCreatePaymentSchedulerRequest()
	request.Sum = 0

	payment, err := service.Add(request)

	assert.Equal(t, errors.New("sum should not be zero of negative"), err)
	assert.Equal(t, paymentScheduler.PaymentSchedulerResponse{}, payment)
	assert.Len(t, service.(*paymentSchedulerServiceObject).payments, 0)
	assert.Len(t, service.(*paymentSchedulerServiceObject).housePayments, 0)
	assert.Len(t, service.(*paymentSchedulerServiceObject).userPayments, 0)
}

func Test_Add_WithUserNotExists(t *testing.T) {
	service := serviceGenerator()

	users.On("ExistsById", userId).Return(false)

	request := generateCreatePaymentSchedulerRequest()

	payment, err := service.Add(request)

	assert.Equal(t, errors.New(fmt.Sprintf("user with id %s in not exists", request.UserId)), err)
	assert.Equal(t, paymentScheduler.PaymentSchedulerResponse{}, payment)
	assert.Len(t, service.(*paymentSchedulerServiceObject).payments, 0)
	assert.Len(t, service.(*paymentSchedulerServiceObject).housePayments, 0)
	assert.Len(t, service.(*paymentSchedulerServiceObject).userPayments, 0)
}

func Test_Add_WithHouseNotExists(t *testing.T) {
	service := serviceGenerator()

	users.On("ExistsById", userId).Return(true)
	houses.On("ExistsById", houseId).Return(false)

	request := generateCreatePaymentSchedulerRequest()

	payment, err := service.Add(request)

	assert.Equal(t, errors.New(fmt.Sprintf("house with id %s in not exists", request.HouseId)), err)
	assert.Equal(t, paymentScheduler.PaymentSchedulerResponse{}, payment)
	assert.Len(t, service.(*paymentSchedulerServiceObject).payments, 0)
	assert.Len(t, service.(*paymentSchedulerServiceObject).housePayments, 0)
	assert.Len(t, service.(*paymentSchedulerServiceObject).userPayments, 0)
}

func Test_Add_WithInvalidSpec(t *testing.T) {
	service := serviceGenerator()

	users.On("ExistsById", userId).
		Return(true)
	houses.On("ExistsById", houseId).
		Return(true)
	schedulers.On("Add", mock.AnythingOfType("uuid.UUID"), "@daily", mock.Anything).
		Return(cron.EntryID(0), nil)

	request := generateCreatePaymentSchedulerRequest()
	request.Spec = ""

	payment, err := service.Add(request)

	assert.Equal(t, errors.New("scheduler configuration not provided"), err)
	assert.Equal(t, paymentScheduler.PaymentSchedulerResponse{}, payment)
	assert.Len(t, service.(*paymentSchedulerServiceObject).payments, 0)
	assert.Len(t, service.(*paymentSchedulerServiceObject).housePayments, 0)
	assert.Len(t, service.(*paymentSchedulerServiceObject).userPayments, 0)
}

func Test_Add_WithErrorDuringScheduling(t *testing.T) {
	service := serviceGenerator()

	users.On("ExistsById", userId).
		Return(true)
	houses.On("ExistsById", houseId).
		Return(true)
	schedulers.On("Add", mock.AnythingOfType("uuid.UUID"), "@daily", mock.Anything).
		Return(cron.EntryID(0), errors.New("error"))

	request := generateCreatePaymentSchedulerRequest()

	payment, err := service.Add(request)

	assert.Equal(t, errors.New("error"), err)
	assert.Equal(t, paymentScheduler.PaymentSchedulerResponse{}, payment)
	assert.Len(t, service.(*paymentSchedulerServiceObject).payments, 0)
	assert.Len(t, service.(*paymentSchedulerServiceObject).housePayments, 0)
	assert.Len(t, service.(*paymentSchedulerServiceObject).userPayments, 0)
}

func Test_Remove(t *testing.T) {
	service := serviceGenerator()

	scheduler := generatePaymentScheduler()

	service.(*paymentSchedulerServiceObject).payments[scheduler.Id] = scheduler
	service.(*paymentSchedulerServiceObject).housePayments[scheduler.HouseId] = []paymentScheduler.PaymentScheduler{scheduler}
	service.(*paymentSchedulerServiceObject).userPayments[scheduler.UserId] = []paymentScheduler.PaymentScheduler{scheduler}

	schedulers.On("Remove", scheduler.Id).Return(nil)

	err := service.Remove(scheduler.Id)

	assert.Nil(t, err)

	assert.Len(t, service.(*paymentSchedulerServiceObject).payments, 0)
	assert.Len(t, service.(*paymentSchedulerServiceObject).housePayments, 0)
	assert.Len(t, service.(*paymentSchedulerServiceObject).userPayments, 0)
}

func Test_Remove_WithErrorFromScheduler(t *testing.T) {
	service := serviceGenerator()

	scheduler := generatePaymentScheduler()

	service.(*paymentSchedulerServiceObject).payments[scheduler.Id] = scheduler
	service.(*paymentSchedulerServiceObject).housePayments[scheduler.HouseId] = []paymentScheduler.PaymentScheduler{scheduler}
	service.(*paymentSchedulerServiceObject).userPayments[scheduler.UserId] = []paymentScheduler.PaymentScheduler{scheduler}

	schedulers.On("Remove", scheduler.Id).Return(errors.New("error"))

	err := service.Remove(scheduler.Id)

	assert.Nil(t, err)

	assert.Len(t, service.(*paymentSchedulerServiceObject).payments, 0)
	assert.Len(t, service.(*paymentSchedulerServiceObject).housePayments, 0)
	assert.Len(t, service.(*paymentSchedulerServiceObject).userPayments, 0)
}

func Test_Remove_WithMissingRecord(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()

	err := service.Remove(id)

	assert.Equal(t, errors.New(fmt.Sprintf("payment scheduler with id %s not found", id)), err)
}

func Test_FindById(t *testing.T) {
	service := serviceGenerator()

	scheduler := generatePaymentScheduler()

	service.(*paymentSchedulerServiceObject).payments[scheduler.Id] = scheduler

	actual, err := service.FindById(scheduler.Id)

	assert.Nil(t, err)
	assert.Equal(t, scheduler.ToResponse(), actual)
}

func Test_FindById_WithNotExistingId(t *testing.T) {
	paymentService := serviceGenerator()

	id := uuid.New()
	actual, err := paymentService.FindById(id)

	assert.Equal(t, errors.New(fmt.Sprintf("payment scheduler with id %s not found", id)), err)
	assert.Equal(t, paymentScheduler.PaymentSchedulerResponse{}, actual)
}

func Test_FindByHouseId(t *testing.T) {
	service := serviceGenerator()

	scheduler := generatePaymentScheduler()

	service.(*paymentSchedulerServiceObject).housePayments[scheduler.HouseId] = []paymentScheduler.PaymentScheduler{scheduler}

	actual := service.FindByHouseId(scheduler.HouseId)

	assert.Equal(t, []paymentScheduler.PaymentSchedulerResponse{scheduler.ToResponse()}, actual)
}

func Test_FindByHouseId_WithNotExistingRecords(t *testing.T) {
	paymentService := serviceGenerator()

	actual := paymentService.FindByHouseId(uuid.New())

	assert.Equal(t, []paymentScheduler.PaymentSchedulerResponse{}, actual)
}

func Test_FindByUserId(t *testing.T) {
	service := serviceGenerator()

	scheduler := generatePaymentScheduler()

	service.(*paymentSchedulerServiceObject).userPayments[scheduler.UserId] = []paymentScheduler.PaymentScheduler{scheduler}

	actual := service.FindByUserId(scheduler.UserId)

	assert.Equal(t, []paymentScheduler.PaymentSchedulerResponse{scheduler.ToResponse()}, actual)
}

func Test_FindPaymentByUserId_WithNotExistingRecords(t *testing.T) {
	service := serviceGenerator()

	actual := service.FindByUserId(uuid.New())

	assert.Equal(t, []paymentScheduler.PaymentSchedulerResponse{}, actual)
}

func generateCreatePaymentSchedulerRequest() paymentScheduler.CreatePaymentSchedulerRequest {
	return paymentScheduler.CreatePaymentSchedulerRequest{
		Name:        "Test Payment",
		Description: "Test Payment Description",
		HouseId:     houseId,
		UserId:      userId,
		Sum:         1000,
		Spec:        scheduler2.DAILY,
	}
}

func generatePaymentScheduler() paymentScheduler.PaymentScheduler {
	return paymentScheduler.PaymentScheduler{
		Payment: paymentModel.Payment{
			Id:          uuid.New(),
			Name:        "Test Payment",
			Description: "Test Payment Description",
			HouseId:     houseId,
			UserId:      userId,
			Sum:         1000,
		},
		Spec: scheduler2.DAILY,
	}
}
