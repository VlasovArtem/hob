package cache

import (
	"github.com/VlasovArtem/hob/src/pivotal/model"
	"github.com/google/uuid"
)

type PivotalCacheObject struct {
	pivotalCache map[uuid.UUID]*model.Pivotal
	invalidCache map[uuid.UUID]bool
}

type PivotalCache interface {
	Find(id uuid.UUID) *model.Pivotal
	Add(id uuid.UUID, pivotal *model.Pivotal)
	Invalidate(id uuid.UUID)
}
