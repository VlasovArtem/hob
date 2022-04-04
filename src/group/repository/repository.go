package repository

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/group/model"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var entity = model.Group{}

type GroupRepositoryObject struct {
	database db.ModeledDatabase
}

func (g *GroupRepositoryObject) Initialize(factory dependency.DependenciesProvider) any {
	return NewGroupRepository(dependency.FindRequiredDependency[db.DatabaseObject, db.DatabaseService](factory))
}

func (g *GroupRepositoryObject) GetEntity() any {
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
	CreateBatch(entities []model.Group) ([]model.Group, error)
	FindById(id uuid.UUID) (model.GroupDto, error)
	FindByOwnerId(ownerId uuid.UUID) (response []model.GroupDto)
	ExistsById(id uuid.UUID) bool
	ExistsByIds(ids []uuid.UUID) bool
	DeleteById(id uuid.UUID) error
	Update(id uuid.UUID, request model.UpdateGroupRequest) error
}

func (g *GroupRepositoryObject) Create(entity model.Group) (model.Group, error) {
	return entity, g.database.Create(&entity)
}

func (g *GroupRepositoryObject) CreateBatch(entities []model.Group) ([]model.Group, error) {
	return entities, g.database.Create(&entities)
}

func (g *GroupRepositoryObject) FindById(id uuid.UUID) (response model.GroupDto, err error) {
	return response, g.database.Find(&response, id)
}

func (g *GroupRepositoryObject) FindByOwnerId(ownerId uuid.UUID) (response []model.GroupDto) {
	err := g.database.FindBy(&response, "owner_id = ?", ownerId)

	if err != nil {
		return []model.GroupDto{}
	}

	return response
}

func (g *GroupRepositoryObject) ExistsById(id uuid.UUID) bool {
	return g.database.Exists(id)
}

func (g *GroupRepositoryObject) ExistsByIds(ids []uuid.UUID) bool {
	var count int64
	err := g.database.Modeled().Where("id IN ?", ids).Count(&count).Error

	if err != nil {
		log.Error().Err(err)
	}

	return int64(len(ids)) == count
}

func (g *GroupRepositoryObject) DeleteById(id uuid.UUID) error {
	return g.database.Delete(id)
}

func (g *GroupRepositoryObject) Update(id uuid.UUID, request model.UpdateGroupRequest) error {
	return g.database.Update(id, request)
}
