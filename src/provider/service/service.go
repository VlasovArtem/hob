package service

import (
	"common/dependency"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"provider/model"
	"provider/repository"
)

type ProviderServiceObject struct {
	repository repository.ProviderRepository
}

func NewProviderService(repository repository.ProviderRepository) ProviderService {
	return &ProviderServiceObject{repository}
}

func (p *ProviderServiceObject) Initialize(factory dependency.DependenciesFactory) {
	factory.Add(
		NewProviderService(factory.FindRequiredByObject(repository.ProviderRepositoryObject{}).(repository.ProviderRepository)),
	)
}

type ProviderService interface {
	Add(request model.CreateProviderRequest) (dto model.ProviderDto, err error)
	FindById(id uuid.UUID) (dto model.ProviderDto, err error)
	FindByNameLike(namePattern string, page int, size int) []model.ProviderDto
}

func (p *ProviderServiceObject) Add(request model.CreateProviderRequest) (dto model.ProviderDto, err error) {
	if p.repository.ExistsByName(request.Name) {
		return dto, errors.New(fmt.Sprintf("provider with name '%s' already exists", request.Name))
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
			return dto, errors.New(fmt.Sprintf("provider with id %s in not exists", id))
		}
		return dto, err
	} else {
		return provider.ToDto(), err
	}
}

func (p *ProviderServiceObject) FindByNameLike(namePattern string, page int, size int) []model.ProviderDto {
	return convert(p.repository.FindByNameLike(namePattern, page, size))
}

func convert(entities []model.Provider) (dtos []model.ProviderDto) {
	if len(entities) == 0 {
		return make([]model.ProviderDto, 0)
	}

	for _, entity := range entities {
		dtos = append(dtos, entity.ToDto())
	}
	return dtos
}
