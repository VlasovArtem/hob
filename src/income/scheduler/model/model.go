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

type IncomeSchedulerDto struct {
	Id          uuid.UUID
	Name        string
	Description string
	Sum         float32
	HouseId     uuid.UUID
	Spec        scheduler.SchedulingSpecification
}

func (i IncomeScheduler) ToDto() IncomeSchedulerDto {
	return IncomeSchedulerDto{
		Id:          i.Id,
		Name:        i.Name,
		Description: i.Description,
		Sum:         i.Sum,
		HouseId:     i.HouseId,
		Spec:        i.Spec,
	}
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
