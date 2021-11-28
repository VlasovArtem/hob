package service

import (
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/provider/custom/mocks"
	"github.com/VlasovArtem/hob/src/provider/custom/model"
	userMocks "github.com/VlasovArtem/hob/src/user/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"testing"
)

var (
	providerRepository *mocks.CustomProviderRepository
	userRepositoryMock *userMocks.UserRepository
)

func generateService() CustomProviderService {
	providerRepository = new(mocks.CustomProviderRepository)
	userRepositoryMock = new(userMocks.UserRepository)

	return NewCustomProviderService(providerRepository, userRepositoryMock)
}

func Test_Add(t *testing.T) {
	service := generateService()

	request := mocks.GenerateCustomProviderRequest()

	var expected model.CustomProvider

	userRepositoryMock.On("ExistsById", request.UserId).Return(true)
	providerRepository.On("ExistsByNameAndUserId", request.Name, request.UserId).Return(false)
	providerRepository.On("Create", mock.Anything).Return(
		func(entity model.CustomProvider) model.CustomProvider {
			expected = entity
			return entity
		}, nil)

	response, err := service.Add(request)

	assert.Nil(t, err)
	assert.Equal(t, expected.ToDto(), response)
}

func Test_Add_WithNotExistingUser(t *testing.T) {
	service := generateService()

	request := mocks.GenerateCustomProviderRequest()

	userRepositoryMock.On("ExistsById", request.UserId).Return(false)

	response, err := service.Add(request)

	assert.Equal(t, errors.New(fmt.Sprintf("user with %s not exists", request.UserId)), err)
	assert.Equal(t, model.CustomProviderDto{}, response)
	providerRepository.AssertNotCalled(t, "ExistsByNameAndUserId")
	providerRepository.AssertNotCalled(t, "Create")
}

func Test_Add_WithExistingNameByUser(t *testing.T) {
	service := generateService()

	request := mocks.GenerateCustomProviderRequest()

	userRepositoryMock.On("ExistsById", request.UserId).Return(true)
	providerRepository.On("ExistsByNameAndUserId", request.Name, request.UserId).Return(true)

	response, err := service.Add(request)

	assert.Equal(t, errors.New(fmt.Sprintf("user already have provider with name %s", request.Name)), err)
	assert.Equal(t, model.CustomProviderDto{}, response)
	providerRepository.AssertNotCalled(t, "Create")
}

func Test_Add_WithError(t *testing.T) {
	service := generateService()

	request := mocks.GenerateCustomProviderRequest()

	expectedError := errors.New("error")

	userRepositoryMock.On("ExistsById", request.UserId).Return(true)
	providerRepository.On("ExistsByNameAndUserId", request.Name, request.UserId).Return(false)
	providerRepository.On("Create", mock.Anything).Return(model.CustomProvider{}, expectedError)

	response, err := service.Add(request)

	assert.Equal(t, expectedError, err)
	assert.Equal(t, model.CustomProviderDto{}, response)
}

func Test_FindById(t *testing.T) {
	service := generateService()

	entity := mocks.GenerateCustomProvider(uuid.New())

	providerRepository.On("FindById", entity.Id).Return(entity, nil)

	dto, err := service.FindById(entity.Id)

	assert.Nil(t, err)
	assert.Equal(t, entity.ToDto(), dto)
}

func Test_FindById_WithNotExistsRecord(t *testing.T) {
	service := generateService()

	id := uuid.New()

	providerRepository.On("FindById", id).Return(model.CustomProvider{}, gorm.ErrRecordNotFound)

	response, err := service.FindById(id)

	assert.Equal(t, errors.New(fmt.Sprintf("custom provider with id %s in not exists", id)), err)
	assert.Equal(t, model.CustomProviderDto{}, response)
}

func Test_FindById_WithError(t *testing.T) {
	service := generateService()

	id := uuid.New()
	expectedError := errors.New("error")

	providerRepository.On("FindById", id).Return(model.CustomProvider{}, expectedError)

	response, err := service.FindById(id)

	assert.Equal(t, expectedError, err)
	assert.Equal(t, model.CustomProviderDto{}, response)
}

func Test_FindByUserId(t *testing.T) {
	service := generateService()

	entity := mocks.GenerateCustomProvider(uuid.New())

	providerRepository.On("FindByUserId", entity.UserId).Return([]model.CustomProvider{entity}, nil)

	dto := service.FindByUserId(entity.UserId)

	assert.Equal(t, []model.CustomProviderDto{entity.ToDto()}, dto)
}

func Test_FindByUserId_WithEmptyResponse(t *testing.T) {
	service := generateService()

	userId := uuid.New()
	providerRepository.On("FindByUserId", userId).Return([]model.CustomProvider{}, nil)

	dto := service.FindByUserId(userId)

	assert.Equal(t, []model.CustomProviderDto{}, dto)
}
