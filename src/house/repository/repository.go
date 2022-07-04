package repository

import (
	"github.com/VlasovArtem/hob/src/common"
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/db"
	groupModel "github.com/VlasovArtem/hob/src/group/model"
	"github.com/VlasovArtem/hob/src/house/model"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type houseRepositoryStruct struct {
	db.ModeledDatabase[model.House]
}

func NewHouseRepository(database db.DatabaseService) HouseRepository {
	return &houseRepositoryStruct{
		ModeledDatabase: db.NewModeledDatabase(model.House{}, database),
	}
}

func (h *houseRepositoryStruct) Initialize(factory dependency.DependenciesProvider) any {
	return NewHouseRepository(factory.FindRequiredByObject(db.Database{}).(db.DatabaseService))
}

type HouseRepository interface {
	db.ModeledDatabase[model.House]
	CreateBatch(entities []model.House) ([]model.House, error)
	FindById(id uuid.UUID) (response model.House, err error)
	FindByUserId(id uuid.UUID) []model.House
	FindHousesByGroupId(groupId uuid.UUID) []model.House
	UpdateByRequest(id uuid.UUID, request model.UpdateHouseRequest) error
}

func (h *houseRepositoryStruct) CreateBatch(entities []model.House) ([]model.House, error) {
	return entities, h.DB().Omit("Groups.*").Create(&entities).Error
}

func (h *houseRepositoryStruct) FindById(id uuid.UUID) (response model.House, err error) {
	response.Id = id
	if err = h.DB().Preload("Groups").First(&response).Error; err != nil {
		return model.House{}, err
	}
	return response, err
}

func (h *houseRepositoryStruct) FindByUserId(id uuid.UUID) (response []model.House) {
	if err := h.DB().Preload("Groups").Where("user_id = ?", id).Find(&response).Error; err != nil {
		log.Error().Err(err)
	}

	return response
}

func (h *houseRepositoryStruct) UpdateByRequest(id uuid.UUID, request model.UpdateHouseRequest) error {
	err := h.Update(id, struct {
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

	_, err = h.FindById(id)

	if err != nil {
		return err
	}

	if groups := common.MapSlice(request.GroupIds, groupModel.GroupIdToGroup); len(groups) > 0 {
		return h.Modeled().Association("Groups").Replace(groups)
	}

	return nil
}

func (h *houseRepositoryStruct) FindHousesByGroupId(groupId uuid.UUID) (response []model.House) {
	if err := h.DB().
		Preload("Groups").
		Joins("FULL JOIN house_groups hg ON hg.house_id = houses.id").
		Where("hg.group_id = ?", groupId).
		Find(&response).Error; err != nil {
		log.Error().Err(err)
	}

	return response
}
