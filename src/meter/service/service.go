package service

import (
	"fmt"
	"github.com/VlasovArtem/hob/src/common/database"
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/int-errors"
	"github.com/VlasovArtem/hob/src/common/transactional"
	"github.com/VlasovArtem/hob/src/meter/model"
	"github.com/VlasovArtem/hob/src/meter/repository"
	paymentService "github.com/VlasovArtem/hob/src/payment/service"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MeterServiceStr struct {
	paymentService paymentService.PaymentService
	repository     repository.MeterRepository
}

func NewMeterService(paymentService paymentService.PaymentService, repository repository.MeterRepository) MeterService {
	return &MeterServiceStr{paymentService, repository}
}

func (m *MeterServiceStr) GetRequiredDependencies() []dependency.Requirements {
	return []dependency.Requirements{
		dependency.FindNameAndType(paymentService.PaymentServiceStr{}),
		dependency.FindNameAndType(repository.MeterRepositoryStr{}),
	}
}

func (m *MeterServiceStr) Initialize(factory dependency.DependenciesProvider) any {
	return NewMeterService(
		dependency.FindRequiredDependency[paymentService.PaymentServiceStr, paymentService.PaymentService](factory),
		dependency.FindRequiredDependency[repository.MeterRepositoryStr, repository.MeterRepository](factory),
	)
}

type MeterService interface {
	transactional.Transactional[MeterService]
	Add(request model.CreateMeterRequest) (model.MeterDto, error)
	Update(id uuid.UUID, request model.UpdateMeterRequest) error
	DeleteById(id uuid.UUID) error
	FindById(id uuid.UUID) (model.MeterDto, error)
	FindByPaymentId(id uuid.UUID) (model.MeterDto, error)
}

func (m *MeterServiceStr) Add(request model.CreateMeterRequest) (response model.MeterDto, err error) {
	if !m.paymentService.ExistsById(request.PaymentId) {
		return response, fmt.Errorf("payment with id %s not found", request.PaymentId)
	}

	entity := request.ToEntity()
	if err := m.repository.Create(&entity); err != nil {
		return response, err
	} else {
		return entity.ToDto(), nil
	}
}

func (m *MeterServiceStr) Update(id uuid.UUID, request model.UpdateMeterRequest) error {
	if !m.repository.Exists(id) {
		return int_errors.NewErrNotFound("meter with id %s not found", id)
	}

	return m.repository.Update(id, request.ToEntity())
}

func (m *MeterServiceStr) DeleteById(id uuid.UUID) error {
	if !m.repository.Exists(id) {
		return int_errors.NewErrNotFound("meter with id %s not found", id)
	}

	return m.repository.Delete(id)
}

func (m *MeterServiceStr) FindById(id uuid.UUID) (dto model.MeterDto, err error) {
	if meter, err := m.repository.Find(id); err != nil {
		return dto, database.HandlerFindError(err, "meter with id %s in not exists", id)
	} else {
		return meter.ToDto(), err
	}
}

func (m *MeterServiceStr) FindByPaymentId(id uuid.UUID) (dto model.MeterDto, err error) {
	if meter, err := m.repository.FindByPaymentId(id); err != nil {
		return dto, database.HandlerFindError(err, "meter with payment id %s in not exists", id)
	} else {
		return meter.ToDto(), err
	}
}

func (m *MeterServiceStr) Transactional(tx *gorm.DB) MeterService {
	return &MeterServiceStr{
		paymentService: m.paymentService.Transactional(tx),
		repository:     m.repository.Transactional(tx),
	}
}
