package model

import (
	"github.com/google/uuid"
	houseModel "house/model"
	"scheduler"
	userModel "user/model"
)

type PaymentScheduler struct {
	Id          uuid.UUID `gorm:"primarykey"`
	Name        string
	Description string
	HouseId     uuid.UUID
	UserId      uuid.UUID
	Sum         float32
	User        userModel.User   `gorm:"foreignKey:UserId"`
	House       houseModel.House `gorm:"foreignKey:HouseId"`
	Spec        scheduler.SchedulingSpecification
}

type CreatePaymentSchedulerRequest struct {
	Name        string
	Description string
	HouseId     uuid.UUID
	UserId      uuid.UUID
	Sum         float32
	Spec        scheduler.SchedulingSpecification
}

type PaymentSchedulerDto struct {
	Id          uuid.UUID
	Name        string
	Description string
	HouseId     uuid.UUID
	UserId      uuid.UUID
	Sum         float32
	Spec        scheduler.SchedulingSpecification
}

func (ps PaymentScheduler) ToDto() PaymentSchedulerDto {
	return PaymentSchedulerDto{
		Id:          ps.Id,
		Name:        ps.Name,
		Description: ps.Description,
		HouseId:     ps.HouseId,
		UserId:      ps.UserId,
		Sum:         ps.Sum,
		Spec:        ps.Spec,
	}
}

func (request CreatePaymentSchedulerRequest) ToEntity() PaymentScheduler {
	return PaymentScheduler{
		Id:          uuid.New(),
		Name:        request.Name,
		Description: request.Description,
		HouseId:     request.HouseId,
		UserId:      request.UserId,
		Sum:         request.Sum,
		Spec:        request.Spec,
	}
}
