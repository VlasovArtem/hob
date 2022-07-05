package repository

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/transactional"
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/income/scheduler/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IncomeRepositorySchedulerStr struct {
	db.ModeledDatabase[model.IncomeScheduler]
}

func NewIncomeSchedulerRepository(database db.DatabaseService) IncomeSchedulerRepository {
	return &IncomeRepositorySchedulerStr{db.NewModeledDatabase(model.IncomeScheduler{}, database)}
}

func (i *IncomeRepositorySchedulerStr) Initialize(factory dependency.DependenciesProvider) any {
	return NewIncomeSchedulerRepository(dependency.FindRequiredDependency[db.Database, db.DatabaseService](factory))
}

type IncomeSchedulerRepository interface {
	db.ModeledDatabase[model.IncomeScheduler]
	transactional.Transactional[IncomeSchedulerRepository]
	FindByHouseId(houseId uuid.UUID) ([]model.IncomeSchedulerDto, error)
}

func (i *IncomeRepositorySchedulerStr) FindByHouseId(houseId uuid.UUID) (response []model.IncomeSchedulerDto, err error) {
	return response, i.FindReceiverBy(&response, "house_id = ?", houseId)
}

func (i *IncomeRepositorySchedulerStr) Transactional(tx *gorm.DB) IncomeSchedulerRepository {
	return &IncomeRepositorySchedulerStr{db.NewTransactionalModeledDatabase(i.GetEntity(), tx)}
}
