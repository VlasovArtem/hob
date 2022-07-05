package db

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type modeledDatabaseStruct[T any] struct {
	dependency.EntityProvider[T]
	DatabaseService
}

func NewTransactionalModeledDatabase[T any](model T, tx *gorm.DB) ModeledDatabase[T] {
	return &modeledDatabaseStruct[T]{
		EntityProvider: dependency.NewEntity[T](model),
		DatabaseService: &Database{&Provider{
			db: tx,
		}},
	}
}

func NewModeledDatabase[T any](model T, service DatabaseService) ModeledDatabase[T] {
	return &modeledDatabaseStruct[T]{
		EntityProvider:  dependency.NewEntity[T](model),
		DatabaseService: service,
	}
}

type ModeledDatabase[T any] interface {
	DatabaseService
	dependency.EntityProvider[T]
	Find(id uuid.UUID) (T, error)
	FindReceiver(receiver any, id uuid.UUID) error
	FindBy(query any, conditions ...any) (T, error)
	FindReceiverBy(receiver any, query any, conditions ...any) error
	First(id uuid.UUID) (T, error)
	FirstReceiver(receiver any, id uuid.UUID) error
	FirstBy(query any, conditions ...any) (T, error)
	FirstReceiverBy(receiver any, query any, conditions ...any) error
	Exists(id uuid.UUID) bool
	ExistsBy(query any, args ...any) (exists bool)
	Delete(id uuid.UUID) error
	DeleteBy(query any, args ...any) error
	Update(id uuid.UUID, entity any, omit ...string) error
	Modeled() *gorm.DB
}

func (m *modeledDatabaseStruct[T]) Find(id uuid.UUID) (model T, err error) {
	return model, m.DB().Model(model).Find(&model, id).Error
}

func (m *modeledDatabaseStruct[T]) FindReceiver(receiver any, id uuid.UUID) error {
	return m.Modeled().Find(receiver, id).Error
}

func (m *modeledDatabaseStruct[T]) FindBy(query any, conditions ...any) (model T, err error) {
	return model, m.Modeled().Where(query, conditions...).Find(&model).Error
}

func (m *modeledDatabaseStruct[T]) FindReceiverBy(receiver any, query any, conditions ...any) error {
	return m.Modeled().Where(query, conditions...).Find(receiver).Error
}

func (m *modeledDatabaseStruct[T]) First(id uuid.UUID) (model T, err error) {
	return model, m.Modeled().First(&model, id).Error
}

func (m *modeledDatabaseStruct[T]) FirstReceiver(receiver any, id uuid.UUID) error {
	return m.Modeled().First(receiver, id).Error
}

func (m *modeledDatabaseStruct[T]) FirstBy(query any, conditions ...any) (model T, err error) {
	return model, m.Modeled().Where(query, conditions...).First(&model).Error
}

func (m *modeledDatabaseStruct[T]) FirstReceiverBy(receiver any, query any, conditions ...any) error {
	return m.Modeled().Where(query, conditions...).First(receiver).Error
}

func (m *modeledDatabaseStruct[T]) Exists(id uuid.UUID) (exists bool) {
	tx := m.Modeled().Select("count(*) > 0").Where("id = ?", id)
	if err := tx.Find(&exists).Error; err != nil {
		log.Error().Err(err).Msg("")
	}
	return exists
}

func (m *modeledDatabaseStruct[T]) ExistsBy(query any, args ...any) (exists bool) {
	tx := m.Modeled().Select("count(*) > 0").Where(query, args...)
	if err := tx.Find(&exists).Error; err != nil {
		log.Error().Err(err).Msg("")
	}
	return exists
}

func (m *modeledDatabaseStruct[T]) Delete(id uuid.UUID) error {
	return m.DB().Delete(m.GetEntity(), id).Error
}

func (m *modeledDatabaseStruct[T]) DeleteBy(query any, args ...any) error {
	conditions := []any{query}
	conditions = append(conditions, args...)
	return m.DB().Delete(m.GetEntity(), conditions...).Error
}

func (m *modeledDatabaseStruct[T]) Update(id uuid.UUID, entity any, omit ...string) error {
	omitColumns := []string{"Id"}
	omitColumns = append(omitColumns, omit...)

	return m.Modeled().Where("id = ?", id).Omit(omitColumns...).Updates(entity).Error
}

func (m *modeledDatabaseStruct[T]) Modeled() *gorm.DB {
	return m.DBModeled(m.GetEntity())
}
