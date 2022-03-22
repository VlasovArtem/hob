package db

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ModeledDatabase struct {
	DatabaseService
	Model any
}

func (m *ModeledDatabase) Find(receiver any, id uuid.UUID) error {
	return m.FindByIdModeled(m.Model, receiver, id)
}

func (m *ModeledDatabase) FindBy(receiver any, query any, conditions ...any) error {
	return m.FindByModeled(m.Model, receiver, query, conditions...)
}

func (m *ModeledDatabase) First(receiver any, id uuid.UUID) error {
	return m.FindByIdModeled(m.Model, receiver, id)
}

func (m *ModeledDatabase) FirstBy(receiver any, query any, conditions ...any) error {
	return m.FirstByModeled(m.Model, receiver, query, conditions...)
}

func (m *ModeledDatabase) Exists(id uuid.UUID) bool {
	return m.ExistsById(m.Model, id)
}

func (m *ModeledDatabase) ExistsBy(query any, args ...any) (exists bool) {
	return m.ExistsByQuery(m.Model, query, args...)
}

func (m *ModeledDatabase) Delete(id uuid.UUID) error {
	return m.DeleteById(m.Model, id)
}

func (m *ModeledDatabase) Update(id uuid.UUID, entity any, omit ...string) error {
	return m.UpdateById(m.Model, id, entity, omit...)
}

func (m *ModeledDatabase) Modeled() *gorm.DB {
	return m.DM(m.Model)
}
