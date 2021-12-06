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
	return NewIncomeRepository(factory.FindRequiredByObject(db.DatabaseObject{}).(db.DatabaseService))
}

func (i *IncomeRepositoryObject) GetEntity() interface{} {
	return model.Income{}
}

func NewIncomeRepository(database db.DatabaseService) IncomeRepository {
	return &IncomeRepositoryObject{database}
}

type IncomeRepository interface {
	Create(entity model.Income) (model.Income, error)
	FindDtoById(id uuid.UUID) (model.IncomeDto, error)
	FindResponseByHouseId(id uuid.UUID) []model.IncomeDto
	ExistsById(id uuid.UUID) bool
	DeleteById(id uuid.UUID) error
	Update(entity model.Income) error
}

func (i *IncomeRepositoryObject) Create(entity model.Income) (model.Income, error) {
	return entity, i.database.Create(&entity)
}

func (i *IncomeRepositoryObject) FindDtoById(id uuid.UUID) (response model.IncomeDto, err error) {
	return response, i.database.FindByIdModeled(model.Income{}, &response, id)
}

func (i *IncomeRepositoryObject) FindResponseByHouseId(id uuid.UUID) (response []model.IncomeDto) {
	i.database.DM(model.Income{}).Find(&response, "house_id = ?", id)

	return response
}

func (i *IncomeRepositoryObject) ExistsById(id uuid.UUID) bool {
	return i.database.ExistsById(model.Income{}, id)
}

func (i *IncomeRepositoryObject) DeleteById(id uuid.UUID) error {
	return i.database.D().Delete(model.Income{}, id).Error
}

func (i *IncomeRepositoryObject) Update(entity model.Income) error {
	return i.database.UpdateById(model.Income{}, entity.Id, entity, "HouseId", "House")
}
