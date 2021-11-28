package mocks

import (
	"github.com/google/uuid"
	"payment/scheduler/model"
	"scheduler"
	"test/testhelper"
)

var (
	HouseId = testhelper.ParseUUID("d077adaa-00d7-4e80-ac86-57512267505d")
	UserId  = testhelper.ParseUUID("ad2c5035-6745-48d0-9eee-fd22f5dae8e0")
)

func GenerateCreatePaymentSchedulerRequest(houseId uuid.UUID, userId uuid.UUID) model.CreatePaymentSchedulerRequest {
	return model.CreatePaymentSchedulerRequest{
		Name:        "Test Payment",
		Description: "Test Payment Description",
		HouseId:     houseId,
		UserId:      userId,
		Sum:         1000,
		Spec:        scheduler.DAILY,
	}
}

func GeneratePaymentScheduler(houseId uuid.UUID, userId uuid.UUID) model.PaymentScheduler {
	return model.PaymentScheduler{
		Id:          uuid.New(),
		Name:        "Test Payment",
		Description: "Test Payment Description",
		HouseId:     houseId,
		UserId:      userId,
		Sum:         1000,
		Spec:        scheduler.DAILY,
	}
}

func GeneratePaymentSchedulerResponse(id uuid.UUID, houseId uuid.UUID, userId uuid.UUID) model.PaymentSchedulerDto {
	paymentScheduler := GeneratePaymentScheduler(houseId, userId)
	paymentScheduler.Id = id
	return model.PaymentSchedulerDto{
		PaymentScheduler: paymentScheduler,
	}
}
