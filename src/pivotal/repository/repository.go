package repository

import (
	"github.com/VlasovArtem/hob/src/common/transactional"
	"github.com/VlasovArtem/hob/src/db"
	"github.com/google/uuid"
)

type PivotalRepository[T any] interface {
	db.ModeledDatabase[T]
	FindBySourceId(sourceId uuid.UUID, source *T) error
	transactional.Transactional[PivotalRepository[T]]
}
