package repository

import (
	"github.com/VlasovArtem/hob/src/db"
	pivotalModel "github.com/VlasovArtem/hob/src/pivotal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PivotalRepository[T any] interface {
	Create(entity T) (T, error)
	FindBySourceId(sourceId uuid.UUID, source *T) error
	FindBySourceIdTransactional(db *gorm.DB, sourceId uuid.UUID, source *T) error
	UpdateTransactional(db *gorm.DB, sourceId uuid.UUID, pivotal pivotalModel.Pivotal) error
}

func update(db *gorm.DB, database db.modeledDatabase, id uuid.UUID, pivotal pivotalModel.Pivotal) error {
	return database.UpdateTransactional(db, id, pivotal)
}
