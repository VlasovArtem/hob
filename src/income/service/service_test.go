package service

import (
	"errors"
	"github.com/VlasovArtem/hob/src/common/int-errors"
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

func Test_Add_WithHouseNotExists(t *testing.T) {
	service := serviceGenerator()

	request := mocks.GenerateCreateIncomeRequest()

	houses.On("ExistsById", request.HouseId).Return(false)

	payment, err := service.Add(request)

	assert.Equal(t, int_errors.NewErrNotFound("house with id %s not exists", request.HouseId), err)
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

	income := mocks.GenerateIncome(uuid.New())

	incomeRepository.On("FindById", income.Id).Return(income, nil)

	actual, err := service.FindById(income.Id)

	assert.Nil(t, err)
	assert.Equal(t, income.ToDto(), actual)
}

func Test_FindById_WithNotExistingId(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()

	incomeRepository.On("FindById", id).Return(model.Income{}, gorm.ErrRecordNotFound)

	actual, err := service.FindById(id)

	assert.Equal(t, int_errors.NewErrNotFound("income with id %s not found", id), err)
	assert.Equal(t, model.IncomeDto{}, actual)
}

func Test_FindById_WithError(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()
	expectedError := errors.New("test")

	incomeRepository.On("FindById", id).Return(model.Income{}, expectedError)

	actual, err := service.FindById(id)

	assert.Equal(t, expectedError, err)
	assert.Equal(t, model.IncomeDto{}, actual)
}

func Test_FindByHouseId(t *testing.T) {
	service := serviceGenerator()

	income := []model.IncomeDto{mocks.GenerateIncomeDto()}

	incomeRepository.On("FindByHouseId", income[0].HouseId).Return(income, nil)

	actual := service.FindByHouseId(income[0].HouseId)

	assert.Equal(t, income, actual)
}

func Test_FindByHouseId_WithNotExistingRecords(t *testing.T) {
	service := serviceGenerator()

	var income []model.IncomeDto

	houseId := uuid.New()

	incomeRepository.On("FindByHouseId", houseId).Return(income, nil)

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

	assert.Equal(t, int_errors.NewErrNotFound("income with id %s not found", id), paymentService.DeleteById(id))

	incomeRepository.AssertNotCalled(t, "DeleteById", id)
}

func Test_Update(t *testing.T) {
	houseService := serviceGenerator()

	id, request := mocks.GenerateUpdateIncomeRequest()

	incomeRepository.On("ExistsById", id).Return(true)
	incomeRepository.On("Update", id, request).Return(nil)

	assert.Nil(t, houseService.Update(id, request))

	incomeRepository.AssertCalled(t, "Update", id, request)
}

func Test_Update_WithErrorFromDatabase(t *testing.T) {
	houseService := serviceGenerator()

	id, request := mocks.GenerateUpdateIncomeRequest()

	incomeRepository.On("ExistsById", id).Return(true)
	incomeRepository.On("Update", id, request).Return(errors.New("test"))

	err := houseService.Update(id, request)
	assert.Equal(t, errors.New("test"), err)
}

func Test_Update_WithNotExists(t *testing.T) {
	houseService := serviceGenerator()

	id, request := mocks.GenerateUpdateIncomeRequest()

	incomeRepository.On("ExistsById", id).Return(false)

	err := houseService.Update(id, request)
	assert.Equal(t, int_errors.NewErrNotFound("income with id %s not found", id), err)

	incomeRepository.AssertNotCalled(t, "Update", id, request)
}

func Test_Update_WithDateAfterCurrentDate(t *testing.T) {
	houseService := serviceGenerator()

	id, request := mocks.GenerateUpdateIncomeRequest()
	request.Date = time.Now().Add(time.Hour)

	incomeRepository.On("ExistsById", id).Return(true)

	err := houseService.Update(id, request)
	assert.Equal(t, errors.New("date should not be after current date"), err)

	incomeRepository.AssertNotCalled(t, "Update", id, request)
}
