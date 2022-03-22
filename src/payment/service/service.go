package service

import (
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/common/database"
	"github.com/VlasovArtem/hob/src/common/dependency"
	houses "github.com/VlasovArtem/hob/src/house/service"
	"github.com/VlasovArtem/hob/src/payment/model"
	"github.com/VlasovArtem/hob/src/payment/repository"
	providers "github.com/VlasovArtem/hob/src/provider/service"
	users "github.com/VlasovArtem/hob/src/user/service"
	"github.com/google/uuid"
	"reflect"
	"time"
)

var defaultUUID = uuid.UUID{}

var PaymentServiceType = reflect.TypeOf(PaymentServiceObject{})

type PaymentServiceObject struct {
	userService       users.UserService
	houseService      houses.HouseService
	providerService   providers.ProviderService
	paymentRepository repository.PaymentRepository
}

func NewPaymentService(
	userService users.UserService,
	houseService houses.HouseService,
	providerService providers.ProviderService,
	paymentRepository repository.PaymentRepository) PaymentService {
	return &PaymentServiceObject{
		userService:       userService,
		houseService:      houseService,
		providerService:   providerService,
		paymentRepository: paymentRepository,
	}
}

func (p *PaymentServiceObject) Initialize(factory dependency.DependenciesProvider) any {
	return NewPaymentService(
		factory.FindRequiredByType(users.UserServiceType).(users.UserService),
		factory.FindRequiredByType(houses.HouseServiceType).(houses.HouseService),
		factory.FindRequiredByType(providers.ProviderServiceType).(providers.ProviderService),
		factory.FindRequiredByType(repository.PaymentRepositoryType).(repository.PaymentRepository),
	)
}

type PaymentService interface {
	Add(request model.CreatePaymentRequest) (model.PaymentDto, error)
	FindById(id uuid.UUID) (model.PaymentDto, error)
	FindByHouseId(id uuid.UUID) []model.PaymentDto
	FindByUserId(id uuid.UUID) []model.PaymentDto
	FindByProviderId(id uuid.UUID) []model.PaymentDto
	ExistsById(id uuid.UUID) bool
	DeleteById(id uuid.UUID) error
	Update(id uuid.UUID, request model.UpdatePaymentRequest) error
}

func (p *PaymentServiceObject) Add(request model.CreatePaymentRequest) (response model.PaymentDto, err error) {
	if !p.userService.ExistsById(request.UserId) {
		return response, fmt.Errorf("user with id %s in not exists", request.UserId)
	}
	if !p.houseService.ExistsById(request.HouseId) {
		return response, fmt.Errorf("house with id %s in not exists", request.HouseId)
	}

	if request.ProviderId != defaultUUID {
		if !p.providerService.ExistsById(request.ProviderId) {
			return response, fmt.Errorf("provider with id %s in not exists", request.ProviderId)
		}
	}

	payment, err := p.paymentRepository.Create(request.CreateToEntity())

	return payment.ToDto(), err
}

func (p *PaymentServiceObject) FindById(id uuid.UUID) (model.PaymentDto, error) {
	if payment, err := p.paymentRepository.FindById(id); err != nil {
		return model.PaymentDto{}, database.HandlerFindError(err, fmt.Sprintf("payment with id %s not found", id))
	} else {
		return payment.ToDto(), nil
	}
}

func (p *PaymentServiceObject) FindByHouseId(houseId uuid.UUID) []model.PaymentDto {
	return p.paymentRepository.FindByHouseId(houseId)
}

func (p *PaymentServiceObject) FindByUserId(userId uuid.UUID) []model.PaymentDto {
	return p.paymentRepository.FindByUserId(userId)
}

func (p *PaymentServiceObject) FindByProviderId(id uuid.UUID) []model.PaymentDto {
	return p.paymentRepository.FindByProviderId(id)
}

func (p *PaymentServiceObject) ExistsById(id uuid.UUID) bool {
	return p.paymentRepository.ExistsById(id)
}

func (p *PaymentServiceObject) DeleteById(id uuid.UUID) error {
	if !p.ExistsById(id) {
		return fmt.Errorf("payment with id %s not found", id)
	}
	return p.paymentRepository.DeleteById(id)
}

func (p *PaymentServiceObject) Update(id uuid.UUID, request model.UpdatePaymentRequest) error {
	if !p.ExistsById(id) {
		return fmt.Errorf("payment with id %s not found", id)
	}
	if !p.providerService.ExistsById(request.ProviderId) {
		return fmt.Errorf("provider with id %s not found", request.ProviderId)
	}
	if request.Date.After(time.Now()) {
		return errors.New("date should not be after current date")
	}
	return p.paymentRepository.Update(request.UpdateToEntity(id))
}
