package service

import (
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/common/database"
	"github.com/VlasovArtem/hob/src/common/dependency"
	int_errors "github.com/VlasovArtem/hob/src/common/int-errors"
	houses "github.com/VlasovArtem/hob/src/house/service"
	paymentModel "github.com/VlasovArtem/hob/src/payment/model"
	"github.com/VlasovArtem/hob/src/payment/scheduler/model"
	"github.com/VlasovArtem/hob/src/payment/scheduler/repository"
	payments "github.com/VlasovArtem/hob/src/payment/service"
	providers "github.com/VlasovArtem/hob/src/provider/service"
	"github.com/VlasovArtem/hob/src/scheduler"
	users "github.com/VlasovArtem/hob/src/user/service"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"time"
)

type PaymentSchedulerServiceObject struct {
	userService      users.UserService
	houseService     houses.HouseService
	paymentService   payments.PaymentService
	providerService  providers.ProviderService
	serviceScheduler scheduler.ServiceScheduler
	repository       repository.PaymentSchedulerRepository
}

func NewPaymentSchedulerService(
	userService users.UserService,
	houseService houses.HouseService,
	paymentService payments.PaymentService,
	providerService providers.ProviderService,
	serviceScheduler scheduler.ServiceScheduler,
	repository repository.PaymentSchedulerRepository,
) PaymentSchedulerService {
	return &PaymentSchedulerServiceObject{
		userService:      userService,
		houseService:     houseService,
		paymentService:   paymentService,
		providerService:  providerService,
		serviceScheduler: serviceScheduler,
		repository:       repository,
	}
}

func (p *PaymentSchedulerServiceObject) Initialize(factory dependency.DependenciesProvider) interface{} {
	return NewPaymentSchedulerService(
		factory.FindRequiredByType(users.UserServiceType).(users.UserService),
		factory.FindRequiredByType(houses.HouseServiceType).(houses.HouseService),
		factory.FindRequiredByType(payments.PaymentServiceType).(payments.PaymentService),
		factory.FindRequiredByType(providers.ProviderServiceType).(providers.ProviderService),
		factory.FindRequiredByType(scheduler.SchedulerServiceType).(scheduler.ServiceScheduler),
		factory.FindRequiredByType(repository.PaymentSchedulerRepositoryType).(repository.PaymentSchedulerRepository),
	)
}

type PaymentSchedulerService interface {
	Add(request model.CreatePaymentSchedulerRequest) (model.PaymentSchedulerDto, error)
	Remove(id uuid.UUID) error
	FindById(id uuid.UUID) (model.PaymentSchedulerDto, error)
	FindByHouseId(houseId uuid.UUID) []model.PaymentSchedulerDto
	FindByUserId(userId uuid.UUID) []model.PaymentSchedulerDto
	FindByProviderId(providerId uuid.UUID) []model.PaymentSchedulerDto
	Update(id uuid.UUID, request model.UpdatePaymentSchedulerRequest) error
}

func (p *PaymentSchedulerServiceObject) Add(request model.CreatePaymentSchedulerRequest) (response model.PaymentSchedulerDto, err error) {
	if err = p.validateCreateRequest(request); err != nil {
		return response, err
	}

	entity := request.ToEntity()

	if entity, err = p.repository.Create(entity); err != nil {
		return response, err
	} else if _, err = p.serviceScheduler.Add(entity.Id, string(entity.Spec), p.schedulerFunc(entity)); err != nil {
		p.repository.DeleteById(entity.Id)

		return response, err
	} else {
		return entity.ToDto(), err
	}
}

func (p *PaymentSchedulerServiceObject) validateCreateRequest(request model.CreatePaymentSchedulerRequest) error {
	if request.Sum <= 0 {
		return errors.New("sum should not be zero of negative")
	}
	if !p.userService.ExistsById(request.UserId) {
		return int_errors.NewErrNotFound("user with id %s in not exists", request.UserId)
	}
	if !p.houseService.ExistsById(request.HouseId) {
		return int_errors.NewErrNotFound("house with id %s in not exists", request.HouseId)
	}
	if !p.providerService.ExistsById(request.ProviderId) {
		return int_errors.NewErrNotFound("provider with id %s in not exists", request.ProviderId)
	}
	if request.Spec == "" {
		return errors.New("scheduler configuration not provided")
	}
	return nil
}

func (p *PaymentSchedulerServiceObject) Remove(id uuid.UUID) error {
	if !p.repository.ExistsById(id) {
		return int_errors.NewErrNotFound("payment scheduler with id %s not found", id)
	} else {
		if err := p.serviceScheduler.Remove(id); err != nil {
			log.Error().Err(err).Msg("")
		}
		p.repository.DeleteById(id)
	}
	return nil
}

func (p *PaymentSchedulerServiceObject) FindById(id uuid.UUID) (response model.PaymentSchedulerDto, err error) {
	if paymentScheduler, err := p.repository.FindById(id); err != nil {
		return response, database.HandlerFindError(err, fmt.Sprintf("payment scheduler with id %s not found", id))
	} else {
		return paymentScheduler.ToDto(), err
	}
}

func (p *PaymentSchedulerServiceObject) FindByHouseId(houseId uuid.UUID) []model.PaymentSchedulerDto {
	return p.repository.FindByHouseId(houseId)
}

func (p *PaymentSchedulerServiceObject) FindByUserId(userId uuid.UUID) []model.PaymentSchedulerDto {
	return p.repository.FindByUserId(userId)
}

func (p *PaymentSchedulerServiceObject) FindByProviderId(providerId uuid.UUID) []model.PaymentSchedulerDto {
	return p.repository.FindByProviderId(providerId)
}

func (p *PaymentSchedulerServiceObject) Update(id uuid.UUID, request model.UpdatePaymentSchedulerRequest) error {
	if err, _ := p.validateUpdateRequest(id, request); err != nil {
		return err
	}

	updatedEntity, err := p.repository.Update(id, request)

	if err != nil {
		return err
	}

	if _, err := p.serviceScheduler.Update(updatedEntity.Id, string(updatedEntity.Spec), p.schedulerFunc(updatedEntity)); err != nil {
		p.repository.DeleteById(updatedEntity.Id)

		return err
	}
	return nil
}

func (p *PaymentSchedulerServiceObject) validateUpdateRequest(id uuid.UUID, request model.UpdatePaymentSchedulerRequest) (error, bool) {
	if request.Sum <= 0 {
		return errors.New("sum should not be zero of negative"), true
	}
	if !p.repository.ExistsById(id) {
		return int_errors.NewErrNotFound("payment schedule with id %s not found", id), true
	}
	if !p.providerService.ExistsById(request.ProviderId) {
		return int_errors.NewErrNotFound("provider with id %s not found", request.ProviderId), true
	}
	if request.Spec == "" {
		return errors.New("scheduler configuration not provided"), true
	}
	return nil, false
}

func (p *PaymentSchedulerServiceObject) schedulerFunc(payment model.PaymentScheduler) func() {
	return func() {
		if _, err := p.paymentService.Add(
			paymentModel.CreatePaymentRequest{
				Name:        payment.Name,
				Description: payment.Description,
				HouseId:     payment.HouseId,
				UserId:      payment.UserId,
				ProviderId:  payment.ProviderId,
				Date:        time.Now(),
				Sum:         payment.Sum,
			},
		); err != nil {
			log.Error().Err(err).Msg("")
		} else {
			log.Info().Msgf("New payment added to the house %s and user %s via scheduler %s", payment.HouseId, payment.UserId, payment.Id)
		}
	}
}
