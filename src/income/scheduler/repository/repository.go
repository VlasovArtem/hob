package repository

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/income/scheduler/model"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type IncomeSchedulerRepositoryObject struct {
	database db.DatabaseService
}

func NewIncomeSchedulerRepository(database db.DatabaseService) IncomeSchedulerRepository {
	return &IncomeSchedulerRepositoryObject{database}
}

func (i *IncomeSchedulerRepositoryObject) Initialize(factory dependency.DependenciesFactory) interface{} {
	return factory.Add(NewIncomeSchedulerRepository(factory.FindRequiredByObject(db.DatabaseObject{}).(db.DatabaseService)))
}

func (i *IncomeSchedulerRepositoryObject) GetEntity() interface{} {
	return model.IncomeScheduler{}
}

type IncomeSchedulerRepository interface {
	Create(scheduler model.IncomeScheduler) (model.IncomeScheduler, error)
	ExistsById(id uuid.UUID) bool
	DeleteById(id uuid.UUID)
	FindById(id uuid.UUID) (model.IncomeScheduler, error)
	FindByHouseId(houseId uuid.UUID) []model.IncomeScheduler
}

func (i *IncomeSchedulerRepositoryObject) Create(scheduler model.IncomeScheduler) (model.IncomeScheduler, error) {
	return scheduler, i.database.Create(&scheduler)
}

func (i *IncomeSchedulerRepositoryObject) ExistsById(id uuid.UUID) bool {
	return i.database.ExistsById(model.IncomeScheduler{}, id)
}

func (i *IncomeSchedulerRepositoryObject) DeleteById(id uuid.UUID) {
	i.database.D().Delete(model.IncomeScheduler{}, id)
}

func (i *IncomeSchedulerRepositoryObject) FindById(id uuid.UUID) (response model.IncomeScheduler, err error) {
	return response, i.database.FindById(&response, id)
}

func (i *IncomeSchedulerRepositoryObject) FindByHouseId(houseId uuid.UUID) (response []model.IncomeScheduler) {
	if tx := i.database.D().Find(&response, "house_id = ?", houseId); tx.Error != nil {
		log.Error().Err(tx.Error).Msg("")
	}
	return response
}
