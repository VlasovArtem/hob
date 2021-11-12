package model

import (
	"github.com/google/uuid"
	"time"
)

type Payment struct {
	Id          uuid.UUID
	Name        string
	Description string
	HouseId     uuid.UUID
	UserId      uuid.UUID
	Date        time.Time
	Sum         float32
}

type CreatePaymentRequest struct {
	Name        string
	Description string
	HouseId     uuid.UUID
	UserId      uuid.UUID
	Date        time.Time
	Sum         float32
}

type PaymentResponse struct {
	Payment
}

func (p Payment) ToResponse() PaymentResponse {
	return PaymentResponse{p}
}

func (c CreatePaymentRequest) ToEntity() Payment {
	return Payment{
		Id:          uuid.New(),
		Name:        c.Name,
		Description: c.Description,
		HouseId:     c.HouseId,
		UserId:      c.UserId,
		Date:        c.Date,
		Sum:         c.Sum,
	}
}
