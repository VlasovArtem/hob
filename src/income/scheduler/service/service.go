package service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	hs "house/service"
	im "income/model"
	"income/scheduler/model"
	is "income/service"
	"log"
	"scheduler"
	"time"
)

type incomeSchedulerServiceObject struct {
	houseService          hs.HouseService
	incomeService         is.IncomeService
	serviceScheduler      scheduler.ServiceScheduler
	incomeSchedulers      map[uuid.UUID]model.IncomeScheduler
	houseIncomeSchedulers map[uuid.UUID]model.IncomeScheduler
}

func NewIncomeSchedulerService(
	houseService hs.HouseService,
	incomeService is.IncomeService,
	serviceScheduler scheduler.ServiceScheduler) IncomeSchedulerService {
	return &incomeSchedulerServiceObject{
		houseService:          houseService,
		incomeService:         incomeService,
		serviceScheduler:      serviceScheduler,
		incomeSchedulers:      make(map[uuid.UUID]model.IncomeScheduler),
		houseIncomeSchedulers: make(map[uuid.UUID]model.IncomeScheduler),
	}
}

type IncomeSchedulerService interface {
	Add(request model.CreateIncomeSchedulerRequest) (model.IncomeSchedulerResponse, error)
	Remove(id uuid.UUID) error
	FindById(id uuid.UUID) (model.IncomeSchedulerResponse, error)
	FindByHouseId(id uuid.UUID) (model.IncomeSchedulerResponse, error)
}

func (i *incomeSchedulerServiceObject) Add(request model.CreateIncomeSchedulerRequest) (response model.IncomeSchedulerResponse, err error) {
	if !i.houseService.ExistsById(request.HouseId) {
		return response, errors.New(fmt.Sprintf("house with id %s not found", request.HouseId))
	}
	if request.Spec == "" {
		return response, errors.New("scheduler configuration not provided")
	}

	entity := request.ToEntity()

	if _, err = i.serviceScheduler.Add(entity.Id, string(entity.Spec), i.schedulerFunc(entity.Income)); err != nil {
		return response, err
	}

	i.incomeSchedulers[entity.Id] = entity
	i.houseIncomeSchedulers[entity.HouseId] = entity

	return entity.ToResponse(), nil
}

func (i *incomeSchedulerServiceObject) Remove(id uuid.UUID) error {
	if incomeScheduler, ok := i.incomeSchedulers[id]; !ok {
		return errors.New(fmt.Sprintf("income scheduler with id %s not found", id))
	} else {
		if err := i.serviceScheduler.Remove(id); err != nil {
			log.Println(err)
		}
		delete(i.incomeSchedulers, id)
		delete(i.houseIncomeSchedulers, incomeScheduler.HouseId)
	}
	return nil
}

func (i *incomeSchedulerServiceObject) FindById(id uuid.UUID) (response model.IncomeSchedulerResponse, err error) {
	if income, ok := i.incomeSchedulers[id]; !ok {
		return response, errors.New(fmt.Sprintf("income scheduler with id %s not found", id))
	} else {
		return income.ToResponse(), nil
	}
}

func (i *incomeSchedulerServiceObject) FindByHouseId(id uuid.UUID) (response model.IncomeSchedulerResponse, err error) {
	if income, ok := i.houseIncomeSchedulers[id]; !ok {
		return response, errors.New(fmt.Sprintf("income scheduler with house id %s not found", id))

	} else {
		return income.ToResponse(), nil
	}
}

func (i *incomeSchedulerServiceObject) schedulerFunc(income im.Income) func() {
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
