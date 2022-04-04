package model

import (
	houseModel "github.com/VlasovArtem/hob/src/house/model"
	providerModel "github.com/VlasovArtem/hob/src/provider/model"
	userModel "github.com/VlasovArtem/hob/src/user/model"
	"github.com/google/uuid"
	"time"
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
	ProviderId  uuid.UUID
	Provider    providerModel.Provider `gorm:"foreignKey:ProviderId"`
}

type CreatePaymentRequest struct {
	Name        string
	Description string
	HouseId     uuid.UUID
	UserId      uuid.UUID
	ProviderId  uuid.UUID
	Date        time.Time
	Sum         float32
}

type CreatePaymentBatchRequest struct {
	Payments []CreatePaymentRequest
}

type UpdatePaymentRequest struct {
	Name        string
	Description string
	Date        time.Time
	Sum         float32
	ProviderId  uuid.UUID
}

type PaymentDto struct {
	Id          uuid.UUID
	Name        string
	Description string
	HouseId     uuid.UUID
	UserId      uuid.UUID
	ProviderId  uuid.UUID
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
		ProviderId:  p.ProviderId,
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
		ProviderId:  c.ProviderId,
		Date:        c.Date,
		Sum:         c.Sum,
	}
}

func (u UpdatePaymentRequest) UpdateToEntity(id uuid.UUID) Payment {
	return Payment{
		Id:          id,
		Name:        u.Name,
		Description: u.Description,
		ProviderId:  u.ProviderId,
		Date:        u.Date,
		Sum:         u.Sum,
	}
}

func EntityToDto(entity Payment) PaymentDto {
	return entity.ToDto()
}
