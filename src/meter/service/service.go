package service

import (
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/common/dependency"
	houseService "github.com/VlasovArtem/hob/src/house/service"
	"github.com/VlasovArtem/hob/src/meter/model"
	"github.com/VlasovArtem/hob/src/meter/repository"
	paymentService "github.com/VlasovArtem/hob/src/payment/service"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MeterServiceObject struct {
	paymentService paymentService.PaymentService
	houseService   houseService.HouseService
	repository     repository.MeterRepository
}

func NewMeterService(paymentService paymentService.PaymentService, houseService houseService.HouseService, repository repository.MeterRepository) MeterService {
	return &MeterServiceObject{paymentService, houseService, repository}
}

func (m *MeterServiceObject) Initialize(factory dependency.DependenciesFactory) {
	factory.Add(
		NewMeterService(
			factory.FindRequiredByObject(paymentService.PaymentServiceObject{}).(paymentService.PaymentService),
			factory.FindRequiredByObject(houseService.HouseServiceObject{}).(houseService.HouseService),
			factory.FindRequiredByObject(repository.MeterRepositoryObject{}).(repository.MeterRepository),
		),
	)
}

type MeterService interface {
	Add(request model.CreateMeterRequest) (model.MeterResponse, error)
	FindById(id uuid.UUID) (model.MeterResponse, error)
	FindByPaymentId(id uuid.UUID) (model.MeterResponse, error)
	FindByHouseId(id uuid.UUID) []model.MeterResponse
}

func (m *MeterServiceObject) Add(request model.CreateMeterRequest) (response model.MeterResponse, err error) {
	if !m.paymentService.ExistsById(request.PaymentId) {
		return response, errors.New(fmt.Sprintf("payment with id %s in not exists", request.PaymentId))
	}
	if !m.houseService.ExistsById(request.HouseId) {
		return response, errors.New(fmt.Sprintf("house with id %s in not exists", request.HouseId))
	}

	if entity, err := m.repository.Create(request.ToEntity()); err != nil {
		return response, err
	} else {
		return entity.ToResponse(), nil
	}
}

func (m *MeterServiceObject) FindById(id uuid.UUID) (model.MeterResponse, error) {
	if meter, err := m.repository.FindById(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.MeterResponse{}, errors.New(fmt.Sprintf("meter with id %s in not exists", id))
		} else {
			return model.MeterResponse{}, err
		}
	} else {
		return meter.ToResponse(), nil
	}
}

func (m *MeterServiceObject) FindByPaymentId(id uuid.UUID) (model.MeterResponse, error) {
	if meter, err := m.repository.FindByPaymentId(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.MeterResponse{}, errors.New(fmt.Sprintf("meter with payment id %s in not exists", id))
		} else {
			return model.MeterResponse{}, err
		}
	} else {
		return meter.ToResponse(), nil
	}
}

func (m *MeterServiceObject) FindByHouseId(id uuid.UUID) (response []model.MeterResponse) {
	for _, meter := range m.repository.FindByHouseId(id) {
		response = append(response, meter.ToResponse())
	}

	if response == nil {
		return make([]model.MeterResponse, 0)
	}

	return response
}
