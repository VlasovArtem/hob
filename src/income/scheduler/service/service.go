package service

import (
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/common/dependency"
	houseService "github.com/VlasovArtem/hob/src/house/service"
	incomeModel "github.com/VlasovArtem/hob/src/income/model"
	"github.com/VlasovArtem/hob/src/income/scheduler/model"
	"github.com/VlasovArtem/hob/src/income/scheduler/repository"
	incomeService "github.com/VlasovArtem/hob/src/income/service"
	"github.com/VlasovArtem/hob/src/scheduler"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"time"
)

type IncomeSchedulerServiceObject struct {
	houseService     houseService.HouseService
	incomeService    incomeService.IncomeService
	serviceScheduler scheduler.ServiceScheduler
	repository       repository.IncomeSchedulerRepository
}

func NewIncomeSchedulerService(
	houseService houseService.HouseService,
	incomeService incomeService.IncomeService,
	serviceScheduler scheduler.ServiceScheduler,
	repository repository.IncomeSchedulerRepository,
) IncomeSchedulerService {
	return &IncomeSchedulerServiceObject{
		houseService:     houseService,
		incomeService:    incomeService,
		serviceScheduler: serviceScheduler,
		repository:       repository,
	}
}

func (i *IncomeSchedulerServiceObject) Initialize(factory dependency.DependenciesProvider) interface{} {
	return NewIncomeSchedulerService(
		factory.FindRequiredByObject(houseService.HouseServiceObject{}).(houseService.HouseService),
		factory.FindRequiredByObject(incomeService.IncomeServiceObject{}).(incomeService.IncomeService),
		factory.FindRequiredByObject(scheduler.SchedulerServiceObject{}).(scheduler.ServiceScheduler),
		factory.FindRequiredByObject(repository.IncomeSchedulerRepositoryObject{}).(repository.IncomeSchedulerRepository),
	)
}

type IncomeSchedulerService interface {
	Add(request model.CreateIncomeSchedulerRequest) (model.IncomeSchedulerDto, error)
	Remove(id uuid.UUID) error
	FindById(id uuid.UUID) (model.IncomeSchedulerDto, error)
	FindByHouseId(id uuid.UUID) []model.IncomeSchedulerDto
}

func (i *IncomeSchedulerServiceObject) Add(request model.CreateIncomeSchedulerRequest) (response model.IncomeSchedulerDto, err error) {
	if !i.houseService.ExistsById(request.HouseId) {
		return response, errors.New(fmt.Sprintf("house with id %s not found", request.HouseId))
	} else if err != nil {
		return response, err
	}
	if request.Spec == "" {
		return response, errors.New("scheduler configuration not provided")
	}

	entity := request.ToEntity()

	if _, err = i.serviceScheduler.Add(entity.Id, string(entity.Spec), i.schedulerFunc(entity.Income)); err != nil {
		return response, err
	}

	if entity, err = i.repository.Create(entity); err != nil {
		return response, err
	} else {
		return entity.ToDto(), nil
	}
}

func (i *IncomeSchedulerServiceObject) Remove(id uuid.UUID) error {
	if !i.repository.ExistsById(id) {
		return errors.New(fmt.Sprintf("income scheduler with id %s not found", id))
	} else {
		if err := i.serviceScheduler.Remove(id); err != nil {
			log.Error().Err(err).Msg("")
		}
		i.repository.DeleteById(id)
	}
	return nil
}

func (i *IncomeSchedulerServiceObject) FindById(id uuid.UUID) (response model.IncomeSchedulerDto, err error) {
	if !i.repository.ExistsById(id) {
		return response, errors.New(fmt.Sprintf("income scheduler with id %s not found", id))
	} else {
		if paymentScheduler, err := i.repository.FindById(id); err != nil {
			return response, err
		} else {
			return paymentScheduler.ToDto(), err
		}
	}
}

func (i *IncomeSchedulerServiceObject) FindByHouseId(id uuid.UUID) []model.IncomeSchedulerDto {
	return convert(i.repository.FindByHouseId(id))
}

func (i *IncomeSchedulerServiceObject) schedulerFunc(income incomeModel.Income) func() {
	return func() {
		if _, err := i.incomeService.Add(
			incomeModel.CreateIncomeRequest{
				Name:        income.Name,
				Description: income.Description,
				Date:        time.Now(),
				Sum:         income.Sum,
				HouseId:     income.HouseId,
			},
		); err != nil {
			log.Error().Err(err).Msg("")
		} else {
			log.Info().Msgf("New income added to the house %s via scheduler %s", income.HouseId, income.Id)
		}
	}
}

func convert(payments []model.IncomeScheduler) []model.IncomeSchedulerDto {
	if len(payments) == 0 {
		return make([]model.IncomeSchedulerDto, 0)
	}

	var paymentsResponse []model.IncomeSchedulerDto

	for _, payment := range payments {
		paymentsResponse = append(paymentsResponse, payment.ToDto())
	}

	return paymentsResponse
}
