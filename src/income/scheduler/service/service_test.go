package service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/robfig/cron"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	incomeModel "income/model"
	incomeSchedulerModel "income/scheduler/model"
	scheduler2 "scheduler"
	innerMock "test/mock"
	"test/testhelper"
	"testing"
)

var (
	houses     *innerMock.HouseServiceMock
	incomes    *innerMock.IncomeServiceMock
	schedulers *innerMock.SchedulerServiceMock
	houseId    = testhelper.ParseUUID("99c48818-ea50-4f56-8d02-44e55e3bfc32")
)

func serviceGenerator() IncomeSchedulerService {
	houses = new(innerMock.HouseServiceMock)
	incomes = new(innerMock.IncomeServiceMock)
	schedulers = new(innerMock.SchedulerServiceMock)

	return NewIncomeSchedulerService(houses, incomes, schedulers)
}

func Test_Add(t *testing.T) {
	service := serviceGenerator()

	houses.On("ExistsById", houseId).
		Return(true)
	schedulers.On("Add", mock.AnythingOfType("uuid.UUID"), "@daily", mock.Anything).
		Return(cron.EntryID(0), nil)

	request := generateCreateIncomeSchedulerRequest()

	payment, err := service.Add(request)

	expectedEntity := request.ToEntity()
	expectedEntity.Id = payment.Id
	expectedResponse := expectedEntity.ToResponse()

	assert.Nil(t, err)
	assert.Equal(t, expectedResponse, payment)
	schedulers.AssertCalled(t, "Add", expectedEntity.Id, "@daily", mock.Anything)

	incomes.On("Add", mock.Anything).Return(incomeModel.IncomeResponse{}, nil)

	function := schedulers.Calls[0].Arguments.Get(2).(func())
	function()

	createIncomeRequest := incomes.Calls[0].Arguments.Get(0).(incomeModel.CreateIncomeRequest)

	assert.Equal(t, incomeModel.CreateIncomeRequest{
		Name:        "Test Income",
		Description: "Test Income Description",
		HouseId:     houseId,
		Date:        createIncomeRequest.Date,
		Sum:         1000,
	}, createIncomeRequest)

	serviceObject := service.(*incomeSchedulerServiceObject)

	_, paymentExists := serviceObject.incomeSchedulers[payment.Id]
	assert.True(t, paymentExists)

	_, housePaymentExists := serviceObject.houseIncomeSchedulers[payment.HouseId]
	assert.True(t, housePaymentExists)
}

func Test_Add_WithHouseNotExists(t *testing.T) {
	service := serviceGenerator()

	houses.On("ExistsById", houseId).Return(false)

	request := generateCreateIncomeSchedulerRequest()

	payment, err := service.Add(request)

	assert.Equal(t, errors.New(fmt.Sprintf("house with id %s not found", request.HouseId)), err)
	assert.Equal(t, incomeSchedulerModel.IncomeSchedulerResponse{}, payment)
	assert.Len(t, service.(*incomeSchedulerServiceObject).incomeSchedulers, 0)
	assert.Len(t, service.(*incomeSchedulerServiceObject).houseIncomeSchedulers, 0)
}

func Test_Add_WithInvalidSpec(t *testing.T) {
	service := serviceGenerator()

	houses.On("ExistsById", houseId).
		Return(true)
	schedulers.On("Add", mock.AnythingOfType("uuid.UUID"), "@daily", mock.Anything).
		Return(cron.EntryID(0), nil)

	request := generateCreateIncomeSchedulerRequest()
	request.Spec = ""

	payment, err := service.Add(request)

	assert.Equal(t, errors.New("scheduler configuration not provided"), err)
	assert.Equal(t, incomeSchedulerModel.IncomeSchedulerResponse{}, payment)
	assert.Len(t, service.(*incomeSchedulerServiceObject).incomeSchedulers, 0)
	assert.Len(t, service.(*incomeSchedulerServiceObject).houseIncomeSchedulers, 0)
}

func Test_Add_WithErrorDuringScheduling(t *testing.T) {
	service := serviceGenerator()

	houses.On("ExistsById", houseId).
		Return(true)
	schedulers.On("Add", mock.AnythingOfType("uuid.UUID"), "@daily", mock.Anything).
		Return(cron.EntryID(0), errors.New("error"))

	request := generateCreateIncomeSchedulerRequest()

	payment, err := service.Add(request)

	assert.Equal(t, errors.New("error"), err)
	assert.Equal(t, incomeSchedulerModel.IncomeSchedulerResponse{}, payment)
	assert.Len(t, service.(*incomeSchedulerServiceObject).incomeSchedulers, 0)
	assert.Len(t, service.(*incomeSchedulerServiceObject).houseIncomeSchedulers, 0)
}

