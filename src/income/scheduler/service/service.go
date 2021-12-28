package service

import (
	"errors"
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/int-errors"
	houseService "github.com/VlasovArtem/hob/src/house/service"
	incomeModel "github.com/VlasovArtem/hob/src/income/model"
	"github.com/VlasovArtem/hob/src/income/scheduler/model"
	"github.com/VlasovArtem/hob/src/income/scheduler/repository"
	incomeService "github.com/VlasovArtem/hob/src/income/service"
	"github.com/VlasovArtem/hob/src/scheduler"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"reflect"
	"time"
)

var IncomeSchedulerServiceType = reflect.TypeOf(IncomeSchedulerServiceObject{})

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

func (i *IncomeSchedulerServiceObject) Initialize(factory dependency.DependenciesProvider) any {
	return NewIncomeSchedulerService(
		factory.FindRequiredByType(houseService.HouseServiceType).(houseService.HouseService),
		factory.FindRequiredByType(incomeService.IncomeServiceType).(incomeService.IncomeService),
		factory.FindRequiredByType(scheduler.SchedulerServiceType).(scheduler.ServiceScheduler),
		factory.FindRequiredByType(repository.IncomeSchedulerRepositoryType).(repository.IncomeSchedulerRepository),
	)
}

type IncomeSchedulerService interface {
	Add(request model.CreateIncomeSchedulerRequest) (model.IncomeSchedulerDto, error)
	DeleteById(id uuid.UUID) error
	Update(id uuid.UUID, request model.UpdateIncomeSchedulerRequest) error
	FindById(id uuid.UUID) (model.IncomeSchedulerDto, error)
	FindByHouseId(id uuid.UUID) []model.IncomeSchedulerDto
}

func (i *IncomeSchedulerServiceObject) Add(request model.CreateIncomeSchedulerRequest) (response model.IncomeSchedulerDto, err error) {
	if err := i.validateCreateRequest(request); err != nil {
		return response, err
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

func (i *IncomeSchedulerServiceObject) Update(id uuid.UUID, request model.UpdateIncomeSchedulerRequest) error {
	if err := i.validateUpdateRequest(id, request); err != nil {
		return err
	}

	updatedEntity, err := i.repository.Update(id, request)
	if err != nil {
		return err
	}

	if _, err := i.serviceScheduler.Update(id, string(updatedEntity.Spec), i.schedulerFunc(updatedEntity.Income)); err != nil {
		if err := i.repository.DeleteById(id); err != nil {
			log.Err(err)
		}
		return err
	}

	return nil
}

func (i *IncomeSchedulerServiceObject) DeleteById(id uuid.UUID) error {
	if !i.repository.ExistsById(id) {
		return int_errors.NewErrNotFound("income scheduler with id %s not found", id)
	} else {
		if err := i.serviceScheduler.Remove(id); err != nil {
			log.Error().Err(err)
		}
		return i.repository.DeleteById(id)
	}
}

func (i *IncomeSchedulerServiceObject) FindById(id uuid.UUID) (response model.IncomeSchedulerDto, err error) {
	if !i.repository.ExistsById(id) {
		return response, int_errors.NewErrNotFound("income scheduler with id %s not found", id)
	} else {
		if paymentScheduler, err := i.repository.FindById(id); err != nil {
			return response, err
		} else {
			return paymentScheduler.ToDto(), err
		}
	}
}

func (i *IncomeSchedulerServiceObject) FindByHouseId(id uuid.UUID) []model.IncomeSchedulerDto {
	responses, err := i.repository.FindByHouseId(id)
	if err != nil {
		log.Err(err)
	}
	return responses
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

func (i *IncomeSchedulerServiceObject) validateCreateRequest(request model.CreateIncomeSchedulerRequest) error {
	if request.Sum <= 0 {
		return errors.New("sum should not be zero of negative")
	}
	if !i.houseService.ExistsById(request.HouseId) {
		return int_errors.NewErrNotFound("house with id %s not found", request.HouseId)
	}

	if request.Spec == "" {
		return errors.New("scheduler configuration not provided")
	}

	return nil
}

func (i *IncomeSchedulerServiceObject) validateUpdateRequest(id uuid.UUID, request model.UpdateIncomeSchedulerRequest) error {
	if request.Sum <= 0 {
		return errors.New("sum should not be zero of negative")
	}
	if !i.repository.ExistsById(id) {
		return int_errors.NewErrNotFound("income scheduler with id %s not found", id)
	}

	if request.Spec == "" {
		return errors.New("scheduler configuration not provided")
	}

	return nil
}
