package service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"meter/model"
	p "payment/service"
)

type meterServiceObject struct {
	paymentService p.PaymentService
	meters         map[uuid.UUID]model.Meter
	paymentMeter   map[uuid.UUID]model.Meter
}

func NewMeterService(service p.PaymentService) MeterService {
	return &meterServiceObject{
		paymentService: service,
		meters:         make(map[uuid.UUID]model.Meter),
		paymentMeter:   make(map[uuid.UUID]model.Meter),
	}
}

type MeterService interface {
	AddMeter(request model.CreateMeterRequest) (model.MeterResponse, error)
	FindById(id uuid.UUID) (model.MeterResponse, error)
	FindByPaymentId(id uuid.UUID) (model.MeterResponse, error)
}

func (m *meterServiceObject) AddMeter(request model.CreateMeterRequest) (response model.MeterResponse, err error) {
	if !m.paymentService.ExistsById(request.PaymentId) {
		return response, errors.New(fmt.Sprintf("payment with id %s in not exists", request.PaymentId))
	}

	entity := request.ToEntity()

	m.meters[entity.Id] = entity
	m.paymentMeter[entity.PaymentId] = entity

	return entity.ToResponse(), nil
}

func (m *meterServiceObject) FindById(id uuid.UUID) (response model.MeterResponse, err error) {
	if meter, ok := m.meters[id]; !ok {
		return response, errors.New(fmt.Sprintf("meter with id %s in not exists", id))
	} else {
		return meter.ToResponse(), nil
	}
}

func (m *meterServiceObject) FindByPaymentId(id uuid.UUID) (model.MeterResponse, error) {
	if response, ok := m.paymentMeter[id]; ok {
		return response.ToResponse(), nil
	}
	return model.MeterResponse{}, errors.New(fmt.Sprintf("meters with payment id %s not found", id))
}
