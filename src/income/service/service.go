package service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	hs "house/service"
	"income/model"
)

type incomeServiceObject struct {
	houseService hs.HouseService
	incomes      map[uuid.UUID]model.Income
	houseIncomes map[uuid.UUID]model.Income
}

func NewIncomeService(houseService hs.HouseService) IncomeService {
	return &incomeServiceObject{
		houseService: houseService,
		incomes:      make(map[uuid.UUID]model.Income),
		houseIncomes: make(map[uuid.UUID]model.Income),
	}
}

type IncomeService interface {
	Add(request model.CreateIncomeRequest) (model.IncomeResponse, error)
	FindById(id uuid.UUID) (model.IncomeResponse, error)
	FindByHouseId(id uuid.UUID) (model.IncomeResponse, error)
}

func (i *incomeServiceObject) Add(request model.CreateIncomeRequest) (response model.IncomeResponse, err error) {
	if !i.houseService.ExistsById(request.HouseId) {
		return response, errors.New(fmt.Sprintf("house with id %s not exists", request.HouseId))
	}

	entity := request.ToEntity()

	i.incomes[entity.Id] = entity
	i.houseIncomes[entity.HouseId] = entity

	return entity.ToResponse(), nil
}

func (i *incomeServiceObject) FindById(id uuid.UUID) (response model.IncomeResponse, err error) {
	if income, ok := i.incomes[id]; !ok {
		return response, errors.New(fmt.Sprintf("income with id %s not exists", id))
	} else {
		return income.ToResponse(), nil
	}
}

func (i *incomeServiceObject) FindByHouseId(id uuid.UUID) (response model.IncomeResponse, err error) {
	if income, ok := i.houseIncomes[id]; !ok {
		return response, errors.New(fmt.Sprintf("income with house id %s not found", id))

	} else {
		return income.ToResponse(), nil
	}
}
