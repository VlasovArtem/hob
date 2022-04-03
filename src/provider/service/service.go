package service

import (
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/common/dependency"
	int_errors "github.com/VlasovArtem/hob/src/common/int-errors"
	"github.com/VlasovArtem/hob/src/provider/model"
	"github.com/VlasovArtem/hob/src/provider/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProviderServiceObject struct {
	repository repository.ProviderRepository
}

func NewProviderService(repository repository.ProviderRepository) ProviderService {
	return &ProviderServiceObject{repository}
}

func (p *ProviderServiceObject) Initialize(factory dependency.DependenciesProvider) any {
	return NewProviderService(dependency.FindRequiredDependency[repository.ProviderRepositoryObject, repository.ProviderRepository](factory))
}

type ProviderService interface {
	Add(request model.CreateProviderRequest) (dto model.ProviderDto, err error)
	ExistsById(id uuid.UUID) bool
	Update(id uuid.UUID, request model.UpdateProviderRequest) error
	Delete(id uuid.UUID) error
	FindById(id uuid.UUID) (dto model.ProviderDto, err error)
	FindByUserId(id uuid.UUID) []model.ProviderDto
	FindByNameLikeAndUserId(namePattern string, userId uuid.UUID, page, size int) []model.ProviderDto
}

func (p *ProviderServiceObject) Add(request model.CreateProviderRequest) (dto model.ProviderDto, err error) {
	if p.repository.ExistsByNameAndUserId(request.Name, request.UserId) {
		return dto, errors.New(fmt.Sprintf("provider with name '%s' for user already exists", request.Name))
	}

	if entity, err := p.repository.Create(request.ToEntity()); err != nil {
		return dto, err
	} else {
		return entity.ToDto(), err
	}
}

func (p *ProviderServiceObject) FindById(id uuid.UUID) (dto model.ProviderDto, err error) {
	if provider, err := p.repository.FindById(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto, notFoundError(id)
		}
		return dto, err
	} else {
		return provider.ToDto(), err
	}
}

func (p *ProviderServiceObject) FindByUserId(id uuid.UUID) (response []model.ProviderDto) {
	return p.repository.FindByUserId(id)
}

func (p *ProviderServiceObject) FindByNameLikeAndUserId(namePattern string, userId uuid.UUID, page, size int) []model.ProviderDto {
	return p.repository.FindByNameLikeAndUserId(namePattern, page, size, userId)
}

func (p *ProviderServiceObject) Update(id uuid.UUID, request model.UpdateProviderRequest) error {
	if !p.repository.ExistsById(id) {
		return notFoundError(id)
	}
	return p.repository.Update(request.ToEntity(id))
}

func (p *ProviderServiceObject) ExistsById(id uuid.UUID) bool {
	return p.repository.ExistsById(id)
}

func (p *ProviderServiceObject) Delete(id uuid.UUID) error {
	if !p.repository.ExistsById(id) {
		return notFoundError(id)
	}
	return p.repository.Delete(id)
}

func notFoundError(id uuid.UUID) error {
	return int_errors.NewErrNotFound("provider with id %s not found", id)
}
