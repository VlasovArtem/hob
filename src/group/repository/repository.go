package repository

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/group/model"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"reflect"
)

var (
	GroupRepositoryType = reflect.TypeOf(GroupRepositoryObject{})
	entity              = model.Group{}
)

type GroupRepositoryObject struct {
	database db.ModeledDatabase
}

func (i *GroupRepositoryObject) Initialize(factory dependency.DependenciesProvider) any {
	return NewGroupRepository(factory.FindRequiredByObject(db.DatabaseObject{}).(db.DatabaseService))
}

func (i *GroupRepositoryObject) GetEntity() any {
	return entity
}

func NewGroupRepository(database db.DatabaseService) GroupRepository {
	return &GroupRepositoryObject{
		db.ModeledDatabase{
			DatabaseService: database,
			Model:           entity,
		},
	}
}

type GroupRepository interface {
	Create(entity model.Group) (model.Group, error)
	FindById(id uuid.UUID) (model.GroupDto, error)
	FindByOwnerId(ownerId uuid.UUID) (response []model.GroupDto)
	ExistsById(id uuid.UUID) bool
	ExistsByIds(ids []uuid.UUID) bool
	DeleteById(id uuid.UUID) error
	Update(id uuid.UUID, request model.UpdateGroupRequest) error
}

func (i *GroupRepositoryObject) Create(entity model.Group) (model.Group, error) {
	return entity, i.database.Create(&entity)
}

func (i *GroupRepositoryObject) FindById(id uuid.UUID) (response model.GroupDto, err error) {
	return response, i.database.Find(&response, id)
}

func (i *GroupRepositoryObject) FindByOwnerId(ownerId uuid.UUID) (response []model.GroupDto) {
	err := i.database.FindBy(&response, "owner_id = ?", ownerId)

	if err != nil {
		return []model.GroupDto{}
	}

	return response
}

func (i *GroupRepositoryObject) ExistsById(id uuid.UUID) bool {
	return i.database.Exists(id)
}

func (i *GroupRepositoryObject) ExistsByIds(ids []uuid.UUID) bool {
	var count int64
	err := i.database.Modeled().Where("id IN ?", ids).Count(&count).Error

	if err != nil {
		log.Error().Err(err)
	}

	return int64(len(ids)) == count
}

func (i *GroupRepositoryObject) DeleteById(id uuid.UUID) error {
	return i.database.Delete(id)
}

func (i *GroupRepositoryObject) Update(id uuid.UUID, request model.UpdateGroupRequest) error {
	return i.database.Update(id, request)
}
