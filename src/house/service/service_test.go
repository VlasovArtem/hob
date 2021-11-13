package service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"house/model"
	"test"
	"test/mock"
	"test/testhelper"
	"testing"
)

var (
	users            *mock.UserServiceMock
	countriesService = testhelper.InitCountryService()
	userId           = testhelper.ParseUUID("0757088a-ed8e-465e-b9ed-34ebacbfb3be")
)

func serviceGenerator() HouseService {
	users = new(mock.UserServiceMock)

	return NewHouseService(countriesService, users)
}

func Test_Add(t *testing.T) {
	service := serviceGenerator()

	users.On("ExistsById", userId).Return(true)

	request := generateCreateHouseRequest()

	actual, err := service.Add(request)

	assert.Nil(t, err)
	assert.Equal(t, model.HouseResponse{
		Id:          actual.Id,
		Name:        "Test House",
		Country:     "Ukraine",
		City:        "City",
		StreetLine1: "StreetLine1",
		StreetLine2: "StreetLine2",
		UserId:      userId,
	}, actual)

	serviceObject := service.(*houseServiceObject)

	_, houseExists := serviceObject.houses[actual.Id]
	assert.True(t, houseExists)

	_, userHouseExists := serviceObject.userHouses[actual.UserId]
	assert.True(t, userHouseExists)
}

func Test_FindById(t *testing.T) {
	houseService := serviceGenerator()

	house := generateCreateHouse()

	houseService.(*houseServiceObject).houses[house.Id] = house

	actual, err := houseService.FindById(house.Id)

	assert.Nil(t, err)
	assert.Equal(t, house.ToResponse(), actual)
}

func Test_FindById_WithNotExists(t *testing.T) {
	houseService := serviceGenerator()

	id := uuid.New()

	actual, err := houseService.FindById(id)

	assert.Equal(t, errors.New(fmt.Sprintf("house with id %s not found", id)), err)
	assert.Equal(t, model.HouseResponse{}, actual)
}

func Test_FindByUserId(t *testing.T) {
	houseService := serviceGenerator()

	house := generateCreateHouse()

	houseService.(*houseServiceObject).userHouses[house.UserId] = []model.House{house}

	actual := houseService.FindByUserId(house.UserId)

	assert.Equal(t, []model.HouseResponse{house.ToResponse()}, actual)
}

func Test_FindByUserId_WithNotExists(t *testing.T) {
	houseService := serviceGenerator()

	actual := houseService.FindByUserId(uuid.New())

	assert.Equal(t, []model.HouseResponse{}, actual)
}

func Test_ExistsById(t *testing.T) {
	houseService := serviceGenerator()

	house := generateCreateHouse()

	houseService.(*houseServiceObject).houses[house.Id] = house

	assert.True(t, houseService.ExistsById(house.Id))
}

func Test_ExistsById_WithNotExists(t *testing.T) {
	houseService := serviceGenerator()

	assert.False(t, houseService.ExistsById(uuid.New()))
}

func generateCreateHouseRequest() model.CreateHouseRequest {
	return model.CreateHouseRequest{
		Name:        "Test House",
		Country:     "UA",
		City:        "City",
		StreetLine1: "StreetLine1",
		StreetLine2: "StreetLine2",
		UserId:      userId,
	}
}

func generateCreateHouse() model.House {
	return model.House{
		Id:          uuid.New(),
		Name:        "Test House",
		Country:     test.CountryObject,
		City:        "City",
		StreetLine1: "StreetLine1",
		StreetLine2: "StreetLine2",
		UserId:      userId,
	}
}
