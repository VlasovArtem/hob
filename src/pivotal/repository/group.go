package repository

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/pivotal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GroupPivotalRepository struct {
	db.ModeledDatabase[model.GroupPivotal]
}

func NewGroupPivotalRepository(database db.DatabaseService) PivotalRepository[model.GroupPivotal] {
	return &GroupPivotalRepository{db.NewModeledDatabase(model.GroupPivotal{}, database)}
}

func (g *GroupPivotalRepository) Initialize(factory dependency.DependenciesProvider) any {
	return NewGroupPivotalRepository(dependency.FindRequiredDependency[db.Database, db.DatabaseService](factory))
}

func (g *GroupPivotalRepository) FindBySourceId(sourceId uuid.UUID, source *model.GroupPivotal) error {
	return g.FindReceiverBy(&source, "group_id = ?", sourceId)
}

func (g *GroupPivotalRepository) Transactional(tx *gorm.DB) PivotalRepository[model.GroupPivotal] {
	return &GroupPivotalRepository{db.NewTransactionalModeledDatabase(g.GetEntity(), tx)}
}
