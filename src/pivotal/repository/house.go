package repository

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/pivotal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HousePivotalRepository[T model.HousePivotal] struct {
	db db.modeledDatabase
}

func NewHousePivotalRepository(database db.DatabaseService) PivotalRepository[model.HousePivotal] {
	return &HousePivotalRepository[model.HousePivotal]{
		db.modeledDatabase{
			DatabaseService: database,
			Model:           model.HousePivotal{},
		},
	}
}

func (h *HousePivotalRepository[T]) Initialize(factory dependency.DependenciesProvider) any {
	return NewHousePivotalRepository(dependency.FindRequiredDependency[db.Database, db.DatabaseService](factory))
}

func (h *HousePivotalRepository[T]) GetEntity() any {
	return model.HousePivotal{}
}

func (h *HousePivotalRepository[T]) Create(entity T) (T, error) {
	return entity, h.db.DB().Omit("Groups.*").Create(&entity).Error
}

func (h *HousePivotalRepository[T]) FindBySourceId(sourceId uuid.UUID, source *T) error {
	return h.db.FindBy(&source, "house_id = ?", sourceId)
}

func (h *HousePivotalRepository[T]) FindBySourceIdTransactional(db *gorm.DB, sourceId uuid.UUID, source *T) error {
	return h.db.FindBy(&source, "house_id = ?", sourceId)
}

func (h *HousePivotalRepository[T]) UpdateTransactional(db *gorm.DB, sourceId uuid.UUID, pivotal model.Pivotal) error {
	return update(db, h.db, sourceId, pivotal)
}