func Test_Remove(t *testing.T) {
	service := serviceGenerator()

	scheduler := generatePaymentScheduler()

	service.(*incomeSchedulerServiceObject).incomeSchedulers[scheduler.Id] = scheduler
	service.(*incomeSchedulerServiceObject).houseIncomeSchedulers[scheduler.HouseId] = scheduler

	schedulers.On("Remove", scheduler.Id).Return(nil)

	err := service.Remove(scheduler.Id)

	assert.Nil(t, err)

	assert.Len(t, service.(*incomeSchedulerServiceObject).incomeSchedulers, 0)
	assert.Len(t, service.(*incomeSchedulerServiceObject).houseIncomeSchedulers, 0)
}

func Test_Remove_WithErrorFromScheduler(t *testing.T) {
	service := serviceGenerator()

	scheduler := generatePaymentScheduler()

	service.(*incomeSchedulerServiceObject).incomeSchedulers[scheduler.Id] = scheduler
	service.(*incomeSchedulerServiceObject).houseIncomeSchedulers[scheduler.HouseId] = scheduler

	schedulers.On("Remove", scheduler.Id).Return(errors.New("error"))

	err := service.Remove(scheduler.Id)

	assert.Nil(t, err)

	assert.Len(t, service.(*incomeSchedulerServiceObject).incomeSchedulers, 0)
	assert.Len(t, service.(*incomeSchedulerServiceObject).houseIncomeSchedulers, 0)
}

func Test_Remove_WithMissingRecord(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()

	err := service.Remove(id)

	assert.Equal(t, errors.New(fmt.Sprintf("income scheduler with id %s not found", id)), err)
}

func Test_FindById(t *testing.T) {
	service := serviceGenerator()

	scheduler := generatePaymentScheduler()

	service.(*incomeSchedulerServiceObject).incomeSchedulers[scheduler.Id] = scheduler

	actual, err := service.FindById(scheduler.Id)

	assert.Nil(t, err)
	assert.Equal(t, scheduler.ToResponse(), actual)
}

func Test_FindById_WithNotExistingId(t *testing.T) {
	paymentService := serviceGenerator()

	id := uuid.New()
	actual, err := paymentService.FindById(id)

	assert.Equal(t, errors.New(fmt.Sprintf("income scheduler with id %s not found", id)), err)
	assert.Equal(t, incomeSchedulerModel.IncomeSchedulerResponse{}, actual)
}

func Test_FindByHouseId(t *testing.T) {
	service := serviceGenerator()

	scheduler := generatePaymentScheduler()

	service.(*incomeSchedulerServiceObject).houseIncomeSchedulers[scheduler.HouseId] = scheduler

	actual, err := service.FindByHouseId(scheduler.HouseId)

	assert.Nil(t, err)
	assert.Equal(t, scheduler.ToResponse(), actual)
}

func Test_FindByHouseId_WithNotExistingRecords(t *testing.T) {
	paymentService := serviceGenerator()

	id := uuid.New()
	actual, err := paymentService.FindByHouseId(id)

	assert.Equal(t, errors.New(fmt.Sprintf("income scheduler with house id %s not found", id)), err)
	assert.Equal(t, incomeSchedulerModel.IncomeSchedulerResponse{}, actual)
}

func generateCreateIncomeSchedulerRequest() incomeSchedulerModel.CreateIncomeSchedulerRequest {
	return incomeSchedulerModel.CreateIncomeSchedulerRequest{
		Name:        "Test Income",
		Description: "Test Income Description",
		HouseId:     houseId,
		Sum:         1000,
		Spec:        scheduler2.DAILY,
	}
}

func generatePaymentScheduler() incomeSchedulerModel.IncomeScheduler {
	return incomeSchedulerModel.IncomeScheduler{
		Income: incomeModel.Income{
			Id:          uuid.New(),
			Name:        "Test Income",
			Description: "Test Income Description",
			HouseId:     houseId,
			Sum:         1000,
		},
		Spec: scheduler2.DAILY,
	}
}
