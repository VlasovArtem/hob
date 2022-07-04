package repository

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/pivotal/model"
	"github.com/google/uuid"
	"time"
)

type GroupPivotalRepository[T model.GroupPivotal] struct {
	db db.modeledDatabase
}

func NewGroupPivotalRepository(database db.DatabaseService) PivotalRepository[model.GroupPivotal] {
	return &GroupPivotalRepository[model.GroupPivotal]{
		db.modeledDatabase{
			DatabaseService: database,
			Model:           model.GroupPivotal{},
		},
	}
}

func (g *GroupPivotalRepository[T]) Initialize(factory dependency.DependenciesProvider) any {
	return NewGroupPivotalRepository(dependency.FindRequiredDependency[db.Database, db.DatabaseService](factory))
}

func (g *GroupPivotalRepository[T]) GetEntity() any {
	return model.GroupPivotal{}
}

func (g *GroupPivotalRepository[T]) Create(entity T) (T, error) {
	return entity, g.db.DB().Omit("Groups.*").Create(&entity).Error
}

func (g *GroupPivotalRepository[T]) FindBySourceId(sourceId uuid.UUID, source *T) error {
	return g.db.FindBy(&source, "group_id = ?", sourceId)
}

func (g *GroupPivotalRepository[T]) Update(id uuid.UUID, total float64, latestIncomeUpdate time.Time, latestPaymentUpdate time.Time) error {
	return update(g.db, id, total, latestIncomeUpdate, latestPaymentUpdate)
}
