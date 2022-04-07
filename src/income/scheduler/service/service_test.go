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
	"github.com/VlasovArtem/hob/src/test/testhelper"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type IncomeSchedulerServiceTestSuite struct {
	testhelper.MockTestSuite[IncomeSchedulerService]
	houses              *houseMocks.HouseService
	incomes             *incomeMocks.IncomeService
	schedulers          *schedulerMocks.ServiceScheduler
	schedulerRepository *mocks.IncomeSchedulerRepository
}

func TestIncomeSchedulerServiceTestSuite(t *testing.T) {
	ts := &IncomeSchedulerServiceTestSuite{}
	ts.TestObjectGenerator = func() IncomeSchedulerService {
		ts.houses = new(houseMocks.HouseService)
		ts.incomes = new(incomeMocks.IncomeService)
		ts.schedulers = new(schedulerMocks.ServiceScheduler)
		ts.schedulerRepository = new(mocks.IncomeSchedulerRepository)
		return NewIncomeSchedulerService(ts.houses, ts.incomes, ts.schedulers, ts.schedulerRepository)
	}

	suite.Run(t, ts)
}

func (i *IncomeSchedulerServiceTestSuite) Test_Add() {
	request := mocks.GenerateCreateIncomeSchedulerRequest()

	i.houses.On("ExistsById", request.HouseId).Return(true)
	i.schedulers.On("Add", mock.AnythingOfType("uuid.UUID"), "@daily", mock.Anything).Return(cron.EntryID(0), nil)
	i.schedulerRepository.On("Create", mock.Anything).Return(
		func(meter model.IncomeScheduler) model.IncomeScheduler {
			return meter
		},
		nil,
	)

	payment, err := i.TestO.Add(request)

	expectedEntity := request.ToEntity()
	expectedEntity.Id = payment.Id
	expectedResponse := expectedEntity.ToDto()

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), expectedResponse, payment)
	i.schedulers.AssertCalled(i.T(), "Add", expectedEntity.Id, "@daily", mock.Anything)

	i.incomes.On("Add", mock.Anything).Return(incomeModel.IncomeDto{}, nil)

	function := i.schedulers.Calls[0].Arguments.Get(2).(func())
	function()

	createIncomeRequest := i.incomes.Calls[0].Arguments.Get(0).(incomeModel.CreateIncomeRequest)

	assert.Equal(i.T(), incomeModel.CreateIncomeRequest{
		Name:        "Test Income",
		Description: "Test Income Description",
		HouseId:     createIncomeRequest.HouseId,
		Date:        createIncomeRequest.Date,
		Sum:         1000,
	}, createIncomeRequest)
}

func (i *IncomeSchedulerServiceTestSuite) Test_Add_WithHouseNotExists() {
	request := mocks.GenerateCreateIncomeSchedulerRequest()

	i.houses.On("ExistsById", request.HouseId).Return(false)

	payment, err := i.TestO.Add(request)

	assert.Equal(i.T(), int_errors.NewErrNotFound("house with id %s not found", request.HouseId), err)
	assert.Equal(i.T(), model.IncomeSchedulerDto{}, payment)
	i.schedulers.AssertNotCalled(i.T(), "Create", mock.Anything, "@daily", mock.Anything)
}

func (i *IncomeSchedulerServiceTestSuite) Test_Add_WithInvalidSpec() {
	request := mocks.GenerateCreateIncomeSchedulerRequest()

	i.houses.On("ExistsById", request.HouseId).
		Return(true)
	i.schedulers.On("Create", mock.AnythingOfType("uuid.UUID"), "@daily", mock.Anything).
		Return(cron.EntryID(0), nil)

	request.Spec = ""

	payment, err := i.TestO.Add(request)

	assert.Equal(i.T(), errors.New("scheduler configuration not provided"), err)
	assert.Equal(i.T(), model.IncomeSchedulerDto{}, payment)
}

func (i *IncomeSchedulerServiceTestSuite) Test_Add_WithErrorDuringScheduling() {
	request := mocks.GenerateCreateIncomeSchedulerRequest()

	i.houses.On("ExistsById", request.HouseId).
		Return(true)
	i.schedulers.On("Add", mock.AnythingOfType("uuid.UUID"), "@daily", mock.Anything).
		Return(cron.EntryID(0), errors.New("error"))

	payment, err := i.TestO.Add(request)

	assert.Equal(i.T(), errors.New("error"), err)
	assert.Equal(i.T(), model.IncomeSchedulerDto{}, payment)
}

