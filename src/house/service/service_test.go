package service

import (
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/house/mocks"
	"github.com/VlasovArtem/hob/src/house/model"
	"github.com/VlasovArtem/hob/src/test/testhelper"
	userMocks "github.com/VlasovArtem/hob/src/user/mocks"
	userModel "github.com/VlasovArtem/hob/src/user/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"testing"
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

	repository.On("FindById", house.Id).Return(house, nil)

	actual, err := houseService.FindById(house.Id)

	assert.Nil(t, err)
	assert.Equal(t, house, actual)
}

func Test_FindById_WithRecordNotFound(t *testing.T) {
	houseService := serviceGenerator()

	id := uuid.New()

	repository.On("FindById", id).Return(model.HouseDto{}, gorm.ErrRecordNotFound)

	actual, err := houseService.FindById(id)

	assert.Equal(t, errors.New(fmt.Sprintf("house with id %s not found", id)), err)
	assert.Equal(t, model.HouseDto{}, actual)
}

func Test_FindById_WithRecordNotFoundExists(t *testing.T) {
	houseService := serviceGenerator()

	id := uuid.New()

	expectedError := errors.New("error")
	repository.On("FindById", id).Return(model.HouseDto{}, expectedError)

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

func Test_DeleteById(t *testing.T) {
	houseService := serviceGenerator()

	id := uuid.New()

	repository.On("ExistsById", id).Return(true)
	repository.On("DeleteById", id).Return(nil)

	assert.Nil(t, houseService.DeleteById(id))
}

func Test_DeleteById_WithNotExists(t *testing.T) {
	houseService := serviceGenerator()

	id := uuid.New()

	repository.On("ExistsById", id).Return(false)

	assert.Equal(t, errors.New(fmt.Sprintf("house with id %s not found", id)), houseService.DeleteById(id))

	repository.AssertNotCalled(t, "DeleteById", id)
}

func Test_Update(t *testing.T) {
	houseService := serviceGenerator()

	request := mocks.GenerateUpdateHouseRequest()

	repository.On("ExistsById", request.Id).Return(true)
	repository.On("Update", mock.Anything).Return(nil)

	assert.Nil(t, houseService.Update(request))

	repository.AssertCalled(t, "Update", model.House{
		Id:          request.Id,
		Name:        request.Name,
		CountryCode: request.Country,
		City:        request.City,
		StreetLine1: request.StreetLine1,
		StreetLine2: request.StreetLine2,
		UserId:      uuid.UUID{},
		User:        userModel.User{},
	})
}

func Test_Update_WithErrorFromDatabase(t *testing.T) {
	houseService := serviceGenerator()

	request := mocks.GenerateUpdateHouseRequest()

	repository.On("ExistsById", request.Id).Return(true)
	repository.On("Update", mock.Anything).Return(errors.New("test"))

	err := houseService.Update(request)
	assert.Equal(t, errors.New("test"), err)
}

func Test_Update_WithNotExists(t *testing.T) {
	houseService := serviceGenerator()

	request := mocks.GenerateUpdateHouseRequest()

	repository.On("ExistsById", request.Id).Return(false)

	err := houseService.Update(request)
	assert.Equal(t, errors.New(fmt.Sprintf("house with id %s not found", request.Id)), err)

	repository.AssertNotCalled(t, "Update", mock.Anything)
}

func Test_Update_WithNotMatchingCountry(t *testing.T) {
	houseService := serviceGenerator()

	request := mocks.GenerateUpdateHouseRequest()
	request.Country = "invalid"

	repository.On("ExistsById", request.Id).Return(true)

	err := houseService.Update(request)
	assert.Equal(t, errors.New(fmt.Sprintf("country with code %s is not found", request.Country)), err)

	repository.AssertNotCalled(t, "Update", mock.Anything)
}
