package cache

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/pivotal/model"
	"github.com/google/uuid"
)

type PivotalCacheObject struct {
	pivotalCache map[uuid.UUID]*model.Pivotal
	invalidCache map[uuid.UUID]bool
}

func NewPivotalCache() PivotalCache {
	return &PivotalCacheObject{
		pivotalCache: make(map[uuid.UUID]*model.Pivotal),
		invalidCache: make(map[uuid.UUID]bool),
	}
}

type PivotalCache interface {
	Find(id uuid.UUID) *model.Pivotal
	Add(id uuid.UUID, pivotal *model.Pivotal)
	Invalidate(id uuid.UUID)
}

func (p *PivotalCacheObject) Initialize(factory dependency.DependenciesProvider) any {
	return NewPivotalCache()
}

func (p *PivotalCacheObject) Find(id uuid.UUID) *model.Pivotal {
	return p.pivotalCache[id]
}

func (p *PivotalCacheObject) Add(id uuid.UUID, pivotal *model.Pivotal) {
	p.pivotalCache[id] = pivotal
}

func (p *PivotalCacheObject) Invalidate(id uuid.UUID) {
	p.invalidCache[id] = true
}
