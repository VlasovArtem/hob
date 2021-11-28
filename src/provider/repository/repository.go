package repository

import (
	"fmt"
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/provider/model"
	"github.com/google/uuid"
)

type ProviderRepositoryObject struct {
	database db.DatabaseService
}

func NewProviderRepository(database db.DatabaseService) ProviderRepository {
	return &ProviderRepositoryObject{database}
}

func (p *ProviderRepositoryObject) Initialize(factory dependency.DependenciesFactory) {
	factory.Add(
		NewProviderRepository(factory.FindRequiredByObject(db.DatabaseObject{}).(db.DatabaseService)),
	)
}

func (p *ProviderRepositoryObject) GetEntity() interface{} {
	return model.Provider{}
}

type ProviderRepository interface {
	Create(provider model.Provider) (model.Provider, error)
	FindById(id uuid.UUID) (model.Provider, error)
	FindByNameLike(namePattern string, page int, limit int) []model.Provider
	ExistsById(id uuid.UUID) bool
	ExistsByName(name string) bool
}

func (p *ProviderRepositoryObject) Create(provider model.Provider) (model.Provider, error) {
	return provider, p.database.Create(&provider)
}

func (p *ProviderRepositoryObject) FindById(id uuid.UUID) (provider model.Provider, err error) {
	return provider, p.database.FindById(&provider, id)
}

func (p *ProviderRepositoryObject) FindByNameLike(namePattern string, page int, limit int) (response []model.Provider) {
	p.database.D().
		Offset(page*limit).
		Limit(limit).
		Order("name asc").
		Find(&response, "name like ?", fmt.Sprintf("%%%s%%", namePattern))

	return response
}

func (p *ProviderRepositoryObject) ExistsById(id uuid.UUID) bool {
	return p.database.ExistsById(model.Provider{}, id)
}

func (p *ProviderRepositoryObject) ExistsByName(name string) bool {
	return p.database.ExistsByQuery(model.Provider{}, "name = ?", name)
}
