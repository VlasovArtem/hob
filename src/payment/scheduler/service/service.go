package service

import (
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/common/database"
	"github.com/VlasovArtem/hob/src/common/dependency"
	intErrors "github.com/VlasovArtem/hob/src/common/int-errors"
	"github.com/VlasovArtem/hob/src/common/transactional"
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
	"gorm.io/gorm"
	"time"
)

type PaymentSchedulerServiceStr struct {
	userService      users.UserService
	houseService     houses.HouseService
	paymentService   payments.PaymentService
	providerService  providers.ProviderService
	serviceScheduler scheduler.ServiceScheduler
	repository       repository.PaymentSchedulerRepository
}

func (p *PaymentSchedulerServiceStr) Transactional(tx *gorm.DB) PaymentSchedulerService {
	return &PaymentSchedulerServiceStr{
		userService:      p.userService.Transactional(tx),
		houseService:     p.houseService.Transactional(tx),
		paymentService:   p.paymentService.Transactional(tx),
		providerService:  p.providerService.Transactional(tx),
		serviceScheduler: p.serviceScheduler,
		repository:       p.repository.Transactional(tx),
	}
}

func NewPaymentSchedulerService(
	userService users.UserService,
	houseService houses.HouseService,
	paymentService payments.PaymentService,
	providerService providers.ProviderService,
	serviceScheduler scheduler.ServiceScheduler,
	repository repository.PaymentSchedulerRepository,
) PaymentSchedulerService {
	return &PaymentSchedulerServiceStr{
		userService:      userService,
		houseService:     houseService,
		paymentService:   paymentService,
		providerService:  providerService,
		serviceScheduler: serviceScheduler,
		repository:       repository,
	}
}

func (p *PaymentSchedulerServiceStr) Initialize(factory dependency.DependenciesProvider) any {
	return NewPaymentSchedulerService(
		dependency.FindRequiredDependency[users.UserServiceObject, users.UserService](factory),
		dependency.FindRequiredDependency[houses.HouseServiceStr, houses.HouseService](factory),
		dependency.FindRequiredDependency[payments.PaymentServiceStr, payments.PaymentService](factory),
		dependency.FindRequiredDependency[providers.ProviderServiceStr, providers.ProviderService](factory),
		dependency.FindRequiredDependency[scheduler.SchedulerServiceObject, scheduler.ServiceScheduler](factory),
		dependency.FindRequiredDependency[repository.PaymentSchedulerRepositoryStr, repository.PaymentSchedulerRepository](factory),
	)
}

type PaymentSchedulerService interface {
	transactional.Transactional[PaymentSchedulerService]
	Add(request model.CreatePaymentSchedulerRequest) (model.PaymentSchedulerDto, error)
	Remove(id uuid.UUID) error
	FindById(id uuid.UUID) (model.PaymentSchedulerDto, error)
	FindByHouseId(houseId uuid.UUID) []model.PaymentSchedulerDto
	FindByUserId(userId uuid.UUID) []model.PaymentSchedulerDto
	FindByProviderId(providerId uuid.UUID) []model.PaymentSchedulerDto
	Update(id uuid.UUID, request model.UpdatePaymentSchedulerRequest) error
}

func (p *PaymentSchedulerServiceStr) Add(request model.CreatePaymentSchedulerRequest) (response model.PaymentSchedulerDto, err error) {
	if err = p.validateCreateRequest(request); err != nil {
		return
	}

	entity := request.ToEntity()

	if err = p.repository.Create(&entity); err != nil {
		return
	} else if _, err = p.serviceScheduler.Add(entity.Id, string(entity.Spec), p.schedulerFunc(entity)); err != nil {
		return response, p.repository.Delete(entity.Id)
	} else {
		return entity.ToDto(), err
	}
}

func (p *PaymentSchedulerServiceStr) validateCreateRequest(request model.CreatePaymentSchedulerRequest) error {
	if request.Sum <= 0 {
		return errors.New("sum should not be zero of negative")
	}
	if !p.userService.ExistsById(request.UserId) {
		return intErrors.NewErrNotFound("user with id %s in not exists", request.UserId)
	}
	if !p.houseService.ExistsById(request.HouseId) {
		return intErrors.NewErrNotFound("house with id %s in not exists", request.HouseId)
	}
	if !p.providerService.ExistsById(request.ProviderId) {
		return intErrors.NewErrNotFound("provider with id %s in not exists", request.ProviderId)
	}
	if request.Spec == "" {
		return errors.New("scheduler configuration not provided")
	}
	return nil
}

func (p *PaymentSchedulerServiceStr) Remove(id uuid.UUID) error {
	if !p.repository.Exists(id) {
		return intErrors.NewErrNotFound("payment scheduler with id %s not found", id)
	} else {
		if err := p.serviceScheduler.Remove(id); err != nil {
			log.Error().Err(err).Msg("")
		} else {
			return p.repository.Delete(id)
		}
	}
	return nil
}

func (p *PaymentSchedulerServiceStr) FindById(id uuid.UUID) (response model.PaymentSchedulerDto, err error) {
	if paymentScheduler, err := p.repository.Find(id); err != nil {
		return response, database.HandlerFindError(err, fmt.Sprintf("payment scheduler with id %s not found", id))
	} else {
		return paymentScheduler.ToDto(), err
	}
}

func (p *PaymentSchedulerServiceStr) FindByHouseId(houseId uuid.UUID) []model.PaymentSchedulerDto {
	return p.repository.FindByHouseId(houseId)
}

func (p *PaymentSchedulerServiceStr) FindByUserId(userId uuid.UUID) []model.PaymentSchedulerDto {
	return p.repository.FindByUserId(userId)
}

func (p *PaymentSchedulerServiceStr) FindByProviderId(providerId uuid.UUID) []model.PaymentSchedulerDto {
	return p.repository.FindByProviderId(providerId)
}

func (p *PaymentSchedulerServiceStr) Update(id uuid.UUID, request model.UpdatePaymentSchedulerRequest) error {
	if err, _ := p.validateUpdateRequest(id, request); err != nil {
		return err
	}

	err := p.repository.Update(id, request)

	if err != nil {
		return err
	}

	paymentScheduler, err := p.repository.Find(id)

	if err != nil {
		return err
	}

	if _, err := p.serviceScheduler.Update(id, string(paymentScheduler.Spec), p.schedulerFunc(paymentScheduler)); err != nil {
		return p.repository.Delete(id)
	}
	return nil
}

func (p *PaymentSchedulerServiceStr) validateUpdateRequest(id uuid.UUID, request model.UpdatePaymentSchedulerRequest) (error, bool) {
	if request.Sum <= 0 {
		return errors.New("sum should not be zero of negative"), true
	}
	if err := p.repository.Delete(id); err != nil {
		return intErrors.NewErrNotFound("payment schedule with id %s not found", id), true
	}
	if !p.providerService.ExistsById(request.ProviderId) {
		return intErrors.NewErrNotFound("provider with id %s not found", request.ProviderId), true
	}
	if request.Spec == "" {
		return errors.New("scheduler configuration not provided"), true
	}
	return nil, false
}

func (p *PaymentSchedulerServiceStr) schedulerFunc(payment model.PaymentScheduler) func() {
	return func() {
		if _, err := p.paymentService.Add(
			paymentModel.CreatePaymentRequest{
				Name:        payment.Name,
				Description: payment.Description,
				HouseId:     payment.HouseId,
				UserId:      payment.UserId,
				ProviderId:  &payment.ProviderId,
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
