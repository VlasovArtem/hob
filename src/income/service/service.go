package service

import (
	"errors"
	"github.com/VlasovArtem/hob/src/common"
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/int-errors"
	groupService "github.com/VlasovArtem/hob/src/group/service"
	houseService "github.com/VlasovArtem/hob/src/house/service"
	"github.com/VlasovArtem/hob/src/income/model"
	"github.com/VlasovArtem/hob/src/income/repository"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"reflect"
	"time"
)

var IncomeServiceType = reflect.TypeOf(IncomeServiceObject{})

type IncomeServiceObject struct {
	houseService houseService.HouseService
	groupService groupService.GroupService
	repository   repository.IncomeRepository
}

func NewIncomeService(
	houseService houseService.HouseService,
	groupService groupService.GroupService,
	repository repository.IncomeRepository,
) IncomeService {
	return &IncomeServiceObject{
		houseService: houseService,
		groupService: groupService,
		repository:   repository,
	}
}

func (i *IncomeServiceObject) Initialize(factory dependency.DependenciesProvider) any {
	return NewIncomeService(
		dependency.FindRequiredDependency[houseService.HouseServiceObject, houseService.HouseService](factory),
		dependency.FindRequiredDependency[groupService.GroupServiceObject, groupService.GroupService](factory),
		dependency.FindRequiredDependency[repository.IncomeRepositoryObject, repository.IncomeRepository](factory),
	)
}

type IncomeService interface {
	Add(request model.CreateIncomeRequest) (model.IncomeDto, error)
	FindById(id uuid.UUID) (model.IncomeDto, error)
	FindByHouseId(id uuid.UUID) []model.IncomeDto
	FindByGroupIds(ids []uuid.UUID) []model.IncomeDto
	ExistsById(id uuid.UUID) bool
	DeleteById(id uuid.UUID) error
	Update(id uuid.UUID, request model.UpdateIncomeRequest) error
}

func (i *IncomeServiceObject) Add(request model.CreateIncomeRequest) (response model.IncomeDto, err error) {
	if !i.houseService.ExistsById(request.HouseId) {
		return response, int_errors.NewErrNotFound("house with id %s not exists", request.HouseId)
	}
	if len(request.GroupIds) != 0 && !i.groupService.ExistsByIds(request.GroupIds) {
		return response, int_errors.NewErrNotFound("not all group with ids %s found", common.Join(request.GroupIds, ","))
	}
	if request.Date.After(time.Now()) {
		return response, errors.New("date should not be after current date")
	}

	if entity, err := i.repository.Create(request.ToEntity()); err != nil {
		return response, err
	} else {
		return entity.ToDto(), nil
	}
}

func (i *IncomeServiceObject) FindById(id uuid.UUID) (response model.IncomeDto, err error) {
	if entity, err := i.repository.FindById(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, int_errors.NewErrNotFound("income with id %s not found", id)
		}
		return response, err
	} else {
		return entity.ToDto(), nil
	}
}

func (i *IncomeServiceObject) FindByHouseId(id uuid.UUID) []model.IncomeDto {
	response, err := i.repository.FindByHouseId(id)

	if err != nil {
		log.Err(err)
	}

	return response
}

func (i *IncomeServiceObject) FindByGroupIds(ids []uuid.UUID) []model.IncomeDto {
	response, err := i.repository.FindByGroupIds(ids)

	if err != nil {
		log.Err(err)
	}

	return response
}

func (i *IncomeServiceObject) ExistsById(id uuid.UUID) bool {
	return i.repository.ExistsById(id)
}

func (i *IncomeServiceObject) DeleteById(id uuid.UUID) error {
	if !i.ExistsById(id) {
		return int_errors.NewErrNotFound("income with id %s not found", id)
	}
	return i.repository.DeleteById(id)
}

func (i *IncomeServiceObject) Update(id uuid.UUID, request model.UpdateIncomeRequest) error {
	if !i.ExistsById(id) {
		return int_errors.NewErrNotFound("income with id %s not found", id)
	}
	if len(request.GroupIds) != 0 && !i.groupService.ExistsByIds(request.GroupIds) {
		return int_errors.NewErrNotFound("not all group with ids %s found", common.Join(request.GroupIds, ","))
	}
	if request.Date.After(time.Now()) {
		return errors.New("date should not be after current date")
	}
	return i.repository.Update(id, request)
}
