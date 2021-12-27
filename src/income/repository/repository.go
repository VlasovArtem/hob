package repository

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/income/model"
	"github.com/google/uuid"
	"reflect"
)

var (
	IncomeRepositoryType = reflect.TypeOf(IncomeRepositoryObject{})
	entity               = model.Income{}
)

type IncomeRepositoryObject struct {
	database db.ModeledDatabase
}

func (i *IncomeRepositoryObject) Initialize(factory dependency.DependenciesProvider) interface{} {
	return NewIncomeRepository(factory.FindRequiredByObject(db.DatabaseObject{}).(db.DatabaseService))
}

func (i *IncomeRepositoryObject) GetEntity() interface{} {
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
	FindById(id uuid.UUID) (model.Income, error)
	FindByHouseId(id uuid.UUID) ([]model.IncomeDto, error)
	ExistsById(id uuid.UUID) bool
	DeleteById(id uuid.UUID) error
	Update(id uuid.UUID, request model.UpdateIncomeRequest) error
}

func (i *IncomeRepositoryObject) Create(entity model.Income) (model.Income, error) {
	return entity, i.database.Create(&entity)
}

func (i *IncomeRepositoryObject) FindById(id uuid.UUID) (response model.Income, err error) {
	return response, i.database.FindById(&response, id)
}

func (i *IncomeRepositoryObject) FindByHouseId(id uuid.UUID) (response []model.IncomeDto, err error) {
	return response, i.database.FindBy(&response, "house_id = ?", id)
}

func (i *IncomeRepositoryObject) ExistsById(id uuid.UUID) bool {
	return i.database.Exists(id)
}

func (i *IncomeRepositoryObject) DeleteById(id uuid.UUID) error {
	return i.database.Delete(id)
}

func (i *IncomeRepositoryObject) Update(id uuid.UUID, request model.UpdateIncomeRequest) error {
	return i.database.Update(id, request)
}