func (i *IncomeSchedulerServiceTestSuite) Test_DeleteById() {

	id := uuid.New()

	i.schedulerRepository.On("ExistsById", id).Return(true)
	i.schedulers.On("Remove", id).Return(nil)
	i.schedulerRepository.On("DeleteById", id).Return(nil)

	err := i.TestO.DeleteById(id)

	assert.Nil(i.T(), err)

	i.schedulerRepository.AssertCalled(i.T(), "ExistsById", id)
	i.schedulers.AssertCalled(i.T(), "Remove", id)
	i.schedulerRepository.AssertCalled(i.T(), "DeleteById", id)
}

func (i *IncomeSchedulerServiceTestSuite) Test_DeleteById_WithNotExistsScheduler() {

	id := uuid.New()

	i.schedulerRepository.On("ExistsById", id).Return(false)

	err := i.TestO.DeleteById(id)

	assert.Equal(i.T(), int_errors.NewErrNotFound("income scheduler with id %s not found", id), err)

	i.schedulerRepository.AssertCalled(i.T(), "ExistsById", id)
	i.schedulers.AssertNotCalled(i.T(), "Remove", id)
	i.schedulerRepository.AssertNotCalled(i.T(), "DeleteById", id)
}

func (i *IncomeSchedulerServiceTestSuite) Test_DeleteById_WithErrorDuringDeleteByIdScheduler() {

	id := uuid.New()
	expectedError := errors.New("test")

	i.schedulerRepository.On("ExistsById", id).Return(true)
	i.schedulers.On("Remove", id).Return(expectedError)
	i.schedulerRepository.On("DeleteById", id).Return(nil)

	err := i.TestO.DeleteById(id)

	assert.Nil(i.T(), err)

	i.schedulerRepository.AssertCalled(i.T(), "ExistsById", id)
	i.schedulers.AssertCalled(i.T(), "Remove", id)
	i.schedulerRepository.AssertCalled(i.T(), "DeleteById", id)
}

func (i *IncomeSchedulerServiceTestSuite) Test_FindById() {

	scheduler := mocks.GenerateIncomeScheduler(uuid.New())

	i.schedulerRepository.On("ExistsById", scheduler.Id).Return(true)
	i.schedulerRepository.On("FindById", scheduler.Id).Return(scheduler, nil)

	actual, err := i.TestO.FindById(scheduler.Id)

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), scheduler.ToDto(), actual)
}

func (i *IncomeSchedulerServiceTestSuite) Test_FindById_WithNotExistingId() {
	id := uuid.New()

	i.schedulerRepository.On("ExistsById", id).Return(false)

	actual, err := i.TestO.FindById(id)

	assert.Equal(i.T(), int_errors.NewErrNotFound("income scheduler with id %s not found", id), err)
	assert.Equal(i.T(), model.IncomeSchedulerDto{}, actual)

	i.schedulerRepository.AssertNotCalled(i.T(), "FindById", id)
}

func (i *IncomeSchedulerServiceTestSuite) Test_FindById_WithErrorFromDatabase() {
	id := uuid.New()
	expectedError := errors.New("error")

	i.schedulerRepository.On("ExistsById", id).Return(true)
	i.schedulerRepository.On("FindById", id).Return(model.IncomeScheduler{}, expectedError)

	actual, err := i.TestO.FindById(id)

	assert.Equal(i.T(), expectedError, err)
	assert.Equal(i.T(), model.IncomeSchedulerDto{}, actual)

	i.schedulerRepository.AssertCalled(i.T(), "FindById", id)
}

func (i *IncomeSchedulerServiceTestSuite) Test_FindByHouseId() {

	scheduler := mocks.GenerateIncomeScheduler(uuid.New())

	i.schedulerRepository.On("FindByHouseId", *scheduler.HouseId).Return([]model.IncomeSchedulerDto{scheduler.ToDto()}, nil)

	actual := i.TestO.FindByHouseId(*scheduler.HouseId)

	assert.Equal(i.T(), []model.IncomeSchedulerDto{scheduler.ToDto()}, actual)
}

func (i *IncomeSchedulerServiceTestSuite) Test_FindByHouseId_WithNotExistingRecords() {
	id := uuid.New()
	i.schedulerRepository.On("FindByHouseId", id).Return([]model.IncomeSchedulerDto{}, errors.New("test"))

	actual := i.TestO.FindByHouseId(id)

	assert.Equal(i.T(), []model.IncomeSchedulerDto{}, actual)
}

