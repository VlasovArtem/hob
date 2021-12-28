package service

import (
	"fmt"
	"github.com/VlasovArtem/hob/src/common/database"
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/int-errors"
	"github.com/VlasovArtem/hob/src/meter/model"
	"github.com/VlasovArtem/hob/src/meter/repository"
	paymentService "github.com/VlasovArtem/hob/src/payment/service"
	"github.com/google/uuid"
	"reflect"
)

var MeterServiceType = reflect.TypeOf(MeterServiceObject{})

type MeterServiceObject struct {
	paymentService paymentService.PaymentService
	repository     repository.MeterRepository
}

func NewMeterService(paymentService paymentService.PaymentService, repository repository.MeterRepository) MeterService {
	return &MeterServiceObject{paymentService, repository}
}

func (m *MeterServiceObject) Initialize(factory dependency.DependenciesProvider) any {
	return NewMeterService(
		factory.FindRequiredByType(paymentService.PaymentServiceType).(paymentService.PaymentService),
		factory.FindRequiredByType(repository.MeterRepositoryType).(repository.MeterRepository),
	)
}

type MeterService interface {
	Add(request model.CreateMeterRequest) (model.MeterDto, error)
	Update(id uuid.UUID, request model.UpdateMeterRequest) error
	DeleteById(id uuid.UUID) error
	FindById(id uuid.UUID) (model.MeterDto, error)
	FindByPaymentId(id uuid.UUID) (model.MeterDto, error)
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

func (m *MeterServiceObject) FindById(id uuid.UUID) (dto model.MeterDto, err error) {
	if meter, err := m.repository.FindById(id); err != nil {
		return dto, database.HandlerFindError(err, "meter with id %s in not exists", id)
	} else {
		return meter.ToDto(), err
	}
}

func (m *MeterServiceObject) FindByPaymentId(id uuid.UUID) (dto model.MeterDto, err error) {
	if meter, err := m.repository.FindByPaymentId(id); err != nil {
		return dto, database.HandlerFindError(err, "meter with payment id %s in not exists", id)
	} else {
		return meter.ToDto(), err
	}
}
