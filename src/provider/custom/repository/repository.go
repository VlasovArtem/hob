package repository

import (
	"common/dependency"
	"db"
	"github.com/google/uuid"
	"provider/custom/model"
)

type CustomProviderRepositoryObject struct {
	database db.DatabaseService
}

func NewCustomProviderRepository(database db.DatabaseService) CustomProviderRepository {
	return &CustomProviderRepositoryObject{database}
}

func (c *CustomProviderRepositoryObject) Initialize(factory dependency.DependenciesFactory) {
	factory.Add(
		NewCustomProviderRepository(factory.FindRequiredByObject(db.DatabaseObject{}).(db.DatabaseService)),
	)
}

func (c *CustomProviderRepositoryObject) GetEntity() interface{} {
	return model.CustomProvider{}
}

type CustomProviderRepository interface {
	Create(provider model.CustomProvider) (model.CustomProvider, error)
	FindById(id uuid.UUID) (model.CustomProvider, error)
	FindByUserId(id uuid.UUID) []model.CustomProvider
	ExistsById(id uuid.UUID) bool
	ExistsByNameAndUserId(name string, userId uuid.UUID) bool
}

func (c *CustomProviderRepositoryObject) Create(provider model.CustomProvider) (model.CustomProvider, error) {
	return provider, c.database.Create(&provider)
}

func (c *CustomProviderRepositoryObject) FindById(id uuid.UUID) (provider model.CustomProvider, err error) {
	return provider, c.database.FindById(&provider, id)
}

func (c *CustomProviderRepositoryObject) FindByUserId(id uuid.UUID) (provider []model.CustomProvider) {
	c.database.DM(provider).
		Where("user_id = ?", id).
		Find(&provider)

	return provider
}

func (c *CustomProviderRepositoryObject) ExistsById(id uuid.UUID) bool {
	return c.database.ExistsById(model.CustomProvider{}, id)
}

func (c *CustomProviderRepositoryObject) ExistsByNameAndUserId(name string, userId uuid.UUID) bool {
	return c.database.ExistsByQuery(model.CustomProvider{}, "name = ? and user_id = ?", name, userId)
}
