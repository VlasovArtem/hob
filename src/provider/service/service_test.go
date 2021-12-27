package service

import (
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/common/int-errors"
	"github.com/VlasovArtem/hob/src/provider/mocks"
	"github.com/VlasovArtem/hob/src/provider/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
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

	providerRepository.On("ExistsByNameAndUserId", request.Name, request.UserId).Return(false)
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

	providerRepository.On("ExistsByNameAndUserId", request.Name, request.UserId).Return(true)

	response, err := service.Add(request)

	assert.Equal(t, errors.New(fmt.Sprintf("provider with name '%s' for user already exists", request.Name)), err)
	assert.Equal(t, model.ProviderDto{}, response)
	providerRepository.AssertNotCalled(t, "Create")
}

func Test_Add_WithError(t *testing.T) {
	service := generateService()

	request := mocks.GenerateCreateProviderRequest()

	expectedError := errors.New("error")

	providerRepository.On("ExistsByNameAndUserId", request.Name, request.UserId).Return(false)
	providerRepository.On("Create", mock.Anything).Return(model.Provider{}, expectedError)

	response, err := service.Add(request)

	assert.Equal(t, expectedError, err)
	assert.Equal(t, model.ProviderDto{}, response)
}

func Test_FindById(t *testing.T) {
	service := generateService()

	entity := mocks.GenerateProvider(uuid.New())

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

	assert.Equal(t, int_errors.NewErrNotFound("provider with id %s not found", id), err)
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

func Test_FindByNameLikeAndUserIds(t *testing.T) {
	service := generateService()

	expected := mocks.GenerateProvider(uuid.New())

	providerRepository.On("FindByNameLikeAndUserId", expected.Name, 0, 25, expected.UserId).Return([]model.ProviderDto{expected.ToDto()})

	actual := service.FindByNameLikeAndUserId(expected.Name, expected.UserId, 0, 25)

	assert.Equal(t, []model.ProviderDto{expected.ToDto()}, actual)
}

func Test_FindByNameLikeAndUserIds_WithoutMatches(t *testing.T) {
	service := generateService()

	userId := uuid.New()

	providerRepository.On("FindByNameLikeAndUserId", "name", 0, 25, userId).Return([]model.ProviderDto{})

	actual := service.FindByNameLikeAndUserId("name", userId, 0, 25)

	assert.Equal(t, []model.ProviderDto{}, actual)
}

func Test_FindByUserId(t *testing.T) {
	service := generateService()

	entity := mocks.GenerateProvider(uuid.New())

	providerRepository.On("FindByUserId", entity.UserId).Return([]model.ProviderDto{entity.ToDto()}, nil)

	dto := service.FindByUserId(entity.UserId)

	assert.Equal(t, []model.ProviderDto{entity.ToDto()}, dto)
}

func Test_FindByUserId_WithEmptyResponse(t *testing.T) {
	service := generateService()

	userId := uuid.New()
	providerRepository.On("FindByUserId", userId).Return([]model.ProviderDto{}, nil)

	dto := service.FindByUserId(userId)

	assert.Equal(t, []model.ProviderDto{}, dto)
}

func Test_Update(t *testing.T) {
	houseService := generateService()

	request := mocks.GenerateUpdateProviderRequest()
	id := uuid.New()

	providerRepository.On("ExistsById", id).Return(true)
	providerRepository.On("Update", mock.Anything).Return(nil)

	assert.Nil(t, houseService.Update(id, request))

	providerRepository.AssertCalled(t, "Update", model.Provider{
		Id:      id,
		Name:    request.Name,
		Details: request.Details,
	})
}

func Test_Update_WithErrorFromDatabase(t *testing.T) {
	houseService := generateService()

	request := mocks.GenerateUpdateProviderRequest()
	id := uuid.New()

	providerRepository.On("ExistsById", id).Return(true)
	providerRepository.On("Update", mock.Anything).Return(errors.New("test"))

	err := houseService.Update(id, request)
	assert.Equal(t, errors.New("test"), err)
}

func Test_Update_WithNotExists(t *testing.T) {
	houseService := generateService()

	request := mocks.GenerateUpdateProviderRequest()
	id := uuid.New()

	providerRepository.On("ExistsById", id).Return(false)

	err := houseService.Update(id, request)
	assert.Equal(t, int_errors.NewErrNotFound("provider with id %s not found", id), err)

	providerRepository.AssertNotCalled(t, "Update", mock.Anything)
}

func Test_Delete(t *testing.T) {
	houseService := generateService()

	id := uuid.New()

	providerRepository.On("ExistsById", id).Return(true)
	providerRepository.On("Delete", id).Return(nil)

	err := houseService.Delete(id)
	assert.Nil(t, err)

	providerRepository.AssertCalled(t, "Delete", id)
}

func Test_Delete_WithError(t *testing.T) {
	houseService := generateService()

	id := uuid.New()

	providerRepository.On("ExistsById", id).Return(true)
	providerRepository.On("Delete", id).Return(errors.New("error"))

	err := houseService.Delete(id)
	assert.Equal(t, errors.New("error"), err)

	providerRepository.AssertCalled(t, "Delete", id)
}

func Test_Delete_WithNotExists(t *testing.T) {
	houseService := generateService()

	id := uuid.New()

	providerRepository.On("ExistsById", id).Return(false)

	err := houseService.Delete(id)
	assert.Equal(t, int_errors.NewErrNotFound("provider with id %s not found", id), err)

	providerRepository.AssertNotCalled(t, "Delete", mock.Anything)
}
