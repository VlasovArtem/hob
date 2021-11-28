package model

import (
	"encoding/json"
	"github.com/google/uuid"
	houseModel "house/model"
	"payment/model"
)

type Meter struct {
	Id          uuid.UUID `gorm:"primarykey;type:uuid"`
	Name        string
	Details     string
	Description string
	PaymentId   uuid.UUID     `gorm:"index:idx_payment_id"`
	Payment     model.Payment `gorm:"foreignKey:PaymentId"`
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

type MeterResponse struct {
	Id          uuid.UUID
	Name        string
	Details     map[string]float64
	Description string
	PaymentId   uuid.UUID
	HouseId     uuid.UUID
}

func (m Meter) ToResponse() MeterResponse {
	details := map[string]float64{}

	_ = json.Unmarshal([]byte(m.Details), &details)

	return MeterResponse{
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
		Details:     string(marshal),
		Description: c.Description,
		PaymentId:   c.PaymentId,
		HouseId:     c.HouseId,
	}
}
