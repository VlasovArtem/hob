package model

import (
	"github.com/google/uuid"
	houseModel "house/model"
	"time"
	userModel "user/model"
)

type Payment struct {
	Id          uuid.UUID `gorm:"primarykey"`
	Name        string
	Description string
	HouseId     uuid.UUID
	UserId      uuid.UUID
	Date        time.Time
	Sum         float32
	User        userModel.User   `gorm:"foreignKey:UserId"`
	House       houseModel.House `gorm:"foreignKey:HouseId"`
}

type CreatePaymentRequest struct {
	Name        string
	Description string
	HouseId     uuid.UUID
	UserId      uuid.UUID
	Date        time.Time
	Sum         float32
}

type PaymentDto struct {
	Id          uuid.UUID
	Name        string
	Description string
	HouseId     uuid.UUID
	UserId      uuid.UUID
	Date        time.Time
	Sum         float32
}

func (p Payment) ToDto() PaymentDto {
	return PaymentDto{
		Id:          p.Id,
		Name:        p.Name,
		Description: p.Description,
		HouseId:     p.HouseId,
		UserId:      p.UserId,
		Date:        p.Date,
		Sum:         p.Sum,
	}
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
