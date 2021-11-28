package service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"provider/mocks"
	"provider/model"
	"testing"
)

var providerRepository *mocks.ProviderRepository

func generateService() ProviderService {
	providerRepository = new(mocks.ProviderRepository)

	return NewProviderService(providerRepository)
}

func Test_Add(t *testing.T) {
	service := generateService()

	request := mocks.GenerateCreateProviderRequest()

	var expected model.Provider

	providerRepository.On("ExistsByName", request.Name).Return(false)
	providerRepository.On("Create", mock.Anything).Return(
		func(entity model.Provider) model.Provider {
			expected = entity
			return entity
		}, nil)

	response, err := service.Add(request)

	assert.Nil(t, err)
	assert.Equal(t, expected.ToDto(), response)
}

func Test_Add_WithExistingName(t *testing.T) {
	service := generateService()

	request := mocks.GenerateCreateProviderRequest()

	providerRepository.On("ExistsByName", request.Name).Return(true)

	response, err := service.Add(request)

	assert.Equal(t, errors.New(fmt.Sprintf("provider with name '%s' already exists", request.Name)), err)
	assert.Equal(t, model.ProviderDto{}, response)
	providerRepository.AssertNotCalled(t, "Create")
}

func Test_Add_WithError(t *testing.T) {
	service := generateService()

	request := mocks.GenerateCreateProviderRequest()

	expectedError := errors.New("error")

	providerRepository.On("ExistsByName", request.Name).Return(false)
	providerRepository.On("Create", mock.Anything).Return(model.Provider{}, expectedError)

	response, err := service.Add(request)

	assert.Equal(t, expectedError, err)
	assert.Equal(t, model.ProviderDto{}, response)
}

func Test_FindById(t *testing.T) {
	service := generateService()

	entity := mocks.GenerateProvider()

	providerRepository.On("FindById", entity.Id).Return(entity, nil)

	dto, err := service.FindById(entity.Id)

	assert.Nil(t, err)
	assert.Equal(t, entity.ToDto(), dto)
}

func Test_FindById_WithNotExistsRecord(t *testing.T) {
	service := generateService()

	id := uuid.New()

	providerRepository.On("FindById", id).Return(model.Provider{}, gorm.ErrRecordNotFound)

	response, err := service.FindById(id)

	assert.Equal(t, errors.New(fmt.Sprintf("provider with id %s in not exists", id)), err)
	assert.Equal(t, model.ProviderDto{}, response)
}

func Test_FindById_WithError(t *testing.T) {
	service := generateService()

	id := uuid.New()
	expectedError := errors.New("error")

	providerRepository.On("FindById", id).Return(model.Provider{}, expectedError)

	response, err := service.FindById(id)

	assert.Equal(t, expectedError, err)
	assert.Equal(t, model.ProviderDto{}, response)
}

func Test_FindByNameLike(t *testing.T) {
	service := generateService()

	expected := mocks.GenerateProvider()

	providerRepository.On("FindByNameLike", expected.Name, 0, 25).Return([]model.Provider{expected})

	actual := service.FindByNameLike(expected.Name, 0, 25)

	assert.Equal(t, []model.ProviderDto{expected.ToDto()}, actual)
}

func Test_FindByNameLike_WithoutMatches(t *testing.T) {
	service := generateService()

	providerRepository.On("FindByNameLike", "name", 0, 25).Return([]model.Provider{})

	actual := service.FindByNameLike("name", 0, 25)

	assert.Equal(t, []model.ProviderDto{}, actual)
}