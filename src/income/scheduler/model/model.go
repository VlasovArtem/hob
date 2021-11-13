package model

import (
	"github.com/google/uuid"
	"income/model"
	"scheduler"
)

type IncomeScheduler struct {
	model.Income
	Spec scheduler.SchedulingSpecification
}

type CreateIncomeSchedulerRequest struct {
	Name        string
	Description string
	Sum         float32
	HouseId     uuid.UUID
	Spec        scheduler.SchedulingSpecification
}

type IncomeSchedulerResponse struct {
	IncomeScheduler
}

func (i IncomeScheduler) ToResponse() IncomeSchedulerResponse {
	return IncomeSchedulerResponse{i}
}

func (c CreateIncomeSchedulerRequest) ToEntity() IncomeScheduler {
	return IncomeScheduler{
		Income: model.Income{
			Id:          uuid.New(),
			Name:        c.Name,
			Description: c.Description,
			Sum:         c.Sum,
			HouseId:     c.HouseId,
		},
		Spec: c.Spec,
	}
}
