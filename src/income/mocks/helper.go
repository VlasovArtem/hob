package mocks

import (
	"github.com/google/uuid"
	"income/model"
	"time"
)

var Date = time.Date(2021, time.January, 1, 0, 0, 0, 0, time.Local)

func GenerateIncome(houseId uuid.UUID) model.Income {
	return model.Income{
		Id:          uuid.New(),
		Name:        "Name",
		Date:        Date,
		Description: "Description",
		Sum:         100.1,
		HouseId:     houseId,
	}
}

func GenerateCreateIncomeRequest() model.CreateIncomeRequest {
	return model.CreateIncomeRequest{
		Name:        "Name",
		Date:        Date,
		Description: "Description",
		Sum:         100.1,
		HouseId:     uuid.New(),
	}
}

func GenerateIncomeResponse() model.IncomeResponse {
	return model.IncomeResponse{
		Id:          uuid.New(),
		Name:        "Name",
		Date:        Date,
		Description: "Description",
		Sum:         100.1,
		HouseId:     uuid.New(),
	}
}
