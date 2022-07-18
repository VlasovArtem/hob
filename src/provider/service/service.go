package service

import (
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/common/dependency"
	int_errors "github.com/VlasovArtem/hob/src/common/int-errors"
	"github.com/VlasovArtem/hob/src/common/transactional"
	"github.com/VlasovArtem/hob/src/provider/model"
	"github.com/VlasovArtem/hob/src/provider/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProviderServiceStr struct {
	repository repository.ProviderRepository
}

func NewProviderService(repository repository.ProviderRepository) ProviderService {
	return &ProviderServiceStr{repository}
}

func (p *ProviderServiceStr) GetRequiredDependencies() []dependency.Requirements {
	return []dependency.Requirements{
		dependency.FindNameAndType(repository.ProviderRepositoryStr{}),
	}
}

func (p *ProviderServiceStr) Initialize(factory dependency.DependenciesProvider) any {
	return NewProviderService(dependency.FindRequiredDependency[repository.ProviderRepositoryStr, repository.ProviderRepository](factory))
}

type ProviderService interface {
	transactional.Transactional[ProviderService]
	Add(request model.CreateProviderRequest) (dto model.ProviderDto, err error)
	ExistsById(id uuid.UUID) bool
	Update(id uuid.UUID, request model.UpdateProviderRequest) error
	Delete(id uuid.UUID) error
	FindById(id uuid.UUID) (dto model.ProviderDto, err error)
	FindByUserId(id uuid.UUID) []model.ProviderDto
	FindByNameLikeAndUserId(namePattern string, userId uuid.UUID, page, size int) []model.ProviderDto
}

func (p *ProviderServiceStr) Add(request model.CreateProviderRequest) (dto model.ProviderDto, err error) {
	if p.repository.ExistsByNameAndUserId(request.Name, request.UserId) {
		return dto, errors.New(fmt.Sprintf("provider with name '%s' for user already exists", request.Name))
	}

	entity := request.ToEntity()
	if err = p.repository.Create(&entity); err != nil {
		return dto, err
	} else {
		return entity.ToDto(), err
	}
}

func (p *ProviderServiceStr) FindById(id uuid.UUID) (dto model.ProviderDto, err error) {
	if err := p.repository.FindReceiver(&dto, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto, notFoundError(id)
		}
		return dto, err
	}
	return
}

func (p *ProviderServiceStr) FindByUserId(id uuid.UUID) (response []model.ProviderDto) {
	return p.repository.FindByUserId(id)
}

func (p *ProviderServiceStr) FindByNameLikeAndUserId(namePattern string, userId uuid.UUID, page, size int) []model.ProviderDto {
	return p.repository.FindByNameLikeAndUserId(namePattern, page, size, userId)
}

func (p *ProviderServiceStr) Update(id uuid.UUID, request model.UpdateProviderRequest) error {
	if !p.repository.Exists(id) {
		return notFoundError(id)
	}
	return p.repository.Update(id, request.ToEntity(id))
}

func (p *ProviderServiceStr) ExistsById(id uuid.UUID) bool {
	return p.repository.Exists(id)
}

func (p *ProviderServiceStr) Delete(id uuid.UUID) error {
	if !p.repository.Exists(id) {
		return notFoundError(id)
	}
	return p.repository.Delete(id)
}

func (p *ProviderServiceStr) Transactional(tx *gorm.DB) ProviderService {
	return &ProviderServiceStr{
		repository: p.repository.Transactional(tx),
	}
}

func notFoundError(id uuid.UUID) error {
	return int_errors.NewErrNotFound("provider with id %s not found", id)
}
