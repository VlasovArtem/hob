package repository

import (
	"github.com/VlasovArtem/hob/src/common"
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/db"
	groupModel "github.com/VlasovArtem/hob/src/group/model"
	"github.com/VlasovArtem/hob/src/income/model"
	"github.com/google/uuid"
	"reflect"
	"time"
)

var (
	IncomeRepositoryType = reflect.TypeOf(IncomeRepositoryObject{})
	entity               = model.Income{}
)

type IncomeRepositoryObject struct {
	db db.ModeledDatabase
}

func (i *IncomeRepositoryObject) Initialize(factory dependency.DependenciesProvider) any {
	return NewIncomeRepository(factory.FindRequiredByObject(db.DatabaseObject{}).(db.DatabaseService))
}

func (i *IncomeRepositoryObject) GetEntity() any {
	return entity
}

func NewIncomeRepository(database db.DatabaseService) IncomeRepository {
	return &IncomeRepositoryObject{
		db.ModeledDatabase{
			DatabaseService: database,
			Model:           entity,
		},
	}
}

type IncomeRepository interface {
	Create(entity model.Income) (model.Income, error)
	CreateBatch(entity []model.Income) ([]model.Income, error)
	FindById(id uuid.UUID) (model.Income, error)
	FindByHouseId(id uuid.UUID) ([]model.IncomeDto, error)
	FindByGroupIds(groupIds []uuid.UUID) ([]model.IncomeDto, error)
	ExistsById(id uuid.UUID) bool
	DeleteById(id uuid.UUID) error
	Update(id uuid.UUID, request model.UpdateIncomeRequest) error
}

func (i *IncomeRepositoryObject) Create(entity model.Income) (model.Income, error) {
	return entity, i.db.D().Omit("Groups.*").Create(&entity).Error
}

func (i *IncomeRepositoryObject) CreateBatch(entities []model.Income) ([]model.Income, error) {
	return entities, i.db.D().Omit("Groups.*").Create(&entities).Error
}

func (i *IncomeRepositoryObject) FindById(id uuid.UUID) (response model.Income, err error) {
	response.Id = id
	if err = i.db.D().Preload("Groups").First(&response).Error; err != nil {
		return model.Income{}, err
	}
	return response, err
}

func (i *IncomeRepositoryObject) FindByHouseId(id uuid.UUID) (response []model.IncomeDto, err error) {
	var responseEntities []model.Income

	if err := i.db.Modeled().Preload("Groups").Find(&responseEntities, "house_id = ?", id).Error; err != nil {
		return []model.IncomeDto{}, err
	}

	return common.MapSlice(responseEntities, func(i model.Income) model.IncomeDto {
		return i.ToDto()
	}), nil
}

func (i *IncomeRepositoryObject) FindByGroupIds(groupIds []uuid.UUID) (response []model.IncomeDto, err error) {
	var responseEntity []model.Income
	if err = i.db.D().Joins("JOIN income_groups ON income_groups.income_id = incomes.id AND income_groups.group_id IN ?", groupIds).Preload("Groups").Find(&responseEntity).Error; err != nil {
		return []model.IncomeDto{}, err
	}
	return common.MapSlice(responseEntity, func(entity model.Income) model.IncomeDto {
		return entity.ToDto()
	}), nil
}

func (i *IncomeRepositoryObject) ExistsById(id uuid.UUID) bool {
	return i.db.Exists(id)
}

func (i *IncomeRepositoryObject) DeleteById(id uuid.UUID) error {
	return i.db.Delete(id)
}

func (i *IncomeRepositoryObject) Update(id uuid.UUID, request model.UpdateIncomeRequest) error {
	err := i.db.Update(id, struct {
		Name        string
		Description string
		Date        time.Time
		Sum         float32
	}{
		request.Name,
		request.Description,
		request.Date,
		request.Sum,
	})

	if err != nil {
		return err
	}

	entity, err := i.FindById(id)

	if err != nil {
		return err
	}

	var groups = common.MapSlice(request.GroupIds, groupModel.GroupIdToGroup)

	return i.db.DM(&entity).Association("Groups").Replace(groups)
}
