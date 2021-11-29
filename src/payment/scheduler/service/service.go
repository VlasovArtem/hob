package service

import (
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/common/database"
	"github.com/VlasovArtem/hob/src/common/dependency"
	houseService "github.com/VlasovArtem/hob/src/house/service"
	paymentModel "github.com/VlasovArtem/hob/src/payment/model"
	"github.com/VlasovArtem/hob/src/payment/scheduler/model"
	"github.com/VlasovArtem/hob/src/payment/scheduler/repository"
	paymentService "github.com/VlasovArtem/hob/src/payment/service"
	"github.com/VlasovArtem/hob/src/scheduler"
	userService "github.com/VlasovArtem/hob/src/user/service"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"time"
)

type PaymentSchedulerServiceObject struct {
	userService      userService.UserService
	houseService     houseService.HouseService
	paymentService   paymentService.PaymentService
	serviceScheduler scheduler.ServiceScheduler
	repository       repository.PaymentSchedulerRepository
}

func NewPaymentSchedulerService(
	userService userService.UserService,
	houseService houseService.HouseService,
	paymentService paymentService.PaymentService,
	serviceScheduler scheduler.ServiceScheduler,
	repository repository.PaymentSchedulerRepository,
) PaymentSchedulerService {
	return &PaymentSchedulerServiceObject{
		userService:      userService,
		houseService:     houseService,
		paymentService:   paymentService,
		serviceScheduler: serviceScheduler,
		repository:       repository,
	}
}

func (p *PaymentSchedulerServiceObject) Initialize(factory dependency.DependenciesFactory) {
	factory.Add(
		NewPaymentSchedulerService(
			factory.FindRequiredByObject(userService.UserServiceObject{}).(userService.UserService),
			factory.FindRequiredByObject(houseService.HouseServiceObject{}).(houseService.HouseService),
			factory.FindRequiredByObject(paymentService.PaymentServiceObject{}).(paymentService.PaymentService),
			factory.FindRequiredByObject(scheduler.SchedulerServiceObject{}).(scheduler.ServiceScheduler),
			factory.FindRequiredByObject(repository.PaymentSchedulerRepositoryObject{}).(repository.PaymentSchedulerRepository),
		),
	)
}

type PaymentSchedulerService interface {
	Add(request model.CreatePaymentSchedulerRequest) (model.PaymentSchedulerDto, error)
	Remove(id uuid.UUID) error
	FindById(id uuid.UUID) (model.PaymentSchedulerDto, error)
	FindByHouseId(houseId uuid.UUID) []model.PaymentSchedulerDto
	FindByUserId(userId uuid.UUID) []model.PaymentSchedulerDto
}

func (p *PaymentSchedulerServiceObject) Add(request model.CreatePaymentSchedulerRequest) (response model.PaymentSchedulerDto, err error) {
	if request.Sum <= 0 {
		return response, errors.New("sum should not be zero of negative")
	}
	if !p.userService.ExistsById(request.UserId) {
		return response, errors.New(fmt.Sprintf("user with id %s in not exists", request.UserId))
	}
	if !p.houseService.ExistsById(request.HouseId) {
		return response, errors.New(fmt.Sprintf("house with id %s in not exists", request.HouseId))
	}
	if request.Spec == "" {
		return response, errors.New("scheduler configuration not provided")
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

func (p *PaymentSchedulerServiceObject) Remove(id uuid.UUID) error {
	if !p.repository.ExistsById(id) {
		return errors.New(fmt.Sprintf("payment scheduler with id %s not found", id))
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
	return convert(p.repository.FindByHouseId(houseId))
}

func (p *PaymentSchedulerServiceObject) FindByUserId(userId uuid.UUID) []model.PaymentSchedulerDto {
	return convert(p.repository.FindByUserId(userId))
}

func (p *PaymentSchedulerServiceObject) schedulerFunc(payment model.PaymentScheduler) func() {
	return func() {
		if _, err := p.paymentService.Add(
			paymentModel.CreatePaymentRequest{
				Name:        payment.Name,
				Description: payment.Description,
				HouseId:     payment.HouseId,
				UserId:      payment.UserId,
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

func convert(payments []model.PaymentScheduler) []model.PaymentSchedulerDto {
	if len(payments) == 0 {
		return make([]model.PaymentSchedulerDto, 0)
	}

	var paymentsResponse []model.PaymentSchedulerDto

	for _, payment := range payments {
		paymentsResponse = append(paymentsResponse, payment.ToDto())
	}

	return paymentsResponse
}
