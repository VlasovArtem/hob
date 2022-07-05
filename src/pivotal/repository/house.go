package repository

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/pivotal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HousePivotalRepository struct {
	db.ModeledDatabase[model.HousePivotal]
}

func NewHousePivotalRepository(database db.DatabaseService) PivotalRepository[model.HousePivotal] {
	return &HousePivotalRepository{db.NewModeledDatabase(model.HousePivotal{}, database)}
}

func (h *HousePivotalRepository) Initialize(factory dependency.DependenciesProvider) any {
	return NewHousePivotalRepository(dependency.FindRequiredDependency[db.Database, db.DatabaseService](factory))
}

func (h *HousePivotalRepository) FindBySourceId(sourceId uuid.UUID, source *model.HousePivotal) error {
	return h.FindReceiverBy(&source, "house_id = ?", sourceId)
}

func (h *HousePivotalRepository) Transactional(tx *gorm.DB) PivotalRepository[model.HousePivotal] {
	return &HousePivotalRepository{db.NewTransactionalModeledDatabase(h.GetEntity(), tx)}
}
