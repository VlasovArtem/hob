package respository

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/house/model"
	"github.com/google/uuid"
)

type HouseRepositoryObject struct {
	database db.DatabaseService
}

func NewHouseRepository(database db.DatabaseService) HouseRepository {
	return &HouseRepositoryObject{database}
}

func (h *HouseRepositoryObject) Initialize(factory dependency.DependenciesProvider) interface{} {
	return NewHouseRepository(factory.FindRequiredByObject(db.DatabaseObject{}).(db.DatabaseService))
}

func (h *HouseRepositoryObject) GetEntity() interface{} {
	return model.House{}
}

type HouseRepository interface {
	Create(entity model.House) (model.House, error)
	FindDtoById(id uuid.UUID) (model.HouseDto, error)
	FindResponseByUserId(id uuid.UUID) []model.HouseDto
	ExistsById(id uuid.UUID) bool
	DeleteById(id uuid.UUID) error
	Update(entity model.House) error
}

func (h *HouseRepositoryObject) Create(entity model.House) (model.House, error) {
	return entity, h.database.Create(&entity)
}

func (h *HouseRepositoryObject) FindDtoById(id uuid.UUID) (response model.HouseDto, err error) {
	return response, h.database.FindByIdModeled(model.House{}, &response, id)
}

func (h *HouseRepositoryObject) FindResponseByUserId(id uuid.UUID) (response []model.HouseDto) {
	h.database.DM(model.House{}).Where("user_id = ?", id).Find(&response)

	return response
}

func (h *HouseRepositoryObject) ExistsById(id uuid.UUID) bool {
	return h.database.ExistsById(model.House{}, id)
}

func (h *HouseRepositoryObject) DeleteById(id uuid.UUID) error {
	return h.database.DeleteById(model.House{}, id)
}

func (h *HouseRepositoryObject) Update(entity model.House) error {
	return h.database.UpdateById(model.House{}, entity.Id, entity, "UserId", "User")
}
