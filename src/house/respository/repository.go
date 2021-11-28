package respository

import (
	"common/dependency"
	"db"
	"github.com/google/uuid"
	"house/model"
)

type HouseRepositoryObject struct {
	database db.DatabaseService
}

func NewHouseRepository(database db.DatabaseService) HouseRepository {
	return &HouseRepositoryObject{database}
}

func (h *HouseRepositoryObject) Initialize(factory dependency.DependenciesFactory) {
	factory.Add(
		NewHouseRepository(
			factory.FindRequiredByObject(db.DatabaseObject{}).(db.DatabaseService),
		),
	)
}

func (h *HouseRepositoryObject) GetEntity() interface{} {
	return model.House{}
}

type HouseRepository interface {
	Create(entity model.House) (model.House, error)
	FindResponseById(id uuid.UUID) (model.HouseDto, error)
	FindResponseByUserId(id uuid.UUID) []model.HouseDto
	ExistsById(id uuid.UUID) bool
}

func (h *HouseRepositoryObject) Create(entity model.House) (model.House, error) {
	return entity, h.database.Create(&entity)
}

func (h *HouseRepositoryObject) FindResponseById(id uuid.UUID) (response model.HouseDto, err error) {
	return response, h.database.FindByIdModeled(model.House{}, &response, id)
}

func (h *HouseRepositoryObject) FindResponseByUserId(id uuid.UUID) (response []model.HouseDto) {
	h.database.DM(model.House{}).Where("user_id = ?", id).Find(&response)

	return response
}

func (h *HouseRepositoryObject) ExistsById(id uuid.UUID) bool {
	return h.database.ExistsById(model.House{}, id)
}
