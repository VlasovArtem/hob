package mocks

import (
	"github.com/VlasovArtem/hob/src/payment/scheduler/model"
	"github.com/VlasovArtem/hob/src/scheduler"
	"github.com/VlasovArtem/hob/src/test/testhelper"
	"github.com/google/uuid"
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
	return model.PaymentSchedulerDto{
		Id:          id,
		Name:        "Test Payment",
		Description: "Test Payment Description",
		HouseId:     houseId,
		UserId:      userId,
		Sum:         1000,
		Spec:        scheduler.DAILY,
	}
}
