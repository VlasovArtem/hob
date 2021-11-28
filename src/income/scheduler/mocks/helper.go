package mocks

import (
	im "github.com/VlasovArtem/hob/src/income/model"
	"github.com/VlasovArtem/hob/src/income/scheduler/model"
	"github.com/VlasovArtem/hob/src/scheduler"
	"github.com/google/uuid"
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
