package mocks

import (
	"github.com/VlasovArtem/hob/src/income/model"
	"github.com/google/uuid"
	"strconv"
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

func GenerateCreateIncomeBatchRequest(number int) model.CreateIncomeBatchRequest {
	request := model.CreateIncomeBatchRequest{}
	for i := 0; i < number; i++ {
		incomeRequest := GenerateCreateIncomeRequest()
		incomeRequest.Name = "Income Name #" + strconv.Itoa(i)
		request.Incomes = append(request.Incomes, incomeRequest)
	}
	return request
}

func GenerateUpdateIncomeRequest() (uuid.UUID, model.UpdateIncomeRequest) {
	return uuid.New(), model.UpdateIncomeRequest{
		Name:        "Name",
		Date:        Date,
		Description: "Description",
		Sum:         100.1,
	}
}

func GenerateIncomeDto() model.IncomeDto {
	return model.IncomeDto{
		Id:          uuid.New(),
		Name:        "Name",
		Date:        Date,
		Description: "Description",
		Sum:         100.1,
		HouseId:     uuid.New(),
	}
}
