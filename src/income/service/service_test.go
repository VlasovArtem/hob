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
	"time"
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
	assert.Equal(t, savedIncome.ToDto(), income)
}

func Test_AddIncome_WithHouseNotExists(t *testing.T) {
	service := serviceGenerator()

	request := mocks.GenerateCreateIncomeRequest()

	houses.On("ExistsById", request.HouseId).Return(false)

	payment, err := service.Add(request)

	assert.Equal(t, errors.New(fmt.Sprintf("house with id %s not exists", request.HouseId)), err)
	assert.Equal(t, model.IncomeDto{}, payment)

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
	assert.Equal(t, model.IncomeDto{}, income)
}

func Test_AddIncome_WithDateAfterCurrentDate(t *testing.T) {
	service := serviceGenerator()

	request := mocks.GenerateCreateIncomeRequest()
	request.Date = time.Now().Add(time.Hour)

	houses.On("ExistsById", request.HouseId).Return(true)

	payment, err := service.Add(request)

	assert.Equal(t, errors.New("date should not be after current date"), err)
	assert.Equal(t, model.IncomeDto{}, payment)

	incomeRepository.AssertNotCalled(t, "Create", mock.Anything)
}

func Test_FindById(t *testing.T) {
	service := serviceGenerator()

	income := mocks.GenerateIncomeResponse()

	incomeRepository.On("FindDtoById", income.Id).Return(income, nil)

	actual, err := service.FindById(income.Id)

	assert.Nil(t, err)
	assert.Equal(t, income, actual)
}

func Test_FindById_WithNotExistingId(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()

	incomeRepository.On("FindDtoById", id).Return(model.IncomeDto{}, gorm.ErrRecordNotFound)

	actual, err := service.FindById(id)

	assert.Equal(t, errors.New(fmt.Sprintf("income with id %s not exists", id)), err)
	assert.Equal(t, model.IncomeDto{}, actual)
}

func Test_FindById_WithError(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()
	expectedError := errors.New("test")

	incomeRepository.On("FindDtoById", id).Return(model.IncomeDto{}, expectedError)

	actual, err := service.FindById(id)

	assert.Equal(t, expectedError, err)
	assert.Equal(t, model.IncomeDto{}, actual)
}

func Test_FindByHouseId(t *testing.T) {
	service := serviceGenerator()

	income := []model.IncomeDto{mocks.GenerateIncomeResponse()}

	incomeRepository.On("FindResponseByHouseId", income[0].HouseId).Return(income, nil)

	actual := service.FindByHouseId(income[0].HouseId)

	assert.Equal(t, income, actual)
}

func Test_FindByHouseId_WithNotExistingRecords(t *testing.T) {
	service := serviceGenerator()

	var income []model.IncomeDto

	houseId := uuid.New()

	incomeRepository.On("FindResponseByHouseId", houseId).Return(income, nil)

	actual := service.FindByHouseId(houseId)

	assert.Equal(t, income, actual)
}

func Test_ExistsById(t *testing.T) {
	paymentService := serviceGenerator()

	id := uuid.New()

	incomeRepository.On("ExistsById", id).Return(true)

	assert.True(t, paymentService.ExistsById(id))
}

func Test_ExistsById_WithNotExists(t *testing.T) {
	paymentService := serviceGenerator()

	id := uuid.New()

	incomeRepository.On("ExistsById", id).Return(false)

	assert.False(t, paymentService.ExistsById(id))
}

func Test_DeleteById(t *testing.T) {
	paymentService := serviceGenerator()

	id := uuid.New()

	incomeRepository.On("ExistsById", id).Return(true)
	incomeRepository.On("DeleteById", id).Return(nil)

	assert.Nil(t, paymentService.DeleteById(id))
}

func Test_DeleteById_WithNotExists(t *testing.T) {
	paymentService := serviceGenerator()

	id := uuid.New()

	incomeRepository.On("ExistsById", id).Return(false)

	assert.Equal(t, errors.New(fmt.Sprintf("income with id %s not found", id)), paymentService.DeleteById(id))

	incomeRepository.AssertNotCalled(t, "DeleteById", id)
}

func Test_Update(t *testing.T) {
	houseService := serviceGenerator()

	request := mocks.GenerateUpdateIncomeRequest()

	incomeRepository.On("ExistsById", request.Id).Return(true)
	incomeRepository.On("Update", mock.Anything).Return(nil)

	assert.Nil(t, houseService.Update(request))

	incomeRepository.AssertCalled(t, "Update", model.Income{
		Id:          request.Id,
		Name:        request.Name,
		Description: request.Description,
		Date:        request.Date,
		Sum:         request.Sum,
	})
}

func Test_Update_WithErrorFromDatabase(t *testing.T) {
	houseService := serviceGenerator()

	request := mocks.GenerateUpdateIncomeRequest()

	incomeRepository.On("ExistsById", request.Id).Return(true)
	incomeRepository.On("Update", mock.Anything).Return(errors.New("test"))

	err := houseService.Update(request)
	assert.Equal(t, errors.New("test"), err)
}

func Test_Update_WithNotExists(t *testing.T) {
	houseService := serviceGenerator()

	request := mocks.GenerateUpdateIncomeRequest()

	incomeRepository.On("ExistsById", request.Id).Return(false)

	err := houseService.Update(request)
	assert.Equal(t, errors.New(fmt.Sprintf("income with id %s not found", request.Id)), err)

	incomeRepository.AssertNotCalled(t, "Update", mock.Anything)
}

func Test_Update_WithDateAfterCurrentDate(t *testing.T) {
	houseService := serviceGenerator()

	request := mocks.GenerateUpdateIncomeRequest()
	request.Date = time.Now().Add(time.Hour)

	incomeRepository.On("ExistsById", request.Id).Return(true)

	err := houseService.Update(request)
	assert.Equal(t, errors.New("date should not be after current date"), err)

	incomeRepository.AssertNotCalled(t, "Update", mock.Anything)
}
