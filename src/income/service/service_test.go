package service

import (
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/common"
	"github.com/VlasovArtem/hob/src/common/int-errors"
	groupMocks "github.com/VlasovArtem/hob/src/group/mocks"
	houseMocks "github.com/VlasovArtem/hob/src/house/mocks"
	"github.com/VlasovArtem/hob/src/income/mocks"
	"github.com/VlasovArtem/hob/src/income/model"
	"github.com/VlasovArtem/hob/src/test/testhelper"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
	"time"
)

type IncomeServiceTestSuite struct {
	testhelper.MockTestSuite[IncomeService]
	houses           *houseMocks.HouseService
	groups           *groupMocks.GroupService
	incomeRepository *mocks.IncomeRepository
}

func TestIncomeServiceTestSuite(t *testing.T) {
	ts := &IncomeServiceTestSuite{}
	ts.TestObjectGenerator = func() IncomeService {
		ts.houses = new(houseMocks.HouseService)
		ts.incomeRepository = new(mocks.IncomeRepository)
		ts.groups = new(groupMocks.GroupService)
		return NewIncomeService(ts.houses, ts.groups, ts.incomeRepository)
	}

	suite.Run(t, ts)
}

func (i *IncomeServiceTestSuite) Test_Add() {
	var savedIncome model.Income
	request := mocks.GenerateCreateIncomeRequest()

	i.houses.On("ExistsById", *request.HouseId).Return(true)
	i.incomeRepository.On("Create", mock.Anything).Return(func(income model.Income) model.Income {
		savedIncome = income

		return income
	}, nil)
	i.groups.On("ExistsByIds", mock.Anything).Return(true)

	income, err := i.TestO.Add(request)

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), savedIncome.ToDto(), income)
}

func (i *IncomeServiceTestSuite) Test_Add_WithoutHouseIdAndWithGroups() {
	var savedIncome model.Income
	request := mocks.GenerateCreateIncomeRequest()
	request.HouseId = nil
	request.GroupIds = []uuid.UUID{uuid.New()}

	i.incomeRepository.On("Create", mock.Anything).Return(func(income model.Income) model.Income {
		savedIncome = income

		return income
	}, nil)
	i.groups.On("ExistsByIds", mock.Anything).Return(true)

	income, err := i.TestO.Add(request)

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), savedIncome.ToDto(), income)

	i.houses.AssertNotCalled(i.T(), "ExistsById", mock.Anything)
}

func (i *IncomeServiceTestSuite) Test_Add_WithoutHouseIdAndGroups() {
	request := mocks.GenerateCreateIncomeRequest()
	request.HouseId = nil
	request.GroupIds = []uuid.UUID{}

	_, err := i.TestO.Add(request)

	assert.EqualError(i.T(), err, "houseId or groupId must be set")

	i.groups.AssertNotCalled(i.T(), "ExistsByIds", mock.Anything)
	i.houses.AssertNotCalled(i.T(), "ExistsById", mock.Anything)
	i.incomeRepository.AssertNotCalled(i.T(), "Create", mock.Anything)
}

func (i *IncomeServiceTestSuite) Test_Add_WithHouseNotExists() {
	request := mocks.GenerateCreateIncomeRequest()

	i.houses.On("ExistsById", *request.HouseId).Return(false)

	payment, err := i.TestO.Add(request)

	assert.Equal(i.T(), int_errors.NewErrNotFound("house with id %s not found", request.HouseId), err)
	assert.Equal(i.T(), model.IncomeDto{}, payment)

	i.incomeRepository.AssertNotCalled(i.T(), "Create", mock.Anything)
}

func (i *IncomeServiceTestSuite) Test_Add_WithErrorFromRepository() {
	expectedError := errors.New("error")
	request := mocks.GenerateCreateIncomeRequest()

	i.houses.On("ExistsById", *request.HouseId).Return(true)
	i.incomeRepository.On("Create", mock.Anything).Return(model.Income{}, expectedError)

	income, err := i.TestO.Add(request)

	assert.Equal(i.T(), expectedError, err)
	assert.Equal(i.T(), model.IncomeDto{}, income)
}

func (i *IncomeServiceTestSuite) Test_Add_WithDateAfterCurrentDate() {
	request := mocks.GenerateCreateIncomeRequest()
	request.Date = time.Now().Add(time.Hour)

	i.houses.On("ExistsById", *request.HouseId).Return(true)

	payment, err := i.TestO.Add(request)

	assert.Equal(i.T(), errors.New("date should not be after current date"), err)
	assert.Equal(i.T(), model.IncomeDto{}, payment)

	i.incomeRepository.AssertNotCalled(i.T(), "Create", mock.Anything)
}

func (i *IncomeServiceTestSuite) Test_Add_WithGroupsNotFound() {
	request := mocks.GenerateCreateIncomeRequest()
	request.GroupIds = []uuid.UUID{uuid.New()}

	i.houses.On("ExistsById", *request.HouseId).Return(true)
	i.groups.On("ExistsByIds", mock.Anything).Return(false)

	income, err := i.TestO.Add(request)

	assert.Equal(i.T(), int_errors.NewErrNotFound("not all group with ids %s found", common.Join(request.GroupIds, ",")), err)
	assert.Equal(i.T(), model.IncomeDto{}, income)

	i.incomeRepository.AssertNotCalled(i.T(), "Create", mock.Anything)
}

