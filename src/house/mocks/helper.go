package mocks

import (
	"github.com/VlasovArtem/hob/src/house/model"
	"github.com/google/uuid"
	"strconv"
)

func GenerateHouse(userId uuid.UUID) model.House {
	return model.House{
		Id:          uuid.New(),
		Name:        "Name",
		CountryCode: "UA",
		City:        "City",
		StreetLine1: "Street Line 1",
		StreetLine2: "Street Line 2",
		UserId:      userId,
	}
}

func GenerateCreateHouseRequest() model.CreateHouseRequest {
	return model.CreateHouseRequest{
		Name:        "Test House",
		CountryCode: "UA",
		City:        "City",
		StreetLine1: "StreetLine1",
		StreetLine2: "StreetLine2",
		UserId:      uuid.New(),
	}
}

func GenerateCreateHouseBatchRequest(number int) model.CreateHouseBatchRequest {
	request := model.CreateHouseBatchRequest{}
	for i := 0; i < number; i++ {
		houseRequest := GenerateCreateHouseRequest()
		houseRequest.Name = "House Name #" + strconv.Itoa(i)
		request.Houses = append(request.Houses, houseRequest)
	}
	return request
}

func GenerateHouseResponse() model.HouseDto {
	return model.HouseDto{
		Id:          uuid.New(),
		Name:        "Test Name",
		CountryCode: "Ukraine",
		City:        "City",
		StreetLine1: "StreetLine1",
		StreetLine2: "StreetLine2",
		UserId:      uuid.New(),
	}
}

func GenerateUpdateHouseRequest() (uuid.UUID, model.UpdateHouseRequest) {
	return uuid.New(), model.UpdateHouseRequest{
		Name:        "Test House NEW",
		CountryCode: "UA",
		City:        "City NEW",
		StreetLine1: "StreetLine1 NEW",
		StreetLine2: "StreetLine2 NEW",
	}
}
