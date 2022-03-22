package repository

import (
	"github.com/VlasovArtem/hob/src/common"
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/db"
	groupModel "github.com/VlasovArtem/hob/src/group/model"
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
	db db.ModeledDatabase
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
	FindById(id uuid.UUID) (model.House, error)
	FindByUserId(id uuid.UUID) []model.House
	ExistsById(id uuid.UUID) bool
	DeleteById(id uuid.UUID) error
	Update(id uuid.UUID, request model.UpdateHouseRequest) error
}

func (h *HouseRepositoryObject) Create(entity model.House) (model.House, error) {
	return entity, h.db.Create(&entity)
}

func (h *HouseRepositoryObject) FindById(id uuid.UUID) (response model.House, err error) {
	response.Id = id
	if err = h.db.D().Preload("Groups").First(&response).Error; err != nil {
		return model.House{}, err
	}
	return response, err
}

func (h *HouseRepositoryObject) FindByUserId(id uuid.UUID) (response []model.House) {
	if err := h.db.D().Preload("Groups").Where("user_id = ?", id).Find(&response).Error; err != nil {
		log.Error().Err(err)
	}

	return response
}

func (h *HouseRepositoryObject) ExistsById(id uuid.UUID) bool {
	return h.db.Exists(id)
}

func (h *HouseRepositoryObject) DeleteById(id uuid.UUID) error {
	return h.db.Delete(id)
}

func (h *HouseRepositoryObject) Update(id uuid.UUID, request model.UpdateHouseRequest) error {
	err := h.db.Update(id, struct {
		Name        string
		CountryCode string
		City        string
		StreetLine1 string
		StreetLine2 string
	}{
		request.Name,
		request.CountryCode,
		request.City,
		request.StreetLine1,
		request.StreetLine2,
	})

	if err != nil {
		return err
	}

	entity, err := h.FindById(id)

	if err != nil {
		return err
	}

	var groups = common.Map(
		request.GroupIds,
		[]groupModel.Group{},
		func(id uuid.UUID) groupModel.Group {
			return groupModel.Group{Id: id}
		})

	return h.db.DM(&entity).Association("Groups").Replace(groups)
}
