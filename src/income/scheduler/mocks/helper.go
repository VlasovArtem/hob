package mocks

import (
	"github.com/google/uuid"
	im "income/model"
	"income/scheduler/model"
	"scheduler"
	"time"
)

var Date = time.Date(2021, time.January, 1, 0, 0, 0, 0, time.Local)

func GenerateIncomeScheduler(houseId uuid.UUID) model.IncomeScheduler {
	return model.IncomeScheduler{
		Income: im.Income{
			Id:          uuid.New(),
			Name:        "Name",
			Description: "Description",
			Date:        Date,
			Sum:         1000,
			HouseId:     houseId,
		},
		Spec: scheduler.DAILY,
	}
}

func GenerateCreateIncomeSchedulerRequest() model.CreateIncomeSchedulerRequest {
	return model.CreateIncomeSchedulerRequest{
		Name:        "Test Income",
		Description: "Test Income Description",
		HouseId:     uuid.New(),
		Sum:         1000,
		Spec:        scheduler.DAILY,
	}
}
