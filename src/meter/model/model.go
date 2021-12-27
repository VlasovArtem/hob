package model

import (
	"encoding/json"
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
}

type CreateMeterRequest struct {
	Name        string
	Details     map[string]float64
	Description string
	PaymentId   uuid.UUID
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
