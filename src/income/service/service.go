package service

import (
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/common"
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/int-errors"
	groupService "github.com/VlasovArtem/hob/src/group/service"
	houseService "github.com/VlasovArtem/hob/src/house/service"
	"github.com/VlasovArtem/hob/src/income/model"
	"github.com/VlasovArtem/hob/src/income/repository"
	"github.com/VlasovArtem/hob/src/pivotal/cache"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"time"
)

type IncomeServiceObject struct {
	houseService houseService.HouseService
	groupService groupService.GroupService
	repository   repository.IncomeRepository
	pivotalCache cache.PivotalCache
}

func NewIncomeService(
	houseService houseService.HouseService,
	groupService groupService.GroupService,
	repository repository.IncomeRepository,
	pivotalCache cache.PivotalCache,
) IncomeService {
	return &IncomeServiceObject{
		houseService: houseService,
		groupService: groupService,
		repository:   repository,
		pivotalCache: pivotalCache,
	}
}

func (i *IncomeServiceObject) Initialize(factory dependency.DependenciesProvider) any {
	return NewIncomeService(
		dependency.FindRequiredDependency[houseService.HouseServiceObject, houseService.HouseService](factory),
		dependency.FindRequiredDependency[groupService.GroupServiceObject, groupService.GroupService](factory),
		dependency.FindRequiredDependency[repository.IncomeRepositoryObject, repository.IncomeRepository](factory),
		dependency.FindRequiredDependency[cache.PivotalCacheObject, cache.PivotalCache](factory),
	)
}

type IncomeService interface {
	Add(request model.CreateIncomeRequest) (model.IncomeDto, error)
	AddBatch(request model.CreateIncomeBatchRequest) ([]model.IncomeDto, error)
	FindById(id uuid.UUID) (model.IncomeDto, error)
	FindByHouseId(id uuid.UUID, limit int, offset int, from, to *time.Time) []model.IncomeDto
	FindByGroupIds(ids []uuid.UUID, limit int, offset int, from, to *time.Time) []model.IncomeDto
	ExistsById(id uuid.UUID) bool
	DeleteById(id uuid.UUID) error
	Update(id uuid.UUID, request model.UpdateIncomeRequest) error
	CalculateSumByHouseId(houseId uuid.UUID, from *time.Time) float64
	CalculateSumByGroupId(groupId uuid.UUID, from *time.Time) float64
}

func (i *IncomeServiceObject) Add(request model.CreateIncomeRequest) (response model.IncomeDto, err error) {
	if request.HouseId == nil && len(request.GroupIds) == 0 {
		return response, errors.New("houseId or groupId must be set")
	}
	if request.HouseId != nil && !i.houseService.ExistsById(*request.HouseId) {
		return response, int_errors.NewErrNotFound("house with id %s not found", request.HouseId)
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
		dto := entity.ToDto()
		i.invalidateCache(dto)
		return dto, nil
	}
}

func (i *IncomeServiceObject) AddBatch(request model.CreateIncomeBatchRequest) (response []model.IncomeDto, err error) {
	if len(request.Incomes) == 0 {
		return make([]model.IncomeDto, 0), nil
	}

	houseIds := make(map[uuid.UUID]bool)
	groups := make(map[uuid.UUID]bool)

	err = common.ForEach(request.Incomes, func(income model.CreateIncomeRequest) error {
		if income.HouseId == nil && len(income.GroupIds) == 0 {
			return errors.New("houseId or groupId must be set")
		}

		if income.HouseId != nil {
			houseIds[*income.HouseId] = true
		}
		_ = common.ForEach(income.GroupIds, func(groupId uuid.UUID) error {
			groups[groupId] = true
			return nil
		})

		return nil
	})

	if err != nil {
		return make([]model.IncomeDto, 0), err
	}

	builder := int_errors.NewBuilder()

	for houseId := range houseIds {
		if !i.houseService.ExistsById(houseId) {
			builder.WithDetail(fmt.Sprintf("house with id %s not found", houseId))
		}
	}

	var groupIds []uuid.UUID

	for groupId := range groups {
		groupIds = append(groupIds, groupId)
	}

	if len(groupIds) != 0 && !i.groupService.ExistsByIds(groupIds) {
		builder.WithDetail(fmt.Sprintf("not all group with ids %s found", common.Join(groupIds, ",")))
	}

	for _, income := range request.Incomes {
		if income.Date.After(time.Now()) {
			builder.WithDetail("date should not be after current date")
		}
	}

	if builder.HasErrors() {
		return nil, int_errors.NewErrResponse(builder.WithMessage("Create income batch failed"))
	}

	entities := common.MapSlice(request.Incomes, func(income model.CreateIncomeRequest) model.Income {
		return income.ToEntity()
	})

	if repositoryResponse, err := i.repository.CreateBatch(entities); err != nil {
		return nil, err
	} else {
		incomes := common.MapSlice(repositoryResponse, model.IncomeToDto)
		for _, dto := range incomes {
			i.invalidateCache(dto)
		}

		return incomes, nil
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

func (i *IncomeServiceObject) FindByHouseId(id uuid.UUID, limit int, offset int, from, to *time.Time) []model.IncomeDto {
	response, err := i.repository.FindByHouseId(id, limit, offset, from, to)

	if err != nil {
		log.Err(err).Msg("failed to find incomes by house id")
	}

	return response
}

func (i *IncomeServiceObject) FindByGroupIds(ids []uuid.UUID, limit int, offset int, from, to *time.Time) []model.IncomeDto {
	response, err := i.repository.FindByGroupIds(ids, limit, offset, from, to)

	if err != nil {
		log.Err(err).Msg("failed to find incomes by group ids")
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
	income, _ := i.FindById(id)
	if err := i.repository.DeleteById(id); err != nil {
		return err
	}

	i.invalidateCache(income)
	return nil
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

func (i *IncomeServiceObject) CalculateSumByHouseId(houseId uuid.UUID, from *time.Time) (sum float64) {
	i.repository.CalculateSumByHouseId(houseId, from, &sum)
	return
}

func (i *IncomeServiceObject) CalculateSumByGroupId(groupId uuid.UUID, from *time.Time) (sum float64) {
	i.repository.CalculateSumByGroupId(groupId, from, &sum)
	return
}

func (i *IncomeServiceObject) invalidateCache(income model.IncomeDto) {
	if income.HouseId != nil {
		i.pivotalCache.Invalidate(*income.HouseId)
	}
	for _, group := range income.Groups {
		i.pivotalCache.Invalidate(group.Id)
	}
}
