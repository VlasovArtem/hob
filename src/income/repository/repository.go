package repository

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/income/model"
	"github.com/google/uuid"
)

type IncomeRepositoryObject struct {
	database db.DatabaseService
}

func (i *IncomeRepositoryObject) Initialize(factory dependency.DependenciesFactory) interface{} {
	return factory.Add(NewIncomeRepository(factory.FindRequiredByObject(db.DatabaseObject{}).(db.DatabaseService)))
}

func (i *IncomeRepositoryObject) GetEntity() interface{} {
	return model.Income{}
}

func NewIncomeRepository(database db.DatabaseService) IncomeRepository {
	return &IncomeRepositoryObject{database}
}

type IncomeRepository interface {
	Create(entity model.Income) (model.Income, error)
	FindResponseById(id uuid.UUID) (model.IncomeResponse, error)
	FindResponseByHouseId(id uuid.UUID) []model.IncomeResponse
}

func (i *IncomeRepositoryObject) Create(entity model.Income) (model.Income, error) {
	return entity, i.database.Create(&entity)
}

func (i *IncomeRepositoryObject) FindResponseById(id uuid.UUID) (response model.IncomeResponse, err error) {
	return response, i.database.FindByIdModeled(model.Income{}, &response, id)
}

func (i *IncomeRepositoryObject) FindResponseByHouseId(id uuid.UUID) (response []model.IncomeResponse) {
	i.database.DM(model.Income{}).Find(&response, "house_id = ?", id)

	return response
}
