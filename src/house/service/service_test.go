package service

import (
	"errors"
	"github.com/VlasovArtem/hob/src/common/int-errors"
	"github.com/VlasovArtem/hob/src/house/mocks"
	"github.com/VlasovArtem/hob/src/house/model"
	"github.com/VlasovArtem/hob/src/test/testhelper"
	userMocks "github.com/VlasovArtem/hob/src/user/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"testing"
)

var (
	users            *userMocks.UserService
	houseRepository  *mocks.HouseRepository
	countriesService = testhelper.InitCountryService()
)

func serviceGenerator() HouseService {
	users = new(userMocks.UserService)
	houseRepository = new(mocks.HouseRepository)

	return NewHouseService(countriesService, users, houseRepository)
}

func Test_Add(t *testing.T) {
	service := serviceGenerator()

	request := mocks.GenerateCreateHouseRequest()

	users.On("ExistsById", request.UserId).Return(true)
	houseRepository.On("Create", mock.Anything).Return(
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

	house := mocks.GenerateHouse(uuid.New())

	houseRepository.On("FindById", house.Id).Return(house, nil)

	actual, err := houseService.FindById(house.Id)

	assert.Nil(t, err)
	assert.Equal(t, house.ToDto(), actual)
}

func Test_FindById_WithRecordNotFound(t *testing.T) {
	houseService := serviceGenerator()

	id := uuid.New()

	houseRepository.On("FindById", id).Return(model.House{}, gorm.ErrRecordNotFound)

	actual, err := houseService.FindById(id)

	assert.Equal(t, int_errors.NewErrNotFound("house with id %s not found", id), err)
	assert.Equal(t, model.HouseDto{}, actual)
}

func Test_FindById_WithRecordNotFoundExists(t *testing.T) {
	houseService := serviceGenerator()

	id := uuid.New()

	expectedError := errors.New("error")
	houseRepository.On("FindById", id).Return(model.House{}, expectedError)

	actual, err := houseService.FindById(id)

	assert.Equal(t, expectedError, err)
	assert.Equal(t, model.HouseDto{}, actual)
}

func Test_FindByUserId(t *testing.T) {
	houseService := serviceGenerator()

	house := mocks.GenerateHouse(uuid.New())

	houseRepository.On("FindByUserId", house.UserId).Return([]model.House{house})

	actual := houseService.FindByUserId(house.UserId)

	assert.Equal(t, []model.HouseDto{house.ToDto()}, actual)
}

func Test_FindByUserId_WithNotExists(t *testing.T) {
	houseService := serviceGenerator()

	id := uuid.New()

	houseRepository.On("FindByUserId", id).Return([]model.House{})

	actual := houseService.FindByUserId(id)

	assert.Equal(t, []model.HouseDto{}, actual)
}

func Test_ExistsById(t *testing.T) {
	houseService := serviceGenerator()

	houseId := uuid.New()

	houseRepository.On("ExistsById", houseId).Return(true)

	assert.True(t, houseService.ExistsById(houseId))
}

func Test_ExistsById_WithNotExists(t *testing.T) {
	houseService := serviceGenerator()

	id := uuid.New()

	houseRepository.On("ExistsById", id).Return(false)

	assert.False(t, houseService.ExistsById(id))
}

func Test_DeleteById(t *testing.T) {
	houseService := serviceGenerator()

	id := uuid.New()

	houseRepository.On("ExistsById", id).Return(true)
	houseRepository.On("DeleteById", id).Return(nil)

	assert.Nil(t, houseService.DeleteById(id))
}

func Test_DeleteById_WithNotExists(t *testing.T) {
	houseService := serviceGenerator()

	id := uuid.New()

	houseRepository.On("ExistsById", id).Return(false)

	assert.Equal(t, int_errors.NewErrNotFound("house with id %s not found", id), houseService.DeleteById(id))

	houseRepository.AssertNotCalled(t, "DeleteById", id)
}

func Test_Update(t *testing.T) {
	houseService := serviceGenerator()

	id, request := mocks.GenerateUpdateHouseRequest()

	houseRepository.On("ExistsById", id).Return(true)
	houseRepository.On("Update", id, request).Return(nil)

	assert.Nil(t, houseService.Update(id, request))
}

func Test_Update_WithErrorFromDatabase(t *testing.T) {
	houseService := serviceGenerator()

	id, request := mocks.GenerateUpdateHouseRequest()

	houseRepository.On("ExistsById", id).Return(true)
	houseRepository.On("Update", id, request).Return(errors.New("test"))

	err := houseService.Update(id, request)
	assert.Equal(t, errors.New("test"), err)
}

func Test_Update_WithNotExists(t *testing.T) {
	houseService := serviceGenerator()

	id, request := mocks.GenerateUpdateHouseRequest()

	houseRepository.On("ExistsById", id).Return(false)

	err := houseService.Update(id, request)
	assert.Equal(t, int_errors.NewErrNotFound("house with id %s not found", id), err)

	houseRepository.AssertNotCalled(t, "Update", id, request)
}

func Test_Update_WithNotMatchingCountry(t *testing.T) {
	houseService := serviceGenerator()

	id, request := mocks.GenerateUpdateHouseRequest()
	request.CountryCode = "invalid"

	houseRepository.On("ExistsById", id).Return(true)

	err := houseService.Update(id, request)
	assert.Equal(t, int_errors.NewErrNotFound("country with code %s is not found", request.CountryCode), err)

	houseRepository.AssertNotCalled(t, "Update", id, request)
}
