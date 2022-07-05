package service

import (
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/common"
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/int-errors"
	"github.com/VlasovArtem/hob/src/common/transactional"
	groupModel "github.com/VlasovArtem/hob/src/group/model"
	groupService "github.com/VlasovArtem/hob/src/group/service"
	houseService "github.com/VlasovArtem/hob/src/house/service"
	"github.com/VlasovArtem/hob/src/income/model"
	"github.com/VlasovArtem/hob/src/income/repository"
	pivotalService "github.com/VlasovArtem/hob/src/pivotal/service"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"time"
)

type IncomeServiceStr struct {
	houseService   houseService.HouseService
	groupService   groupService.GroupService
	repository     repository.IncomeRepository
	pivotalService pivotalService.PivotalService
}

func (i *IncomeServiceStr) Transactional(tx *gorm.DB) IncomeService {
	return &IncomeServiceStr{
		houseService:   i.houseService.Transactional(tx),
		groupService:   i.groupService.Transactional(tx),
		repository:     i.repository.Transactional(tx),
		pivotalService: i.pivotalService.Transactional(tx),
	}
}

func NewIncomeService(
	houseService houseService.HouseService,
	groupService groupService.GroupService,
	repository repository.IncomeRepository,
	pivotalService pivotalService.PivotalService,
) IncomeService {
	return &IncomeServiceStr{
		houseService:   houseService,
		groupService:   groupService,
		repository:     repository,
		pivotalService: pivotalService,
	}
}

func (i *IncomeServiceStr) Initialize(factory dependency.DependenciesProvider) any {
	return NewIncomeService(
		dependency.FindRequiredDependency[houseService.HouseServiceStr, houseService.HouseService](factory),
		dependency.FindRequiredDependency[groupService.GroupServiceStr, groupService.GroupService](factory),
		dependency.FindRequiredDependency[repository.IncomeRepositoryStr, repository.IncomeRepository](factory),
		dependency.FindRequiredDependency[pivotalService.PivotalServiceStr, pivotalService.PivotalService](factory),
	)
}

type IncomeService interface {
	transactional.Transactional[IncomeService]
	Add(request model.CreateIncomeRequest) (model.IncomeDto, error)
	AddBatch(request model.CreateIncomeBatchRequest) ([]model.IncomeDto, error)
	FindById(id uuid.UUID) (model.IncomeDto, error)
	FindByHouseId(id uuid.UUID, limit int, offset int, from, to *time.Time) []model.IncomeDto
	FindByGroupIds(ids []uuid.UUID, limit int, offset int, from, to *time.Time) []model.IncomeDto
	ExistsById(id uuid.UUID) bool
	DeleteById(id uuid.UUID) error
	Update(id uuid.UUID, request model.UpdateIncomeRequest) error
}

func (i *IncomeServiceStr) Add(request model.CreateIncomeRequest) (response model.IncomeDto, err error) {
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

	var income model.Income

	err = i.repository.DB().Transaction(func(tx *gorm.DB) error {
		trx := i.Transactional(tx).(*IncomeServiceStr)

		income = request.ToEntity()

		if err = trx.repository.Create(&income, "Groups.*"); err != nil {
			return err
		}

		if trx.pivotalService.ExistsByHouseId(*request.HouseId) {
			return trx.pivotalService.AddIncome(float64(request.Sum), request.Date.Add(1*time.Microsecond), request.HouseId, request.GroupIds)
		}

		return nil
	})

	return income.ToDto(), err
}

func (i *IncomeServiceStr) AddBatch(request model.CreateIncomeBatchRequest) (response []model.IncomeDto, err error) {
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

	return common.MapSlice(entities, model.IncomeToDto), i.repository.DB().Transaction(func(tx *gorm.DB) error {
		trx := i.Transactional(tx).(*IncomeServiceStr)

		if err = i.repository.Create(&entities, "Groups.*"); err != nil {
			return err
		}

		for _, payment := range entities {
			if trx.pivotalService.ExistsByHouseId(*payment.HouseId) {
				return trx.pivotalService.AddIncome(float64(payment.Sum), payment.Date.Add(1*time.Microsecond), payment.HouseId, common.MapSlice[groupModel.Group, uuid.UUID](payment.Groups, func(g groupModel.Group) uuid.UUID { return g.Id }))
			}
		}

		return nil
	})
}

func (i *IncomeServiceStr) FindById(id uuid.UUID) (response model.IncomeDto, err error) {
	if entity, err := i.repository.FindById(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, int_errors.NewErrNotFound("income with id %s not found", id)
		}
		return response, err
	} else {
		return entity.ToDto(), nil
	}
}

func (i *IncomeServiceStr) FindByHouseId(id uuid.UUID, limit int, offset int, from, to *time.Time) []model.IncomeDto {
	response, err := i.repository.FindByHouseId(id, limit, offset, from, to)

	if err != nil {
		log.Err(err).Msg("failed to find incomes by house id")
	}

	return response
}

func (i *IncomeServiceStr) FindByGroupIds(ids []uuid.UUID, limit int, offset int, from, to *time.Time) []model.IncomeDto {
	response, err := i.repository.FindByGroupIds(ids, limit, offset, from, to)

	if err != nil {
		log.Err(err).Msg("failed to find incomes by group ids")
	}

	return response
}

func (i *IncomeServiceStr) ExistsById(id uuid.UUID) bool {
	return i.repository.Exists(id)
}

func (i *IncomeServiceStr) DeleteById(id uuid.UUID) error {
	if !i.ExistsById(id) {
		return int_errors.NewErrNotFound("income with id %s not found", id)
	}
	if err := i.repository.Delete(id); err != nil {
		return err
	}

	return nil
}

func (i *IncomeServiceStr) Update(id uuid.UUID, request model.UpdateIncomeRequest) error {
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
