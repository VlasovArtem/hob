package mocks

import (
	"github.com/google/uuid"
	"payment/model"
	"test/testhelper"
	"time"
)

var (
	HouseId = testhelper.ParseUUID("d077adaa-00d7-4e80-ac86-57512267505d")
	UserId  = testhelper.ParseUUID("ad2c5035-6745-48d0-9eee-fd22f5dae8e0")
	Date    = time.Date(2021, time.January, 1, 0, 0, 0, 0, time.Local)
)

func GenerateCreatePaymentRequest() model.CreatePaymentRequest {
	return model.CreatePaymentRequest{
		Name:        "Test Payment",
		Description: "Test Payment Description",
		HouseId:     HouseId,
		UserId:      UserId,
		Date:        Date,
		Sum:         1000,
	}
}

func GeneratePayment(houseId uuid.UUID, userId uuid.UUID) model.Payment {
	return model.Payment{
		Id:          uuid.New(),
		Name:        "Test Payment",
		Description: "Test Payment Description",
		HouseId:     houseId,
		UserId:      userId,
		Date:        Date,
		Sum:         1000,
	}
}

func GeneratePaymentResponse() model.PaymentDto {
	return model.PaymentDto{
		Payment: model.Payment{
			Id:          uuid.New(),
			Name:        "Test Payment",
			Description: "Test Payment Description",
			HouseId:     HouseId,
			UserId:      UserId,
			Date:        Date,
			Sum:         1000,
		},
	}
}
