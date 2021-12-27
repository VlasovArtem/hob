package mocks

import (
	"github.com/VlasovArtem/hob/src/payment/scheduler/model"
	"github.com/VlasovArtem/hob/src/scheduler"
	"github.com/VlasovArtem/hob/src/test/testhelper"
	"github.com/google/uuid"
)

var (
	HouseId    = testhelper.ParseUUID("d077adaa-00d7-4e80-ac86-57512267505d")
	UserId     = testhelper.ParseUUID("ad2c5035-6745-48d0-9eee-fd22f5dae8e0")
	ProviderId = testhelper.ParseUUID("50dfd5b5-3d5e-4149-bb3f-003f9296d418")
)

func GenerateCreatePaymentSchedulerRequest() model.CreatePaymentSchedulerRequest {
	return model.CreatePaymentSchedulerRequest{
		Name:        "Test Payment",
		Description: "Test Payment Description",
		HouseId:     HouseId,
		UserId:      UserId,
		ProviderId:  ProviderId,
		Sum:         1000,
		Spec:        scheduler.DAILY,
	}
}

func GenerateUpdatePaymentSchedulerRequest() model.UpdatePaymentSchedulerRequest {
	return model.UpdatePaymentSchedulerRequest{
		Name:        "Test Payment Updated",
		Description: "Test Payment Description Updated",
		ProviderId:  uuid.New(),
		Sum:         1000,
		Spec:        scheduler.DAILY,
	}
}

func GeneratePaymentScheduler(houseId uuid.UUID, userId uuid.UUID, providerId uuid.UUID) model.PaymentScheduler {
	return model.PaymentScheduler{
		Id:          uuid.New(),
		Name:        "Test Payment",
		Description: "Test Payment Description",
		HouseId:     houseId,
		UserId:      userId,
		ProviderId:  providerId,
		Sum:         1000,
		Spec:        scheduler.DAILY,
	}
}

func GeneratePaymentSchedulerResponse(id uuid.UUID) model.PaymentSchedulerDto {
	return model.PaymentSchedulerDto{
		Id:          id,
		Name:        "Test Payment",
		Description: "Test Payment Description",
		HouseId:     HouseId,
		UserId:      UserId,
		ProviderId:  ProviderId,
		Sum:         1000,
		Spec:        scheduler.DAILY,
	}
}
