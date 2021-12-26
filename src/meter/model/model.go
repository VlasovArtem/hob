package model

import (
	"encoding/json"
	houseModel "github.com/VlasovArtem/hob/src/house/model"
	paymentModel "github.com/VlasovArtem/hob/src/payment/model"
	"github.com/google/uuid"
)

type Meter struct {
	Id          uuid.UUID `gorm:"primarykey;type:uuid"`
	Name        string
	Details     []byte
	Description string
	PaymentId   uuid.UUID            `gorm:"index:idx_payment_id"`
	Payment     paymentModel.Payment `gorm:"foreignKey:PaymentId"`
	HouseId     uuid.UUID
	House       houseModel.House `gorm:"foreignKey:HouseId"`
}

type CreateMeterRequest struct {
	Name        string
	Details     map[string]float64
	Description string
	PaymentId   uuid.UUID
	HouseId     uuid.UUID
}

type UpdateMeterRequest struct {
	Name        string
	Details     map[string]float64
	Description string
}

type MeterDto struct {
	Id          uuid.UUID
	Name        string
	Details     map[string]float64
	Description string
	PaymentId   uuid.UUID
	HouseId     uuid.UUID
}

func (m Meter) ToDto() MeterDto {
	details := map[string]float64{}

	_ = json.Unmarshal(m.Details, &details)

	return MeterDto{
		Id:          m.Id,
		Name:        m.Name,
		Details:     details,
		Description: m.Description,
		PaymentId:   m.PaymentId,
		HouseId:     m.HouseId,
	}
}

func (c CreateMeterRequest) ToEntity() Meter {
	marshal, _ := json.Marshal(c.Details)

	return Meter{
		Id:          uuid.New(),
		Name:        c.Name,
		Details:     marshal,
		Description: c.Description,
		PaymentId:   c.PaymentId,
		HouseId:     c.HouseId,
	}
}

func (c UpdateMeterRequest) ToEntity() Meter {
	marshal, _ := json.Marshal(c.Details)

	return Meter{
		Name:        c.Name,
		Details:     marshal,
		Description: c.Description,
	}
}
