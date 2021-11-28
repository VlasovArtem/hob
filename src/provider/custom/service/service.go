package service

import (
	"common/dependency"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"provider/custom/model"
	"provider/custom/repository"
	userRepository "user/repository"
)

type CustomProviderServiceObject struct {
	repository     repository.CustomProviderRepository
	userRepository userRepository.UserRepository
}

func NewCustomProviderService(repository repository.CustomProviderRepository, userRepository userRepository.UserRepository) CustomProviderService {
	return &CustomProviderServiceObject{repository, userRepository}
}

func (c *CustomProviderServiceObject) Initialize(factory dependency.DependenciesFactory) {
	factory.Add(
		NewCustomProviderService(
			factory.FindRequiredByObject(repository.CustomProviderRepositoryObject{}).(repository.CustomProviderRepository),
			factory.FindRequiredByObject(userRepository.UserRepositoryObject{}).(userRepository.UserRepository),
		),
	)
}

type CustomProviderService interface {
	Add(request model.CreateCustomProviderRequest) (model.CustomProviderDto, error)
	FindById(id uuid.UUID) (model.CustomProviderDto, error)
	FindByUserId(id uuid.UUID) []model.CustomProviderDto
}

func (c *CustomProviderServiceObject) Add(request model.CreateCustomProviderRequest) (dto model.CustomProviderDto, err error) {
	if !c.userRepository.ExistsById(request.UserId) {
		return dto, errors.New(fmt.Sprintf("user with %s not exists", request.UserId))
	} else if c.repository.ExistsByNameAndUserId(request.Name, request.UserId) {
		return dto, errors.New(fmt.Sprintf("user already have provider with name %s", request.Name))
	}

	entity := request.ToEntity()

	if created, err := c.repository.Create(entity); err != nil {
		return dto, err
	} else {
		return created.ToDto(), err
	}
}

func (c *CustomProviderServiceObject) FindById(id uuid.UUID) (dto model.CustomProviderDto, err error) {
	if provider, err := c.repository.FindById(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto, errors.New(fmt.Sprintf("custom provider with id %s in not exists", id))
		}
		return dto, err
	} else {
		return provider.ToDto(), err
	}
}

func (c *CustomProviderServiceObject) FindByUserId(id uuid.UUID) (response []model.CustomProviderDto) {
	for _, provider := range c.repository.FindByUserId(id) {
		response = append(response, provider.ToDto())
	}

	if response == nil {
		return make([]model.CustomProviderDto, 0)
	}

	return response
}
