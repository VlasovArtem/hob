package service

import (
	"errors"
	"fmt"
	houseMocks "github.com/VlasovArtem/hob/src/house/mocks"
	"github.com/VlasovArtem/hob/src/income/mocks"
	"github.com/VlasovArtem/hob/src/income/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"testing"
)

var (
	houses           *houseMocks.HouseService
	incomeRepository *mocks.IncomeRepository
)

func serviceGenerator() IncomeService {
	houses = new(houseMocks.HouseService)
	incomeRepository = new(mocks.IncomeRepository)

	return NewIncomeService(houses, incomeRepository)
}

func Test_AddIncome(t *testing.T) {
	service := serviceGenerator()

	var savedIncome model.Income
	request := mocks.GenerateCreateIncomeRequest()

	houses.On("ExistsById", request.HouseId).Return(true)
	incomeRepository.On("Create", mock.Anything).Return(func(income model.Income) model.Income {
		savedIncome = income

		return income
	}, nil)

	income, err := service.Add(request)

	assert.Nil(t, err)
	assert.Equal(t, savedIncome.ToResponse(), income)
}

func Test_AddIncome_WithHouseNotExists(t *testing.T) {
	service := serviceGenerator()

	request := mocks.GenerateCreateIncomeRequest()

	houses.On("ExistsById", request.HouseId).Return(false)

	payment, err := service.Add(request)

	assert.Equal(t, errors.New(fmt.Sprintf("house with id %s not exists", request.HouseId)), err)
	assert.Equal(t, model.IncomeResponse{}, payment)

	incomeRepository.AssertNotCalled(t, "Create", mock.Anything)
}

func Test_AddIncome_WithErrorFromRepository(t *testing.T) {
	service := serviceGenerator()

	expectedError := errors.New("error")
	request := mocks.GenerateCreateIncomeRequest()

	houses.On("ExistsById", request.HouseId).Return(true)
	incomeRepository.On("Create", mock.Anything).Return(model.Income{}, expectedError)

	income, err := service.Add(request)

	assert.Equal(t, expectedError, err)
	assert.Equal(t, model.IncomeResponse{}, income)
}

func Test_FindById(t *testing.T) {
	service := serviceGenerator()

	income := mocks.GenerateIncomeResponse()

	incomeRepository.On("FindResponseById", income.Id).Return(income, nil)

	actual, err := service.FindById(income.Id)

	assert.Nil(t, err)
	assert.Equal(t, income, actual)
}

func Test_FindById_WithNotExistingId(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()

	incomeRepository.On("FindResponseById", id).Return(model.IncomeResponse{}, gorm.ErrRecordNotFound)

	actual, err := service.FindById(id)

	assert.Equal(t, errors.New(fmt.Sprintf("income with id %s not exists", id)), err)
	assert.Equal(t, model.IncomeResponse{}, actual)
}

func Test_FindById_WithError(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()
	expectedError := errors.New("test")

	incomeRepository.On("FindResponseById", id).Return(model.IncomeResponse{}, expectedError)

	actual, err := service.FindById(id)

	assert.Equal(t, expectedError, err)
	assert.Equal(t, model.IncomeResponse{}, actual)
}

func Test_FindByHouseId(t *testing.T) {
	service := serviceGenerator()

	income := []model.IncomeResponse{mocks.GenerateIncomeResponse()}

	incomeRepository.On("FindResponseByHouseId", income[0].HouseId).Return(income, nil)

	actual := service.FindByHouseId(income[0].HouseId)

	assert.Equal(t, income, actual)
}

func Test_FindByHouseId_WithNotExistingRecords(t *testing.T) {
	service := serviceGenerator()

	var income []model.IncomeResponse

	houseId := uuid.New()

	incomeRepository.On("FindResponseByHouseId", houseId).Return(income, nil)

	actual := service.FindByHouseId(houseId)

	assert.Equal(t, income, actual)
}
