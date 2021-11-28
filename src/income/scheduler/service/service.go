package service

import (
	"common/dependency"
	"errors"
	"fmt"
	"github.com/google/uuid"
	hs "house/service"
	im "income/model"
	"income/scheduler/model"
	"income/scheduler/repository"
	is "income/service"
	"log"
	"scheduler"
	"time"
)

type IncomeSchedulerServiceObject struct {
	houseService     hs.HouseService
	incomeService    is.IncomeService
	serviceScheduler scheduler.ServiceScheduler
	repository       repository.IncomeSchedulerRepository
}

func NewIncomeSchedulerService(
	houseService hs.HouseService,
	incomeService is.IncomeService,
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

func (i *IncomeSchedulerServiceObject) Initialize(factory dependency.DependenciesFactory) {
	factory.Add(
		NewIncomeSchedulerService(
			factory.FindRequiredByObject(hs.HouseServiceObject{}).(hs.HouseService),
			factory.FindRequiredByObject(is.IncomeServiceObject{}).(is.IncomeService),
			factory.FindRequiredByObject(scheduler.SchedulerServiceObject{}).(scheduler.ServiceScheduler),
			factory.FindRequiredByObject(repository.IncomeSchedulerRepositoryObject{}).(repository.IncomeSchedulerRepository),
		),
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
			log.Println(err)
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

func (i *IncomeSchedulerServiceObject) schedulerFunc(income im.Income) func() {
	return func() {
		if _, err := i.incomeService.Add(
			im.CreateIncomeRequest{
				Name:        income.Name,
				Description: income.Description,
				Date:        time.Now(),
				Sum:         income.Sum,
				HouseId:     income.HouseId,
			},
		); err != nil {
			log.Println(err)
		} else {
			log.Println(fmt.Sprintf("New income added to the house %s via scheduler %s", income.HouseId, income.Id))
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
