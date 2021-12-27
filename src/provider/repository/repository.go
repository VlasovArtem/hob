package repository

import (
	"fmt"
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/provider/model"
	"github.com/google/uuid"
	"reflect"
)

var (
	PaymentRepositoryType = reflect.TypeOf(ProviderRepositoryObject{})
	entity      = model.Provider{}
	DefaultUser = uuid.UUID{}
)

type ProviderRepositoryObject struct {
	database db.ModeledDatabase
}

func NewProviderRepository(database db.DatabaseService) ProviderRepository {
	return &ProviderRepositoryObject{
		db.ModeledDatabase{
			DatabaseService: database,
			Model:           entity,
		},
	}
}

func (p *ProviderRepositoryObject) Initialize(factory dependency.DependenciesProvider) interface{} {
	return NewProviderRepository(factory.FindRequiredByObject(db.DatabaseObject{}).(db.DatabaseService))
}

func (p *ProviderRepositoryObject) GetEntity() interface{} {
	return entity
}

type ProviderRepository interface {
	Create(provider model.Provider) (model.Provider, error)
	FindById(id uuid.UUID) (model.Provider, error)
	Delete(id uuid.UUID) error
	Update(entity model.Provider) error
	FindByUserId(id uuid.UUID) []model.ProviderDto
	FindByNameLikeAndUserId(namePattern string, page, limit int, userId uuid.UUID) []model.ProviderDto
	ExistsById(id uuid.UUID) bool
	ExistsByNameAndUserId(name string, userId uuid.UUID) bool
}

func (p *ProviderRepositoryObject) Create(provider model.Provider) (model.Provider, error) {
	if provider.UserId == DefaultUser {
		return provider, p.database.D().Omit("UserId", "User").Create(provider).Error
	}
	return provider, p.database.Create(provider)
}

func (p *ProviderRepositoryObject) FindById(id uuid.UUID) (provider model.Provider, err error) {
	return provider, p.database.FindById(&provider, id)
}

func (p *ProviderRepositoryObject) Delete(id uuid.UUID) (err error) {
	return p.database.Delete(id)
}

func (p *ProviderRepositoryObject) Update(entity model.Provider) error {
	return p.database.Update(entity.Id, entity)
}

func (p *ProviderRepositoryObject) FindByUserId(id uuid.UUID) (provider []model.ProviderDto) {
	_ = p.database.FindBy(&provider, "user_id = ? OR user_id IS NULL", id)

	return provider
}

func (p *ProviderRepositoryObject) FindByNameLikeAndUserId(namePattern string, page, limit int, userId uuid.UUID) (response []model.ProviderDto) {
	p.database.DM(model.Provider{}).
		Offset(page*limit).
		Limit(limit).
		Order("name asc").
		Find(&response, "name like ? AND (user_id = ? OR user_id IS NULL)", fmt.Sprintf("%%%s%%", namePattern), userId)

	return response
}

func (p *ProviderRepositoryObject) ExistsById(id uuid.UUID) bool {
	return p.database.Exists(id)
}

func (p *ProviderRepositoryObject) ExistsByNameAndUserId(name string, userId uuid.UUID) bool {
	return p.database.ExistsBy("name = ? and user_id = ?", name, userId)
}
