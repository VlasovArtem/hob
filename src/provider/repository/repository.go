package repository

import (
	"fmt"
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/transactional"
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/provider/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProviderRepositoryStr struct {
	db.ModeledDatabase[model.Provider]
}

func NewProviderRepository(database db.DatabaseService) ProviderRepository {
	return &ProviderRepositoryStr{
		ModeledDatabase: db.NewModeledDatabase(model.Provider{}, database),
	}
}

func (p *ProviderRepositoryStr) Initialize(factory dependency.DependenciesProvider) any {
	return NewProviderRepository(dependency.FindRequiredDependency[db.Database, db.DatabaseService](factory))
}

type ProviderRepository interface {
	db.ModeledDatabase[model.Provider]
	transactional.Transactional[ProviderRepository]
	CreateEntity(provider model.Provider) (model.Provider, error)
	FindByUserId(id uuid.UUID) []model.ProviderDto
	FindByNameLikeAndUserId(namePattern string, page, limit int, userId uuid.UUID) []model.ProviderDto
	ExistsByNameAndUserId(name string, userId uuid.UUID) bool
}

func (p *ProviderRepositoryStr) CreateEntity(provider model.Provider) (model.Provider, error) {
	if provider.UserId == uuid.Nil {
		return provider, p.Create(&provider, "UserId", "User")
	}
	return provider, p.Create(provider)
}

func (p *ProviderRepositoryStr) FindByUserId(id uuid.UUID) (providers []model.ProviderDto) {
	_ = p.FindReceiverBy(&providers, "user_id = ? OR user_id IS NULL", id)

	return providers
}

func (p *ProviderRepositoryStr) FindByNameLikeAndUserId(namePattern string, page, limit int, userId uuid.UUID) (response []model.ProviderDto) {
	p.Modeled().
		Offset(page*limit).
		Limit(limit).
		Order("name asc").
		Find(&response, "name like ? AND (user_id = ? OR user_id IS NULL)", fmt.Sprintf("%%%s%%", namePattern), userId)

	return response
}

func (p *ProviderRepositoryStr) ExistsByNameAndUserId(name string, userId uuid.UUID) bool {
	return p.ExistsBy("name = ? and user_id = ?", name, userId)
}

func (p *ProviderRepositoryStr) Transactional(tx *gorm.DB) ProviderRepository {
	return &ProviderRepositoryStr{
		ModeledDatabase: db.NewTransactionalModeledDatabase(p.GetEntity(), tx),
	}
}
