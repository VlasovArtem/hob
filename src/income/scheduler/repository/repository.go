package repository

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/income/scheduler/model"
	"github.com/google/uuid"
)

var entity = model.IncomeScheduler{}

type IncomeSchedulerRepositoryObject struct {
	database db.ModeledDatabase
}

func NewIncomeSchedulerRepository(database db.DatabaseService) IncomeSchedulerRepository {
	return &IncomeSchedulerRepositoryObject{
		db.ModeledDatabase{
			DatabaseService: database,
			Model:           entity,
		},
	}
}

func (i *IncomeSchedulerRepositoryObject) Initialize(factory dependency.DependenciesProvider) any {
	return NewIncomeSchedulerRepository(dependency.FindRequiredDependency[db.DatabaseObject, db.DatabaseService](factory))
}

func (i *IncomeSchedulerRepositoryObject) GetEntity() any {
	return entity
}

type IncomeSchedulerRepository interface {
	Create(scheduler model.IncomeScheduler) (model.IncomeScheduler, error)
	ExistsById(id uuid.UUID) bool
	DeleteById(id uuid.UUID) error
	FindById(id uuid.UUID) (model.IncomeScheduler, error)
	FindByHouseId(houseId uuid.UUID) ([]model.IncomeSchedulerDto, error)
	Update(id uuid.UUID, scheduler model.UpdateIncomeSchedulerRequest) (model.IncomeScheduler, error)
}

func (i *IncomeSchedulerRepositoryObject) Create(scheduler model.IncomeScheduler) (model.IncomeScheduler, error) {
	return scheduler, i.database.Create(&scheduler)
}

func (i *IncomeSchedulerRepositoryObject) ExistsById(id uuid.UUID) bool {
	return i.database.Exists(id)
}

func (i *IncomeSchedulerRepositoryObject) DeleteById(id uuid.UUID) error {
	return i.database.Delete(id)
}

func (i *IncomeSchedulerRepositoryObject) FindById(id uuid.UUID) (response model.IncomeScheduler, err error) {
	return response, i.database.Find(&response, id)
}

func (i *IncomeSchedulerRepositoryObject) FindByHouseId(houseId uuid.UUID) (response []model.IncomeSchedulerDto, err error) {
	return response, i.database.FindBy(&response, "house_id = ?", houseId)
}

func (i *IncomeSchedulerRepositoryObject) Update(id uuid.UUID, scheduler model.UpdateIncomeSchedulerRequest) (response model.IncomeScheduler, err error) {
	if err = i.database.Update(id, scheduler); err != nil {
		return response, err
	}

	return i.FindById(id)
}
