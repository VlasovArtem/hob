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
	"time"
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
	return NewIncomeService(
		factory.FindRequiredByObject(houseService.HouseServiceObject{}).(houseService.HouseService),
		factory.FindRequiredByObject(repository.IncomeRepositoryObject{}).(repository.IncomeRepository),
	)
}

type IncomeService interface {
	Add(request model.CreateIncomeRequest) (model.IncomeDto, error)
	FindById(id uuid.UUID) (model.IncomeDto, error)
	FindByHouseId(id uuid.UUID) []model.IncomeDto
	ExistsById(id uuid.UUID) bool
	DeleteById(id uuid.UUID) error
	Update(request model.UpdateIncomeRequest) error
}

func (i *IncomeServiceObject) Add(request model.CreateIncomeRequest) (response model.IncomeDto, err error) {
	if !i.houseService.ExistsById(request.HouseId) {
		return response, errors.New(fmt.Sprintf("house with id %s not exists", request.HouseId))
	}
	if request.Date.After(time.Now()) {
		return response, errors.New("date should not be after current date")
	}

	if entity, err := i.repository.Create(request.CreateToEntity()); err != nil {
		return response, err
	} else {
		return entity.ToDto(), nil
	}
}

func (i *IncomeServiceObject) FindById(id uuid.UUID) (response model.IncomeDto, err error) {
	if response, err = i.repository.FindDtoById(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, errors.New(fmt.Sprintf("income with id %s not exists", id))
		}
		return response, err
	} else {
		return response, nil
	}
}

func (i *IncomeServiceObject) FindByHouseId(id uuid.UUID) []model.IncomeDto {
	return i.repository.FindResponseByHouseId(id)
}

func (i *IncomeServiceObject) ExistsById(id uuid.UUID) bool {
	return i.repository.ExistsById(id)
}

func (i *IncomeServiceObject) DeleteById(id uuid.UUID) error {
	if !i.ExistsById(id) {
		return errors.New(fmt.Sprintf("income with id %s not found", id))
	}
	return i.repository.DeleteById(id)
}

func (i *IncomeServiceObject) Update(request model.UpdateIncomeRequest) error {
	if !i.ExistsById(request.Id) {
		return errors.New(fmt.Sprintf("income with id %s not found", request.Id))
	}
	if request.Date.After(time.Now()) {
		return errors.New("date should not be after current date")
	}
	return i.repository.Update(request.UpdateToEntity())
}
