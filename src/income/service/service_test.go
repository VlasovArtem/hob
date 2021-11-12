package service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"income/model"
	"test/mock"
	"test/testhelper"
	"testing"
	"time"
)

var (
	houses  *mock.HouseServiceMock
	houseId = testhelper.ParseUUID("5483a7b3-a7e3-4697-b39e-4dc93519ac38")
	date    = time.Now()
)

func serviceGenerator() IncomeService {
	houses = new(mock.HouseServiceMock)

	return NewIncomeService(houses)
}

func Test_AddIncome(t *testing.T) {
	service := serviceGenerator()

	houses.On("ExistsById", houseId).Return(true)

	request := generateCreateIncomeRequest()

	meter, err := service.AddIncome(request)

	expectedResponse := request.ToEntity().ToResponse()
	expectedResponse.Id = meter.Id

	assert.Nil(t, err)
	assert.Equal(t, expectedResponse, meter)
}

func Test_AddIncome_WithHouseNotExists(t *testing.T) {
	service := serviceGenerator()

	houses.On("ExistsById", houseId).Return(false)

	request := generateCreateIncomeRequest()

	payment, err := service.AddIncome(request)

	assert.Equal(t, errors.New(fmt.Sprintf("house with id %s not exists", request.HouseId)), err)
	assert.Equal(t, model.IncomeResponse{}, payment)
}

func Test_FindById(t *testing.T) {
	service := serviceGenerator()

	income := generateIncome()

	service.(*incomeServiceObject).incomes[income.Id] = income

	actual, err := service.FindById(income.Id)

	assert.Nil(t, err)
	assert.Equal(t, income.ToResponse(), actual)
}

func Test_FindById_WithNotExistingId(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()
	actual, err := service.FindById(id)

	assert.Equal(t, errors.New(fmt.Sprintf("income with id %s not exists", id)), err)
	assert.Equal(t, model.IncomeResponse{}, actual)
}

func Test_FindByHouseId(t *testing.T) {
	service := serviceGenerator()

	income := generateIncome()

	service.(*incomeServiceObject).houseIncomes[income.HouseId] = income

	actual, err := service.FindByHouseId(income.HouseId)

	assert.Nil(t, err)
	assert.Equal(t, income.ToResponse(), actual)
}

func Test_FindByHouseId_WithNotExistingRecords(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()
	actual, err := service.FindByHouseId(id)

	assert.Equal(t, errors.New(fmt.Sprintf("income with house id %s not found", id)), err)
	assert.Equal(t, model.IncomeResponse{}, actual)
}

func generateCreateIncomeRequest() model.CreateIncomeRequest {
	return model.CreateIncomeRequest{
		Name:        "Name",
		Date:        date,
		Description: "Description",
		Sum:         100.1,
		HouseId:     houseId,
	}
}

func generateIncome() model.Income {
	return model.Income{
		Id:          uuid.New(),
		Name:        "Name",
		Date:        date,
		Description: "Description",
		Sum:         100.1,
		HouseId:     houseId,
	}
}
