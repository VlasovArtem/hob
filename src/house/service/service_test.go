package service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"house/mocks"
	"house/model"
	"test/testhelper"
	"testing"
	userMocks "user/mocks"
)

var (
	users            *userMocks.UserService
	repository       *mocks.HouseRepository
	countriesService = testhelper.InitCountryService()
)

func serviceGenerator() HouseService {
	users = new(userMocks.UserService)
	repository = new(mocks.HouseRepository)

	return NewHouseService(countriesService, users, repository)
}

func Test_Add(t *testing.T) {
	service := serviceGenerator()

	request := mocks.GenerateCreateHouseRequest()

	users.On("ExistsById", request.UserId).Return(true)
	repository.On("Create", mock.Anything).Return(
		func(house model.House) model.House { return house },
		nil,
	)

	actual, err := service.Add(request)

	assert.Nil(t, err)
	assert.Equal(t, model.HouseDto{
		Id:          actual.Id,
		Name:        "Test House",
		CountryCode: "UA",
		City:        "City",
		StreetLine1: "StreetLine1",
		StreetLine2: "StreetLine2",
		UserId:      request.UserId,
	}, actual)
}

func Test_FindById(t *testing.T) {
	houseService := serviceGenerator()

	house := mocks.GenerateHouseResponse()

	repository.On("FindResponseById", house.Id).Return(house, nil)

	actual, err := houseService.FindById(house.Id)

	assert.Nil(t, err)
	assert.Equal(t, house, actual)
}

func Test_FindById_WithRecordNotFound(t *testing.T) {
	houseService := serviceGenerator()

	id := uuid.New()

	repository.On("FindResponseById", id).Return(model.HouseDto{}, gorm.ErrRecordNotFound)

	actual, err := houseService.FindById(id)

	assert.Equal(t, errors.New(fmt.Sprintf("house with id %s not found", id)), err)
	assert.Equal(t, model.HouseDto{}, actual)
}

func Test_FindById_WithRecordNotFoundExists(t *testing.T) {
	houseService := serviceGenerator()

	id := uuid.New()

	expectedError := errors.New("error")
	repository.On("FindResponseById", id).Return(model.HouseDto{}, expectedError)

	actual, err := houseService.FindById(id)

	assert.Equal(t, expectedError, err)
	assert.Equal(t, model.HouseDto{}, actual)
}

func Test_FindByUserId(t *testing.T) {
	houseService := serviceGenerator()

	house := mocks.GenerateHouseResponse()

	repository.On("FindResponseByUserId", house.UserId).Return([]model.HouseDto{house})

	actual := houseService.FindByUserId(house.UserId)

	assert.Equal(t, []model.HouseDto{house}, actual)
}

func Test_FindByUserId_WithNotExists(t *testing.T) {
	houseService := serviceGenerator()

	id := uuid.New()

	repository.On("FindResponseByUserId", id).Return([]model.HouseDto{})

	actual := houseService.FindByUserId(id)

	assert.Equal(t, []model.HouseDto{}, actual)
}

func Test_ExistsById(t *testing.T) {
	houseService := serviceGenerator()

	houseId := uuid.New()

	repository.On("ExistsById", houseId).Return(true)

	assert.True(t, houseService.ExistsById(houseId))
}

func Test_ExistsById_WithNotExists(t *testing.T) {
	houseService := serviceGenerator()

	id := uuid.New()

	repository.On("ExistsById", id).Return(false)

	assert.False(t, houseService.ExistsById(id))
}