func (i *IncomeServiceTestSuite) Test_AddBatch() {
	request := mocks.GenerateCreateIncomeBatchRequest(2)
	repositoryResponse := common.MapSlice(request.Incomes, func(income model.CreateIncomeRequest) model.Income {
		return income.ToEntity()
	})

	i.houses.On("ExistsById", mock.Anything).Return(true)
	i.groups.On("ExistsByIds", mock.Anything).Return(true)
	i.incomeRepository.On("CreateBatch", mock.Anything).Return(repositoryResponse, nil)

	batch, err := i.TestO.AddBatch(request)

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), common.MapSlice(repositoryResponse, model.IncomeToDto), batch)
}

func (i *IncomeServiceTestSuite) Test_AddBatch_WithDefaultDetails() {
	request := mocks.GenerateCreateIncomeBatchRequest(2)
	request.Incomes[0].HouseId = nil
	request.Incomes[0].GroupIds = []uuid.UUID{uuid.New()}
	repositoryResponse := common.MapSlice(request.Incomes, func(income model.CreateIncomeRequest) model.Income {
		return income.ToEntity()
	})

	i.houses.On("ExistsById", mock.Anything).Return(true)
	i.groups.On("ExistsByIds", mock.Anything).Return(true)
	i.incomeRepository.On("CreateBatch", mock.Anything).Return(repositoryResponse, nil)

	batch, err := i.TestO.AddBatch(request)

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), common.MapSlice(repositoryResponse, model.IncomeToDto), batch)
}

func (i *IncomeServiceTestSuite) Test_AddBatch_WithMissingGroupIdsAndHouseId() {
	request := mocks.GenerateCreateIncomeBatchRequest(1)
	request.Incomes[0].HouseId = nil
	request.Incomes[0].GroupIds = []uuid.UUID{}

	i.houses.On("ExistsById", mock.Anything).Return(true)
	i.groups.On("ExistsByIds", mock.Anything).Return(true)

	_, err := i.TestO.AddBatch(request)

	assert.EqualError(i.T(), err, "houseId or groupId must be set")
	i.incomeRepository.AssertNotCalled(i.T(), "CreateBatch", mock.Anything)
}

func (i *IncomeServiceTestSuite) Test_AddBatch_WithEmptyData() {
	request := mocks.GenerateCreateIncomeBatchRequest(0)

	batch, err := i.TestO.AddBatch(request)

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), make([]model.IncomeDto, 0), batch)

	i.houses.AssertNotCalled(i.T(), "ExistsById", mock.Anything)
	i.groups.AssertNotCalled(i.T(), "ExistsByIds", mock.Anything)
	i.incomeRepository.AssertNotCalled(i.T(), "CreateBatch", mock.Anything)
}

func (i *IncomeServiceTestSuite) Test_AddBatch_WithInvalidData() {
	request := mocks.GenerateCreateIncomeBatchRequest(3)
	request.Incomes[0].Date = time.Now().Add(time.Hour)
	request.Incomes[2].GroupIds = []uuid.UUID{uuid.New()}

	i.houses.On("ExistsById", *request.Incomes[0].HouseId).Return(true)
	i.houses.On("ExistsById", *request.Incomes[1].HouseId).Return(false)
	i.houses.On("ExistsById", *request.Incomes[2].HouseId).Return(true)
	i.groups.On("ExistsByIds", request.Incomes[0].GroupIds).Return(true)
	i.groups.On("ExistsByIds", request.Incomes[1].GroupIds).Return(true)
	i.groups.On("ExistsByIds", request.Incomes[2].GroupIds).Return(false)

	actual, err := i.TestO.AddBatch(request)

	var expectedResult []model.IncomeDto

	assert.Equal(i.T(), expectedResult, actual)

	builder := int_errors.NewBuilder()
	builder.WithMessage("Create income batch failed")
	builder.WithDetail(fmt.Sprintf("house with id %s not found", request.Incomes[1].HouseId))
	builder.WithDetail(fmt.Sprintf("not all group with ids %s found", common.Join(request.Incomes[2].GroupIds, ",")))
	builder.WithDetail(fmt.Sprintf("date should not be after current date"))

	expectedError := int_errors.NewErrResponse(builder).(*int_errors.ErrResponse)
	actualError := err.(*int_errors.ErrResponse)

	assert.Equal(i.T(), *expectedError, *actualError)

	i.incomeRepository.AssertNotCalled(i.T(), "CreateBatch", mock.Anything)
}

func (i *IncomeServiceTestSuite) Test_FindById() {
	houseId := uuid.New()
	income := mocks.GenerateIncome(&houseId)

	i.incomeRepository.On("FindById", income.Id).Return(income, nil)

	actual, err := i.TestO.FindById(income.Id)

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), income.ToDto(), actual)
}

