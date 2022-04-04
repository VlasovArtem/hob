package service

import (
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/common"
	"github.com/VlasovArtem/hob/src/common/database"
	"github.com/VlasovArtem/hob/src/common/dependency"
	int_errors "github.com/VlasovArtem/hob/src/common/int-errors"
	houses "github.com/VlasovArtem/hob/src/house/service"
	"github.com/VlasovArtem/hob/src/payment/model"
	"github.com/VlasovArtem/hob/src/payment/repository"
	providers "github.com/VlasovArtem/hob/src/provider/service"
	users "github.com/VlasovArtem/hob/src/user/service"
	"github.com/google/uuid"
	"time"
)

var defaultUUID = uuid.UUID{}

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
		dependency.FindRequiredDependency[users.UserServiceObject, users.UserService](factory),
		dependency.FindRequiredDependency[houses.HouseServiceObject, houses.HouseService](factory),
		dependency.FindRequiredDependency[providers.ProviderServiceObject, providers.ProviderService](factory),
		dependency.FindRequiredDependency[repository.PaymentRepositoryObject, repository.PaymentRepository](factory),
	)
}

type PaymentService interface {
	Add(request model.CreatePaymentRequest) (model.PaymentDto, error)
	AddBatch(request model.CreatePaymentBatchRequest) ([]model.PaymentDto, error)
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
		return response, fmt.Errorf("user with id %s not found", request.UserId)
	}
	if !p.houseService.ExistsById(request.HouseId) {
		return response, fmt.Errorf("house with id %s not found", request.HouseId)
	}

	if request.ProviderId != defaultUUID {
		if !p.providerService.ExistsById(request.ProviderId) {
			return response, fmt.Errorf("provider with id %s not found", request.ProviderId)
		}
	}

	payment, err := p.paymentRepository.Create(request.ToEntity())

	return payment.ToDto(), err
}

func (p *PaymentServiceObject) AddBatch(request model.CreatePaymentBatchRequest) (response []model.PaymentDto, err error) {
	if len(request.Payments) == 0 {
		return make([]model.PaymentDto, 0), nil
	}

	userIds := make(map[uuid.UUID]bool)
	houseIds := make(map[uuid.UUID]bool)
	providerIds := make(map[uuid.UUID]bool)

	entities := common.MapSlice(request.Payments, func(paymentRequest model.CreatePaymentRequest) model.Payment {
		userIds[paymentRequest.UserId] = true
		houseIds[paymentRequest.HouseId] = true
		if paymentRequest.ProviderId != defaultUUID {
			providerIds[paymentRequest.ProviderId] = true
		}

		return paymentRequest.ToEntity()
	})

	builder := int_errors.NewBuilder()

	for userId, _ := range userIds {
		if !p.userService.ExistsById(userId) {
			builder.WithDetail(fmt.Sprintf("user with id %s not found", userId))
		}
	}

	for houseId, _ := range houseIds {
		if !p.houseService.ExistsById(houseId) {
			builder.WithDetail(fmt.Sprintf("house with id %s not found", houseId))
		}
	}

	for providerId, _ := range providerIds {
		if !p.providerService.ExistsById(providerId) {
			builder.WithDetail(fmt.Sprintf("provider with id %s not found", providerId))
		}
	}

	if builder.HasErrors() {
		return nil, int_errors.NewErrResponse(builder.WithMessage("Create payment batch failed"))
	}

	if batch, err := p.paymentRepository.CreateBatch(entities); err != nil {
		return response, err
	} else {
		return common.MapSlice(batch, model.EntityToDto), nil
	}
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