func (i *IncomeSchedulerServiceTestSuite) Test_Update() {
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

	i.schedulerRepository.On("ExistsById", id).Return(true)
	i.schedulerRepository.On("Update", id, request).Return(scheduler, nil)
	i.schedulers.On("Update", id, string(request.Spec), mock.Anything).Return(cron.EntryID(0), nil)

	err := i.TestO.Update(id, request)

	assert.Nil(i.T(), err)

	i.schedulerRepository.AssertCalled(i.T(), "ExistsById", id)
	i.schedulerRepository.AssertCalled(i.T(), "Update", id, request)
	i.schedulers.AssertCalled(i.T(), "Update", id, string(request.Spec), mock.Anything)
	i.schedulerRepository.AssertNotCalled(i.T(), "DeleteById", id)
}

func (i *IncomeSchedulerServiceTestSuite) Test_Update_WithNotExists() {
	id, request := mocks.GenerateUpdateIncomeSchedulerRequest()

	i.schedulerRepository.On("ExistsById", id).Return(false)

	err := i.TestO.Update(id, request)

	assert.Equal(i.T(), int_errors.NewErrNotFound("income scheduler with id %s not found", id), err)

	i.schedulerRepository.AssertCalled(i.T(), "ExistsById", id)
	i.schedulerRepository.AssertNotCalled(i.T(), "Update", id, request)
	i.schedulers.AssertNotCalled(i.T(), "Update", id, string(request.Spec), mock.Anything)
	i.schedulerRepository.AssertNotCalled(i.T(), "DeleteById", id)
}

func (i *IncomeSchedulerServiceTestSuite) Test_Update_WithInvalidSum() {
	sums := []float32{-1, 0}

	for _, sum := range sums {
		id, request := mocks.GenerateUpdateIncomeSchedulerRequest()
		request.Sum = sum

		err := i.TestO.Update(id, request)

		assert.Equal(i.T(), errors.New("sum should not be zero of negative"), err)

		i.schedulerRepository.AssertNotCalled(i.T(), "ExistsById", id)
		i.schedulerRepository.AssertNotCalled(i.T(), "Update", id, request)
		i.schedulers.AssertNotCalled(i.T(), "Update", id, string(request.Spec), mock.Anything)
		i.schedulerRepository.AssertNotCalled(i.T(), "DeleteById", id)
	}
}

func (i *IncomeSchedulerServiceTestSuite) Test_Update_WithInvalidScheduler() {

	id, request := mocks.GenerateUpdateIncomeSchedulerRequest()
	request.Spec = ""

	i.schedulerRepository.On("ExistsById", id).Return(true)

	err := i.TestO.Update(id, request)

	assert.Equal(i.T(), errors.New("scheduler configuration not provided"), err)

	i.schedulerRepository.AssertCalled(i.T(), "ExistsById", id)
	i.schedulerRepository.AssertNotCalled(i.T(), "Update", id, request)
	i.schedulers.AssertNotCalled(i.T(), "Update", id, string(request.Spec), mock.Anything)
	i.schedulerRepository.AssertNotCalled(i.T(), "DeleteById", id)
}

func (i *IncomeSchedulerServiceTestSuite) Test_Update_WithErrorFromRepository() {

	id, request := mocks.GenerateUpdateIncomeSchedulerRequest()

	i.schedulerRepository.On("ExistsById", id).Return(true)
	i.schedulerRepository.On("Update", id, request).Return(model.IncomeScheduler{}, errors.New("test"))

	err := i.TestO.Update(id, request)

	assert.Equal(i.T(), errors.New("test"), err)

	i.schedulerRepository.AssertCalled(i.T(), "ExistsById", id)
	i.schedulerRepository.AssertCalled(i.T(), "Update", id, request)
	i.schedulers.AssertNotCalled(i.T(), "Update", id, string(request.Spec), mock.Anything)
	i.schedulerRepository.AssertNotCalled(i.T(), "DeleteById", id)
}

func (i *IncomeSchedulerServiceTestSuite) Test_Update_WithErrorFromScheduler() {

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

	i.schedulerRepository.On("ExistsById", id).Return(true)
	i.schedulerRepository.On("Update", id, request).Return(scheduler, nil)
	i.schedulers.On("Update", id, string(request.Spec), mock.Anything).Return(cron.EntryID(0), errors.New("test2"))
	i.schedulerRepository.On("DeleteById", id).Return(nil)

	err := i.TestO.Update(id, request)

	assert.Equal(i.T(), errors.New("test2"), err)

	i.schedulerRepository.AssertCalled(i.T(), "ExistsById", id)
	i.schedulerRepository.AssertCalled(i.T(), "Update", id, request)
	i.schedulers.AssertCalled(i.T(), "Update", id, string(request.Spec), mock.Anything)
	i.schedulerRepository.AssertCalled(i.T(), "DeleteById", id)
}