func (i *IncomeServiceTestSuite) Test_FindById_WithNotExistingId() {
	id := uuid.New()

	i.incomeRepository.On("FindById", id).Return(model.Income{}, gorm.ErrRecordNotFound)

	actual, err := i.TestO.FindById(id)

	assert.Equal(i.T(), int_errors.NewErrNotFound("income with id %s not found", id), err)
	assert.Equal(i.T(), model.IncomeDto{}, actual)
}

func (i *IncomeServiceTestSuite) Test_FindById_WithError() {
	id := uuid.New()
	expectedError := errors.New("test")

	i.incomeRepository.On("FindById", id).Return(model.Income{}, expectedError)

	actual, err := i.TestO.FindById(id)

	assert.Equal(i.T(), expectedError, err)
	assert.Equal(i.T(), model.IncomeDto{}, actual)
}

func (i *IncomeServiceTestSuite) Test_FindByHouseId() {
	income := []model.IncomeDto{mocks.GenerateIncomeDto()}

	i.incomeRepository.On("FindByHouseId", *income[0].HouseId).Return(income, nil)

	actual := i.TestO.FindByHouseId(*income[0].HouseId)

	assert.Equal(i.T(), income, actual)
}

func (i *IncomeServiceTestSuite) Test_FindByHouseId_WithNotExistingRecords() {
	var income []model.IncomeDto

	houseId := uuid.New()

	i.incomeRepository.On("FindByHouseId", houseId).Return(income, nil)

	actual := i.TestO.FindByHouseId(houseId)

	assert.Equal(i.T(), income, actual)
}

func (i *IncomeServiceTestSuite) Test_ExistsById() {
	id := uuid.New()

	i.incomeRepository.On("ExistsById", id).Return(true)

	assert.True(i.T(), i.TestO.ExistsById(id))
}

func (i *IncomeServiceTestSuite) Test_ExistsById_WithNotExists() {
	id := uuid.New()

	i.incomeRepository.On("ExistsById", id).Return(false)

	assert.False(i.T(), i.TestO.ExistsById(id))
}

func (i *IncomeServiceTestSuite) Test_DeleteById() {
	id := uuid.New()

	i.incomeRepository.On("ExistsById", id).Return(true)
	i.incomeRepository.On("DeleteById", id).Return(nil)

	assert.Nil(i.T(), i.TestO.DeleteById(id))
}

func (i *IncomeServiceTestSuite) Test_DeleteById_WithNotExists() {
	id := uuid.New()

	i.incomeRepository.On("ExistsById", id).Return(false)

	assert.Equal(i.T(), int_errors.NewErrNotFound("income with id %s not found", id), i.TestO.DeleteById(id))

	i.incomeRepository.AssertNotCalled(i.T(), "DeleteById", id)
}

func (i *IncomeServiceTestSuite) Test_Update() {
	id, request := mocks.GenerateUpdateIncomeRequest()

	i.incomeRepository.On("ExistsById", id).Return(true)
	i.incomeRepository.On("Update", id, request).Return(nil)

	assert.Nil(i.T(), i.TestO.Update(id, request))

	i.incomeRepository.AssertCalled(i.T(), "Update", id, request)
}

func (i *IncomeServiceTestSuite) Test_Update_WithErrorFromDatabase() {
	id, request := mocks.GenerateUpdateIncomeRequest()

	i.incomeRepository.On("ExistsById", id).Return(true)
	i.incomeRepository.On("Update", id, request).Return(errors.New("test"))

	err := i.TestO.Update(id, request)
	assert.Equal(i.T(), errors.New("test"), err)
}

func (i *IncomeServiceTestSuite) Test_Update_WithNotExists() {
	id, request := mocks.GenerateUpdateIncomeRequest()

	i.incomeRepository.On("ExistsById", id).Return(false)

	err := i.TestO.Update(id, request)
	assert.Equal(i.T(), int_errors.NewErrNotFound("income with id %s not found", id), err)

	i.incomeRepository.AssertNotCalled(i.T(), "Update", id, request)
}

func (i *IncomeServiceTestSuite) Test_Update_WithDateAfterCurrentDate() {
	id, request := mocks.GenerateUpdateIncomeRequest()
	request.Date = time.Now().Add(time.Hour)

	i.incomeRepository.On("ExistsById", id).Return(true)

	err := i.TestO.Update(id, request)
	assert.Equal(i.T(), errors.New("date should not be after current date"), err)

	i.incomeRepository.AssertNotCalled(i.T(), "Update", id, request)
}

func (i *IncomeServiceTestSuite) Test_Update_WithGroupsIdsNotFound() {
	id, request := mocks.GenerateUpdateIncomeRequest()
	request.GroupIds = []uuid.UUID{uuid.New()}

	i.incomeRepository.On("ExistsById", id).Return(true)
	i.groups.On("ExistsByIds", mock.Anything).Return(false)

	err := i.TestO.Update(id, request)

	assert.Equal(i.T(), int_errors.NewErrNotFound("not all group with ids %s found", common.Join(request.GroupIds, ",")), err)

	i.incomeRepository.AssertNotCalled(i.T(), "Update", id, request)
}
