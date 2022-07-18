package repository

import (
	"github.com/VlasovArtem/hob/src/common"
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/transactional"
	"github.com/VlasovArtem/hob/src/db"
	groupModel "github.com/VlasovArtem/hob/src/group/model"
	"github.com/VlasovArtem/hob/src/house/model"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type HouseRepositoryStr struct {
	db.ModeledDatabase[model.House]
}

func NewHouseRepository(database db.DatabaseService) HouseRepository {
	return &HouseRepositoryStr{
		ModeledDatabase: db.NewModeledDatabase(model.House{}, database),
	}
}

func (h *HouseRepositoryStr) GetRequiredDependencies() []dependency.Requirements {
	return []dependency.Requirements{
		dependency.FindNameAndType(db.Database{}),
	}
}

func (h *HouseRepositoryStr) Initialize(factory dependency.DependenciesProvider) any {
	return NewHouseRepository(factory.FindRequiredByObject(db.Database{}).(db.DatabaseService))
}

type HouseRepository interface {
	db.ModeledDatabase[model.House]
	transactional.Transactional[HouseRepository]
	CreateBatch(entities []model.House) ([]model.House, error)
	FindById(id uuid.UUID) (response model.House, err error)
	FindByUserId(id uuid.UUID) []model.House
	FindHousesByGroupId(groupId uuid.UUID) []model.House
	FindHousesByGroupIds(groupIds []uuid.UUID) []model.House
	UpdateByRequest(id uuid.UUID, request model.UpdateHouseRequest) error
}

func (h *HouseRepositoryStr) CreateBatch(entities []model.House) ([]model.House, error) {
	return entities, h.DB().Omit("Groups.*").Create(&entities).Error
}

func (h *HouseRepositoryStr) FindById(id uuid.UUID) (response model.House, err error) {
	response.Id = id
	if err = h.DB().Preload("Groups").First(&response).Error; err != nil {
		return model.House{}, err
	}
	return response, err
}

func (h *HouseRepositoryStr) FindByUserId(id uuid.UUID) (response []model.House) {
	if err := h.DB().Preload("Groups").Where("user_id = ?", id).Find(&response).Error; err != nil {
		log.Error().Err(err)
	}

	return response
}

func (h *HouseRepositoryStr) UpdateByRequest(id uuid.UUID, request model.UpdateHouseRequest) error {
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

	entity, err := h.FindById(id)

	if err != nil {
		return err
	}

	groups := common.MapSlice(request.GroupIds, groupModel.GroupIdToGroup)

	return h.DBModeled(&entity).Association("Groups").Replace(groups)
}

func (h *HouseRepositoryStr) FindHousesByGroupId(groupId uuid.UUID) (response []model.House) {
	if err := h.DB().
		Preload("Groups").
		Joins("FULL JOIN house_groups hg ON hg.house_id = houses.id").
		Where("hg.group_id = ?", groupId).
		Find(&response).Error; err != nil {
		log.Error().Err(err)
	}

	return response
}

func (h *HouseRepositoryStr) FindHousesByGroupIds(groupIds []uuid.UUID) (response []model.House) {
	if err := h.DB().
		Preload("Groups").
		Joins("FULL JOIN house_groups hg ON hg.house_id = houses.id").
		Where("hg.group_id IN (?)", groupIds).
		Find(&response).Error; err != nil {
		log.Error().Err(err)
	}

	return response
}

func (h *HouseRepositoryStr) Transactional(tx *gorm.DB) HouseRepository {
	return &HouseRepositoryStr{
		ModeledDatabase: db.NewTransactionalModeledDatabase(h.GetEntity(), tx),
	}
}
