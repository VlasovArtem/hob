package model

import (
	"github.com/google/uuid"
	"payment/model"
	"scheduler"
)

type PaymentScheduler struct {
	model.Payment
	Spec scheduler.SchedulingSpecification
}

type CreatePaymentSchedulerRequest struct {
	Name        string
	Description string
	HouseId     uuid.UUID
	UserId      uuid.UUID
	Sum         float32
	Spec        scheduler.SchedulingSpecification
}

type PaymentSchedulerResponse struct {
	PaymentScheduler
}

func (ps PaymentScheduler) ToResponse() PaymentSchedulerResponse {
	return PaymentSchedulerResponse{ps}
}

func (request CreatePaymentSchedulerRequest) ToEntity() PaymentScheduler {
	return PaymentScheduler{
		Payment: model.Payment{
			Id:          uuid.New(),
			Name:        request.Name,
			Description: request.Description,
			HouseId:     request.HouseId,
			UserId:      request.UserId,
			Sum:         request.Sum,
		},
		Spec: request.Spec,
	}
}
