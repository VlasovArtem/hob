package service

import (
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/common/dependency"
	houseService "github.com/VlasovArtem/hob/src/house/service"
	"github.com/VlasovArtem/hob/src/income/model"
	"github.com/VlasovArtem/hob/src/income/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IncomeServiceObject struct {
	houseService houseService.HouseService
	repository   repository.IncomeRepository
}

func NewIncomeService(houseService houseService.HouseService, repository repository.IncomeRepository) IncomeService {
	return &IncomeServiceObject{
		houseService: houseService,
		repository:   repository,
	}
}

func (i *IncomeServiceObject) Initialize(factory dependency.DependenciesFactory) interface{} {
	return factory.Add(
		NewIncomeService(
			factory.FindRequiredByObject(houseService.HouseServiceObject{}).(houseService.HouseService),
			factory.FindRequiredByObject(repository.IncomeRepositoryObject{}).(repository.IncomeRepository),
		),
	)
}

type IncomeService interface {
	Add(request model.CreateIncomeRequest) (model.IncomeResponse, error)
	FindById(id uuid.UUID) (model.IncomeResponse, error)
	FindByHouseId(id uuid.UUID) []model.IncomeResponse
}

func (i *IncomeServiceObject) Add(request model.CreateIncomeRequest) (response model.IncomeResponse, err error) {
	if !i.houseService.ExistsById(request.HouseId) {
		return response, errors.New(fmt.Sprintf("house with id %s not exists", request.HouseId))
	}

	if entity, err := i.repository.Create(request.ToEntity()); err != nil {
		return response, err
	} else {
		return entity.ToResponse(), nil
	}
}

func (i *IncomeServiceObject) FindById(id uuid.UUID) (response model.IncomeResponse, err error) {
	if response, err = i.repository.FindResponseById(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, errors.New(fmt.Sprintf("income with id %s not exists", id))
		}
		return response, err
	} else {
		return response, nil
	}
}

func (i *IncomeServiceObject) FindByHouseId(id uuid.UUID) []model.IncomeResponse {
	return i.repository.FindResponseByHouseId(id)
}
