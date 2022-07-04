package repository

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/group/model"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type GroupRepositoryObject struct {
	db.ModeledDatabase[model.Group]
}

func (g *GroupRepositoryObject) Initialize(factory dependency.DependenciesProvider) any {
	return NewGroupRepository(dependency.FindRequiredDependency[db.Database, db.DatabaseService](factory))
}

func NewGroupRepository(database db.DatabaseService) GroupRepository {
	return &GroupRepositoryObject{
		ModeledDatabase: db.NewModeledDatabase(model.Group{}, database),
	}
}

type GroupRepository interface {
	db.ModeledDatabase[model.Group]
	CreateBatch(entities []model.Group) ([]model.Group, error)
	FindByOwnerId(ownerId uuid.UUID) (response []model.GroupDto)
	ExistsByIds(ids []uuid.UUID) bool
}

func (g *GroupRepositoryObject) CreateBatch(entities []model.Group) ([]model.Group, error) {
	return entities, g.Create(&entities)
}

func (g *GroupRepositoryObject) FindByOwnerId(ownerId uuid.UUID) (response []model.GroupDto) {
	if err := g.FindReceiverBy(&response, "owner_id = ?", ownerId); err != nil {
		response = []model.GroupDto{}
	}

	return
}

func (g *GroupRepositoryObject) ExistsByIds(ids []uuid.UUID) bool {
	var count int64
	if err := g.Modeled().Where("id IN ?", ids).Count(&count).Error; err != nil {
		log.Error().Err(err)
	}

	return int64(len(ids)) == count
}
