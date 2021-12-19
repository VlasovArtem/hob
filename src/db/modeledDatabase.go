package db

import (
	"github.com/google/uuid"
)

type ModeledDatabase struct {
	DatabaseService
	Model interface{}
}

func (m *ModeledDatabase) Find(receiver interface{}, id uuid.UUID) error {
	return m.FindByIdModeled(m.Model, receiver, id)
}

func (m *ModeledDatabase) FindBy(receiver interface{}, query interface{}, conditions ...interface{}) error {
	return m.FindByModeled(m.Model, receiver, query, conditions...)
}

func (m *ModeledDatabase) Exists(id uuid.UUID) bool {
	return m.ExistsById(m.Model, id)
}

func (m *ModeledDatabase) ExistsBy(query interface{}, args ...interface{}) (exists bool) {
	return m.ExistsByQuery(m.Model, query, args...)
}

func (m *ModeledDatabase) Delete(id uuid.UUID) error {
	return m.DeleteById(m.Model, id)
}

func (m *ModeledDatabase) Update(id uuid.UUID, entity interface{}, omit ...string) error {
	return m.UpdateById(m.Model, id, entity, omit...)
}
