package repository

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/pivotal/model"
	"github.com/google/uuid"
	"time"
)

type HousePivotalRepository struct {
	db db.ModeledDatabase
}

func NewHousePivotalRepository(database db.DatabaseService) PivotalRepository {
	return &HousePivotalRepository{
		db.ModeledDatabase{
			DatabaseService: database,
			Model:           model.HousePivotal{},
		},
	}
}

func (h *HousePivotalRepository) Initialize(factory dependency.DependenciesProvider) any {
	return NewHousePivotalRepository(dependency.FindRequiredDependency[db.DatabaseObject, db.DatabaseService](factory))
}

func (h *HousePivotalRepository) GetEntity() any {
	return model.HousePivotal{}
}

type GroupPivotalRepository struct {
	db db.ModeledDatabase
}

func NewGroupPivotalRepository(database db.DatabaseService) PivotalRepository {
	return &GroupPivotalRepository{
		db.ModeledDatabase{
			DatabaseService: database,
			Model:           model.GroupPivotal{},
		},
	}
}

func (g *GroupPivotalRepository) Initialize(factory dependency.DependenciesProvider) any {
	return NewGroupPivotalRepository(dependency.FindRequiredDependency[db.DatabaseObject, db.DatabaseService](factory))
}

func (g *GroupPivotalRepository) GetEntity() any {
	return model.GroupPivotal{}
}

type PivotalRepository interface {
	Create(entity any) (any, error)
	FindBySourceId(sourceId uuid.UUID, source any) error
	Update(sourceId uuid.UUID, total float64, latestIncomeUpdate time.Time, latestPaymentUpdate time.Time) error
}

func (h *HousePivotalRepository) Create(entity any) (any, error) {
	return entity, h.db.Create(&entity)
}

func (h *HousePivotalRepository) FindBySourceId(sourceId uuid.UUID, source any) error {
	return h.db.FindBy(&source, "house_id = ?", sourceId)
}

func (h *HousePivotalRepository) Update(id uuid.UUID, total float64, latestIncomeUpdate time.Time, latestPaymentUpdate time.Time) error {
	return update(h.db, id, total, latestIncomeUpdate, latestPaymentUpdate)
}

func (g *GroupPivotalRepository) Create(entity any) (any, error) {
	return entity, g.db.Create(&entity)
}

func (g *GroupPivotalRepository[T]) FindBySourceId(sourceId uuid.UUID, source any) error {
	return g.db.FindBy(&source, "group_id = ?", sourceId)
}

func (g *GroupPivotalRepository) Update(id uuid.UUID, total float64, latestIncomeUpdate time.Time, latestPaymentUpdate time.Time) error {
	return update(g.db, id, total, latestIncomeUpdate, latestPaymentUpdate)
}

func update(database db.ModeledDatabase, id uuid.UUID, total float64, latestIncomeUpdate time.Time, latestPaymentUpdate time.Time) error {
	return database.Update(id, struct {
		Total                   float64
		LatestIncomeUpdateDate  time.Time
		LatestPaymentUpdateDate time.Time
	}{
		Total:                   total,
		LatestIncomeUpdateDate:  latestIncomeUpdate,
		LatestPaymentUpdateDate: latestPaymentUpdate,
	})
}
