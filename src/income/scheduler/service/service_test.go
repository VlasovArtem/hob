package service

import (
	"errors"
	int_errors "github.com/VlasovArtem/hob/src/common/int-errors"
	houseMocks "github.com/VlasovArtem/hob/src/house/mocks"
	incomeMocks "github.com/VlasovArtem/hob/src/income/mocks"
	incomeModel "github.com/VlasovArtem/hob/src/income/model"
	"github.com/VlasovArtem/hob/src/income/scheduler/mocks"
	"github.com/VlasovArtem/hob/src/income/scheduler/model"
	schedulerMocks "github.com/VlasovArtem/hob/src/scheduler/mocks"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

var (
	houses              *houseMocks.HouseService
	incomes             *incomeMocks.IncomeService
	schedulers          *schedulerMocks.ServiceScheduler
	schedulerRepository *mocks.IncomeSchedulerRepository
)

func serviceGenerator() IncomeSchedulerService {
	houses = new(houseMocks.HouseService)
	incomes = new(incomeMocks.IncomeService)
	schedulers = new(schedulerMocks.ServiceScheduler)
	schedulerRepository = new(mocks.IncomeSchedulerRepository)

	return NewIncomeSchedulerService(houses, incomes, schedulers, schedulerRepository)
}

func Test_Add(t *testing.T) {
	service := serviceGenerator()
	request := mocks.GenerateCreateIncomeSchedulerRequest()

	houses.On("ExistsById", request.HouseId).Return(true)
	schedulers.On("Add", mock.AnythingOfType("uuid.UUID"), "@daily", mock.Anything).Return(cron.EntryID(0), nil)
	schedulerRepository.On("Create", mock.Anything).Return(
		func(meter model.IncomeScheduler) model.IncomeScheduler {
			return meter
		},
		nil,
	)

	payment, err := service.Add(request)

	expectedEntity := request.ToEntity()
	expectedEntity.Id = payment.Id
	expectedResponse := expectedEntity.ToDto()

	assert.Nil(t, err)
	assert.Equal(t, expectedResponse, payment)
	schedulers.AssertCalled(t, "Add", expectedEntity.Id, "@daily", mock.Anything)

	incomes.On("Add", mock.Anything).Return(incomeModel.IncomeDto{}, nil)

	function := schedulers.Calls[0].Arguments.Get(2).(func())
	function()

	createIncomeRequest := incomes.Calls[0].Arguments.Get(0).(incomeModel.CreateIncomeRequest)

	assert.Equal(t, incomeModel.CreateIncomeRequest{
		Name:        "Test Income",
		Description: "Test Income Description",
		HouseId:     createIncomeRequest.HouseId,
		Date:        createIncomeRequest.Date,
		Sum:         1000,
	}, createIncomeRequest)
}

func Test_Add_WithHouseNotExists(t *testing.T) {
	service := serviceGenerator()
	request := mocks.GenerateCreateIncomeSchedulerRequest()

	houses.On("ExistsById", request.HouseId).Return(false)

	payment, err := service.Add(request)

	assert.Equal(t, int_errors.NewErrNotFound("house with id %s not found", request.HouseId), err)
	assert.Equal(t, model.IncomeSchedulerDto{}, payment)
	schedulers.AssertNotCalled(t, "Create", mock.Anything, "@daily", mock.Anything)
}

func Test_Add_WithInvalidSpec(t *testing.T) {
	service := serviceGenerator()
	request := mocks.GenerateCreateIncomeSchedulerRequest()

	houses.On("ExistsById", request.HouseId).
		Return(true)
	schedulers.On("Create", mock.AnythingOfType("uuid.UUID"), "@daily", mock.Anything).
		Return(cron.EntryID(0), nil)

	request.Spec = ""

	payment, err := service.Add(request)

	assert.Equal(t, errors.New("scheduler configuration not provided"), err)
	assert.Equal(t, model.IncomeSchedulerDto{}, payment)
}

func Test_Add_WithErrorDuringScheduling(t *testing.T) {
	service := serviceGenerator()
	request := mocks.GenerateCreateIncomeSchedulerRequest()

	houses.On("ExistsById", request.HouseId).
		Return(true)
	schedulers.On("Add", mock.AnythingOfType("uuid.UUID"), "@daily", mock.Anything).
		Return(cron.EntryID(0), errors.New("error"))

	payment, err := service.Add(request)

	assert.Equal(t, errors.New("error"), err)
	assert.Equal(t, model.IncomeSchedulerDto{}, payment)
}

func Test_DeleteById(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()

	schedulerRepository.On("ExistsById", id).Return(true)
	schedulers.On("Remove", id).Return(nil)
	schedulerRepository.On("DeleteById", id).Return(nil)

	err := service.DeleteById(id)

	assert.Nil(t, err)

	schedulerRepository.AssertCalled(t, "ExistsById", id)
	schedulers.AssertCalled(t, "Remove", id)
	schedulerRepository.AssertCalled(t, "DeleteById", id)
}

func Test_DeleteById_WithNotExistsScheduler(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()

	schedulerRepository.On("ExistsById", id).Return(false)

	err := service.DeleteById(id)

	assert.Equal(t, int_errors.NewErrNotFound("income scheduler with id %s not found", id), err)

	schedulerRepository.AssertCalled(t, "ExistsById", id)
	schedulers.AssertNotCalled(t, "Remove", id)
	schedulerRepository.AssertNotCalled(t, "DeleteById", id)
}

func Test_DeleteById_WithErrorDuringDeleteByIdScheduler(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()
	expectedError := errors.New("test")

	schedulerRepository.On("ExistsById", id).Return(true)
	schedulers.On("Remove", id).Return(expectedError)
	schedulerRepository.On("DeleteById", id).Return(nil)

	err := service.DeleteById(id)

	assert.Nil(t, err)

	schedulerRepository.AssertCalled(t, "ExistsById", id)
	schedulers.AssertCalled(t, "Remove", id)
	schedulerRepository.AssertCalled(t, "DeleteById", id)
}

func Test_FindById(t *testing.T) {
	service := serviceGenerator()

	scheduler := mocks.GenerateIncomeScheduler(uuid.New())

	schedulerRepository.On("ExistsById", scheduler.Id).Return(true)
	schedulerRepository.On("FindById", scheduler.Id).Return(scheduler, nil)

	actual, err := service.FindById(scheduler.Id)

	assert.Nil(t, err)
	assert.Equal(t, scheduler.ToDto(), actual)
}

func Test_FindById_WithNotExistingId(t *testing.T) {
	paymentService := serviceGenerator()

	id := uuid.New()

	schedulerRepository.On("ExistsById", id).Return(false)

	actual, err := paymentService.FindById(id)

	assert.Equal(t, int_errors.NewErrNotFound("income scheduler with id %s not found", id), err)
	assert.Equal(t, model.IncomeSchedulerDto{}, actual)

	schedulerRepository.AssertNotCalled(t, "FindById", id)
}

func Test_FindById_WithErrorFromDatabase(t *testing.T) {
	paymentService := serviceGenerator()

	id := uuid.New()
	expectedError := errors.New("error")

	schedulerRepository.On("ExistsById", id).Return(true)
	schedulerRepository.On("FindById", id).Return(model.IncomeScheduler{}, expectedError)

	actual, err := paymentService.FindById(id)

	assert.Equal(t, expectedError, err)
	assert.Equal(t, model.IncomeSchedulerDto{}, actual)

	schedulerRepository.AssertCalled(t, "FindById", id)
}

func Test_FindByHouseId(t *testing.T) {
	service := serviceGenerator()

	scheduler := mocks.GenerateIncomeScheduler(uuid.New())

	schedulerRepository.On("FindByHouseId", scheduler.HouseId).Return([]model.IncomeSchedulerDto{scheduler.ToDto()}, nil)

	actual := service.FindByHouseId(scheduler.HouseId)

	assert.Equal(t, []model.IncomeSchedulerDto{scheduler.ToDto()}, actual)
}

func Test_FindByHouseId_WithNotExistingRecords(t *testing.T) {
	paymentService := serviceGenerator()

	id := uuid.New()
	schedulerRepository.On("FindByHouseId", id).Return([]model.IncomeSchedulerDto{}, errors.New("test"))

	actual := paymentService.FindByHouseId(id)

	assert.Equal(t, []model.IncomeSchedulerDto{}, actual)
}

func Test_Update(t *testing.T) {
	service := serviceGenerator()
	id, request := mocks.GenerateUpdateIncomeSchedulerRequest()
	scheduler := model.IncomeScheduler{
		Income: incomeModel.Income{
			Id:          id,
			Name:        request.Name,
			Description: request.Description,
			Sum:         request.Sum,
		},
		Spec: request.Spec,
	}

	schedulerRepository.On("ExistsById", id).Return(true)
	schedulerRepository.On("Update", id, request).Return(scheduler, nil)
	schedulers.On("Update", id, string(request.Spec), mock.Anything).Return(cron.EntryID(0), nil)

	err := service.Update(id, request)

	assert.Nil(t, err)

	schedulerRepository.AssertCalled(t, "ExistsById", id)
	schedulerRepository.AssertCalled(t, "Update", id, request)
	schedulers.AssertCalled(t, "Update", id, string(request.Spec), mock.Anything)
	schedulerRepository.AssertNotCalled(t, "DeleteById", id)
}

func Test_Update_WithNotExists(t *testing.T) {
	service := serviceGenerator()
	id, request := mocks.GenerateUpdateIncomeSchedulerRequest()

	schedulerRepository.On("ExistsById", id).Return(false)

	err := service.Update(id, request)

	assert.Equal(t, int_errors.NewErrNotFound("income scheduler with id %s not found", id), err)

	schedulerRepository.AssertCalled(t, "ExistsById", id)
	schedulerRepository.AssertNotCalled(t, "Update", id, request)
	schedulers.AssertNotCalled(t, "Update", id, string(request.Spec), mock.Anything)
	schedulerRepository.AssertNotCalled(t, "DeleteById", id)
}

func Test_Update_WithInvalidSum(t *testing.T) {
	sums := []float32{-1, 0}

	for _, sum := range sums {
		service := serviceGenerator()
		id, request := mocks.GenerateUpdateIncomeSchedulerRequest()
		request.Sum = sum

		err := service.Update(id, request)

		assert.Equal(t, errors.New("sum should not be zero of negative"), err)

		schedulerRepository.AssertNotCalled(t, "ExistsById", id)
		schedulerRepository.AssertNotCalled(t, "Update", id, request)
		schedulers.AssertNotCalled(t, "Update", id, string(request.Spec), mock.Anything)
		schedulerRepository.AssertNotCalled(t, "DeleteById", id)
	}
}

func Test_Update_WithInvalidScheduler(t *testing.T) {
	service := serviceGenerator()

	id, request := mocks.GenerateUpdateIncomeSchedulerRequest()
	request.Spec = ""

	schedulerRepository.On("ExistsById", id).Return(true)

	err := service.Update(id, request)

	assert.Equal(t, errors.New("scheduler configuration not provided"), err)

	schedulerRepository.AssertCalled(t, "ExistsById", id)
	schedulerRepository.AssertNotCalled(t, "Update", id, request)
	schedulers.AssertNotCalled(t, "Update", id, string(request.Spec), mock.Anything)
	schedulerRepository.AssertNotCalled(t, "DeleteById", id)
}

func Test_Update_WithErrorFromRepository(t *testing.T) {
	service := serviceGenerator()

	id, request := mocks.GenerateUpdateIncomeSchedulerRequest()

	schedulerRepository.On("ExistsById", id).Return(true)
	schedulerRepository.On("Update", id, request).Return(model.IncomeScheduler{}, errors.New("test"))

	err := service.Update(id, request)

	assert.Equal(t, errors.New("test"), err)

	schedulerRepository.AssertCalled(t, "ExistsById", id)
	schedulerRepository.AssertCalled(t, "Update", id, request)
	schedulers.AssertNotCalled(t, "Update", id, string(request.Spec), mock.Anything)
	schedulerRepository.AssertNotCalled(t, "DeleteById", id)
}

func Test_Update_WithErrorFromScheduler(t *testing.T) {
	service := serviceGenerator()

	id, request := mocks.GenerateUpdateIncomeSchedulerRequest()
	scheduler := model.IncomeScheduler{
		Income: incomeModel.Income{
			Id:          id,
			Name:        request.Name,
			Description: request.Description,
			Sum:         request.Sum,
		},
		Spec: request.Spec,
	}

	schedulerRepository.On("ExistsById", id).Return(true)
	schedulerRepository.On("Update", id, request).Return(scheduler, nil)
	schedulers.On("Update", id, string(request.Spec), mock.Anything).Return(cron.EntryID(0), errors.New("test2"))
	schedulerRepository.On("DeleteById", id).Return(nil)

	err := service.Update(id, request)

	assert.Equal(t, errors.New("test2"), err)

	schedulerRepository.AssertCalled(t, "ExistsById", id)
	schedulerRepository.AssertCalled(t, "Update", id, request)
	schedulers.AssertCalled(t, "Update", id, string(request.Spec), mock.Anything)
	schedulerRepository.AssertCalled(t, "DeleteById", id)
}
