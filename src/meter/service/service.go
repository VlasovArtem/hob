package service

import (
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/common/dependency"
	int_errors "github.com/VlasovArtem/hob/src/common/int-errors"
	houseService "github.com/VlasovArtem/hob/src/house/service"
	"github.com/VlasovArtem/hob/src/meter/model"
	"github.com/VlasovArtem/hob/src/meter/repository"
	paymentService "github.com/VlasovArtem/hob/src/payment/service"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"reflect"
)

var MeterServiceType = reflect.TypeOf(MeterServiceObject{})

type MeterServiceObject struct {
	paymentService paymentService.PaymentService
	houseService   houseService.HouseService
	repository     repository.MeterRepository
}

func NewMeterService(paymentService paymentService.PaymentService, houseService houseService.HouseService, repository repository.MeterRepository) MeterService {
	return &MeterServiceObject{paymentService, houseService, repository}
}

func (m *MeterServiceObject) Initialize(factory dependency.DependenciesProvider) interface{} {
	return NewMeterService(
		factory.FindRequiredByType(paymentService.PaymentServiceType).(paymentService.PaymentService),
		factory.FindRequiredByType(houseService.HouseServiceType).(houseService.HouseService),
		factory.FindRequiredByType(repository.MeterRepositoryType).(repository.MeterRepository),
	)
}

type MeterService interface {
	Add(request model.CreateMeterRequest) (model.MeterDto, error)
	Update(id uuid.UUID, request model.UpdateMeterRequest) error
	DeleteById(id uuid.UUID) error
	FindById(id uuid.UUID) (model.MeterDto, error)
	FindByPaymentId(id uuid.UUID) (model.MeterDto, error)
	FindByHouseId(id uuid.UUID) []model.MeterDto
}

func (m *MeterServiceObject) Add(request model.CreateMeterRequest) (response model.MeterDto, err error) {
	if !m.paymentService.ExistsById(request.PaymentId) {
		return response, fmt.Errorf("payment with id %s not found", request.PaymentId)
	}

	if entity, err := m.repository.Create(request.ToEntity()); err != nil {
		return response, err
	} else {
		return entity.ToDto(), nil
	}
}

func (m *MeterServiceObject) Update(id uuid.UUID, request model.UpdateMeterRequest) error {
	if !m.repository.ExistsById(id) {
		return int_errors.NewErrNotFound("meter with id %s not found", id)
	}

	return m.repository.Update(id, request.ToEntity())
}

func (m *MeterServiceObject) DeleteById(id uuid.UUID) error {
	if !m.repository.ExistsById(id) {
		return int_errors.NewErrNotFound("meter with id %s not found", id)
	}

	return m.repository.DeleteById(id)
}

func (m *MeterServiceObject) FindById(id uuid.UUID) (model.MeterDto, error) {
	if meter, err := m.repository.FindById(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.MeterDto{}, fmt.Errorf("meter with id %s in not exists", id)
		} else {
			return model.MeterDto{}, err
		}
	} else {
		return meter.ToDto(), nil
	}
}

func (m *MeterServiceObject) FindByPaymentId(id uuid.UUID) (model.MeterDto, error) {
	if meter, err := m.repository.FindByPaymentId(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.MeterDto{}, fmt.Errorf("meter with payment id %s in not exists", id)
		} else {
			return model.MeterDto{}, err
		}
	} else {
		return meter.ToDto(), nil
	}
}

func (m *MeterServiceObject) FindByHouseId(id uuid.UUID) (response []model.MeterDto) {
	for _, meter := range m.repository.FindByHouseId(id) {
		response = append(response, meter.ToDto())
	}

	if response == nil {
		return make([]model.MeterDto, 0)
	}

	return response
}
