package respository

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/house/model"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"reflect"
)

var (
	HouseRepositoryType = reflect.TypeOf(HouseRepositoryObject{})
	entity              = model.House{}
)

type HouseRepositoryObject struct {
	database db.ModeledDatabase
}

func NewHouseRepository(database db.DatabaseService) HouseRepository {
	return &HouseRepositoryObject{
		db.ModeledDatabase{
			DatabaseService: database,
			Model:           entity,
		},
	}
}

func (h *HouseRepositoryObject) Initialize(factory dependency.DependenciesProvider) any {
	return NewHouseRepository(factory.FindRequiredByObject(db.DatabaseObject{}).(db.DatabaseService))
}

func (h *HouseRepositoryObject) GetEntity() any {
	return entity
}

type HouseRepository interface {
	Create(entity model.House) (model.House, error)
	FindDtoById(id uuid.UUID) (model.HouseDto, error)
	FindResponseByUserId(id uuid.UUID) []model.HouseDto
	ExistsById(id uuid.UUID) bool
	DeleteById(id uuid.UUID) error
	Update(id uuid.UUID, request model.UpdateHouseRequest) error
}

func (h *HouseRepositoryObject) Create(entity model.House) (model.House, error) {
	return entity, h.database.Create(&entity)
}

func (h *HouseRepositoryObject) FindDtoById(id uuid.UUID) (response model.HouseDto, err error) {
	return response, h.database.Find(&response, id)
}

func (h *HouseRepositoryObject) FindResponseByUserId(id uuid.UUID) (response []model.HouseDto) {
	if err := h.database.FindBy(&response, "user_id = ?", id); err != nil {
		log.Err(err)
	}
	return response
}

func (h *HouseRepositoryObject) ExistsById(id uuid.UUID) bool {
	return h.database.Exists(id)
}

func (h *HouseRepositoryObject) DeleteById(id uuid.UUID) error {
	return h.database.Delete(id)
}

func (h *HouseRepositoryObject) Update(id uuid.UUID, request model.UpdateHouseRequest) error {
	return h.database.Update(id, request)
}
