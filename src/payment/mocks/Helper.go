package mocks

import (
	"github.com/VlasovArtem/hob/src/payment/model"
	"github.com/VlasovArtem/hob/src/test/testhelper"
	"github.com/google/uuid"
	"strconv"
	"time"
)

var (
	HouseId    = testhelper.ParseUUID("d077adaa-00d7-4e80-ac86-57512267505d")
	UserId     = testhelper.ParseUUID("ad2c5035-6745-48d0-9eee-fd22f5dae8e0")
	ProviderId = testhelper.ParseUUID("a949dbdc-7a7c-4dd5-b224-db536a579d5d")
	Date       = time.Date(2021, time.January, 1, 0, 0, 0, 0, time.Local)
)

func GenerateCreatePaymentRequest() model.CreatePaymentRequest {
	return model.CreatePaymentRequest{
		Name:        "Test Payment",
		Description: "Test Payment Description",
		HouseId:     HouseId,
		UserId:      UserId,
		ProviderId:  &ProviderId,
		Date:        Date,
		Sum:         float32(200 + (time.Now().Unix() % 10)),
	}
}

func GenerateCreatePaymentBatchRequest(number int) model.CreatePaymentBatchRequest {
	request := model.CreatePaymentBatchRequest{}
	for i := 0; i < number; i++ {
		incomeRequest := GenerateCreatePaymentRequest()
		incomeRequest.Name = "Payment Name #" + strconv.Itoa(i)
		request.Payments = append(request.Payments, incomeRequest)
	}
	return request
}

func GenerateUpdatePaymentRequest() model.UpdatePaymentRequest {
	return model.UpdatePaymentRequest{
		Name:        "Test Payment",
		Description: "Test Payment Description",
		Date:        Date,
		Sum:         float32(200 + (time.Now().Unix() % 10)),
		ProviderId:  &ProviderId,
	}
}

func GeneratePayment(houseId uuid.UUID, userId uuid.UUID, providerId uuid.UUID) model.Payment {
	return model.Payment{
		Id:          uuid.New(),
		Name:        "Test Payment",
		Description: "Test Payment Description",
		HouseId:     houseId,
		UserId:      userId,
		Date:        Date,
		Sum:         float32(200 + (time.Now().Unix() % 10)),
		ProviderId:  &providerId,
	}
}

func GeneratePaymentWithoutProvider(houseId uuid.UUID, userId uuid.UUID) model.Payment {
	return model.Payment{
		Id:          uuid.New(),
		Name:        "Test Payment",
		Description: "Test Payment Description",
		HouseId:     houseId,
		UserId:      userId,
		Date:        Date,
		Sum:         float32(200 + (time.Now().Unix() % 10)),
		ProviderId:  nil,
	}
}

func GeneratePaymentResponse() model.PaymentDto {
	return model.PaymentDto{
		Id:          uuid.New(),
		Name:        "Test Payment",
		Description: "Test Payment Description",
		HouseId:     HouseId,
		UserId:      UserId,
		ProviderId:  &ProviderId,
		Date:        Date,
		Sum:         float32(200 + (time.Now().Unix() % 10)),
	}
}
