package model

import "github.com/google/uuid"

type Meter struct {
	Id          uuid.UUID
	Name        string
	Details     map[string]float64
	Description string
	PaymentId   uuid.UUID
}

type CreateMeterRequest struct {
	Name        string
	Details     map[string]float64
	Description string
	PaymentId   uuid.UUID
}

type MeterResponse struct {
	Meter
}

func (m Meter) ToResponse() MeterResponse {
	return MeterResponse{m}
}

func (c CreateMeterRequest) ToEntity() Meter {
	return Meter{
		Id:          uuid.New(),
		Name:        c.Name,
		Details:     c.Details,
		Description: c.Description,
		PaymentId:   c.PaymentId,
	}
}
