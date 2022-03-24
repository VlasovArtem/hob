package service

import (
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/common/int-errors"
	"github.com/VlasovArtem/hob/src/provider/mocks"
	"github.com/VlasovArtem/hob/src/provider/model"
	"github.com/VlasovArtem/hob/src/test/testhelper"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type ProviderServiceTestSuite struct {
	testhelper.MockTestSuite[ProviderService]
	providerRepository *mocks.ProviderRepository
}

func TestProviderServiceTestSuite(t *testing.T) {
	ts := &ProviderServiceTestSuite{}
	ts.TestObjectGenerator = func() ProviderService {
		ts.providerRepository = new(mocks.ProviderRepository)

		return NewProviderService(ts.providerRepository)
	}

	suite.Run(t, ts)
}

func (p *ProviderServiceTestSuite) Test_Add() {
	request := mocks.GenerateCreateProviderRequest()

	var expected model.Provider

	p.providerRepository.On("ExistsByNameAndUserId", request.Name, request.UserId).Return(false)
	p.providerRepository.On("Create", mock.Anything).Return(
		func(entity model.Provider) model.Provider {
			expected = entity
			return entity
		}, nil)

	response, err := p.TestO.Add(request)

	assert.Nil(p.T(), err)
	assert.Equal(p.T(), expected.ToDto(), response)
}

func (p *ProviderServiceTestSuite) Test_Add_WithExistingName() {
	request := mocks.GenerateCreateProviderRequest()

	p.providerRepository.On("ExistsByNameAndUserId", request.Name, request.UserId).Return(true)

	response, err := p.TestO.Add(request)

	assert.Equal(p.T(), errors.New(fmt.Sprintf("provider with name '%s' for user already exists", request.Name)), err)
	assert.Equal(p.T(), model.ProviderDto{}, response)
	p.providerRepository.AssertNotCalled(p.T(), "Create")
}

func (p *ProviderServiceTestSuite) Test_Add_WithError() {
	request := mocks.GenerateCreateProviderRequest()

	expectedError := errors.New("error")

	p.providerRepository.On("ExistsByNameAndUserId", request.Name, request.UserId).Return(false)
	p.providerRepository.On("Create", mock.Anything).Return(model.Provider{}, expectedError)

	response, err := p.TestO.Add(request)

	assert.Equal(p.T(), expectedError, err)
	assert.Equal(p.T(), model.ProviderDto{}, response)
}

func (p *ProviderServiceTestSuite) Test_FindById() {
	entity := mocks.GenerateProvider(uuid.New())

	p.providerRepository.On("FindById", entity.Id).Return(entity, nil)

	dto, err := p.TestO.FindById(entity.Id)

	assert.Nil(p.T(), err)
	assert.Equal(p.T(), entity.ToDto(), dto)
}

func (p *ProviderServiceTestSuite) Test_FindById_WithNotExistsRecord() {
	id := uuid.New()

	p.providerRepository.On("FindById", id).Return(model.Provider{}, gorm.ErrRecordNotFound)

	response, err := p.TestO.FindById(id)

	assert.Equal(p.T(), int_errors.NewErrNotFound("provider with id %s not found", id), err)
	assert.Equal(p.T(), model.ProviderDto{}, response)
}

func (p *ProviderServiceTestSuite) Test_FindById_WithError() {
	id := uuid.New()
	expectedError := errors.New("error")

	p.providerRepository.On("FindById", id).Return(model.Provider{}, expectedError)

	response, err := p.TestO.FindById(id)

	assert.Equal(p.T(), expectedError, err)
	assert.Equal(p.T(), model.ProviderDto{}, response)
}

func (p *ProviderServiceTestSuite) Test_FindByNameLikeAndUserIds() {
	expected := mocks.GenerateProvider(uuid.New())

	p.providerRepository.On("FindByNameLikeAndUserId", expected.Name, 0, 25, expected.UserId).Return([]model.ProviderDto{expected.ToDto()})

	actual := p.TestO.FindByNameLikeAndUserId(expected.Name, expected.UserId, 0, 25)

	assert.Equal(p.T(), []model.ProviderDto{expected.ToDto()}, actual)
}

func (p *ProviderServiceTestSuite) Test_FindByNameLikeAndUserIds_WithoutMatches() {
	userId := uuid.New()

	p.providerRepository.On("FindByNameLikeAndUserId", "name", 0, 25, userId).Return([]model.ProviderDto{})

	actual := p.TestO.FindByNameLikeAndUserId("name", userId, 0, 25)

	assert.Equal(p.T(), []model.ProviderDto{}, actual)
}

func (p *ProviderServiceTestSuite) Test_FindByUserId() {
	entity := mocks.GenerateProvider(uuid.New())

	p.providerRepository.On("FindByUserId", entity.UserId).Return([]model.ProviderDto{entity.ToDto()}, nil)

	dto := p.TestO.FindByUserId(entity.UserId)

	assert.Equal(p.T(), []model.ProviderDto{entity.ToDto()}, dto)
}

func (p *ProviderServiceTestSuite) Test_FindByUserId_WithEmptyResponse() {
	userId := uuid.New()
	p.providerRepository.On("FindByUserId", userId).Return([]model.ProviderDto{}, nil)

	dto := p.TestO.FindByUserId(userId)

	assert.Equal(p.T(), []model.ProviderDto{}, dto)
}

func (p *ProviderServiceTestSuite) Test_Update() {
	request := mocks.GenerateUpdateProviderRequest()
	id := uuid.New()

	p.providerRepository.On("ExistsById", id).Return(true)
	p.providerRepository.On("Update", mock.Anything).Return(nil)

	assert.Nil(p.T(), p.TestO.Update(id, request))

	p.providerRepository.AssertCalled(p.T(), "Update", model.Provider{
		Id:      id,
		Name:    request.Name,
		Details: request.Details,
	})
}

func (p *ProviderServiceTestSuite) Test_Update_WithErrorFromDatabase() {
	request := mocks.GenerateUpdateProviderRequest()
	id := uuid.New()

	p.providerRepository.On("ExistsById", id).Return(true)
	p.providerRepository.On("Update", mock.Anything).Return(errors.New("test"))

	err := p.TestO.Update(id, request)
	assert.Equal(p.T(), errors.New("test"), err)
}

func (p *ProviderServiceTestSuite) Test_Update_WithNotExists() {
	request := mocks.GenerateUpdateProviderRequest()
	id := uuid.New()

	p.providerRepository.On("ExistsById", id).Return(false)

	err := p.TestO.Update(id, request)
	assert.Equal(p.T(), int_errors.NewErrNotFound("provider with id %s not found", id), err)

	p.providerRepository.AssertNotCalled(p.T(), "Update", mock.Anything)
}

func (p *ProviderServiceTestSuite) Test_Delete() {
	id := uuid.New()

	p.providerRepository.On("ExistsById", id).Return(true)
	p.providerRepository.On("Delete", id).Return(nil)

	err := p.TestO.Delete(id)
	assert.Nil(p.T(), err)

	p.providerRepository.AssertCalled(p.T(), "Delete", id)
}

func (p *ProviderServiceTestSuite) Test_Delete_WithError() {
	id := uuid.New()

	p.providerRepository.On("ExistsById", id).Return(true)
	p.providerRepository.On("Delete", id).Return(errors.New("error"))

	err := p.TestO.Delete(id)
	assert.Equal(p.T(), errors.New("error"), err)

	p.providerRepository.AssertCalled(p.T(), "Delete", id)
}

func (p *ProviderServiceTestSuite) Test_Delete_WithNotExists() {
	id := uuid.New()

	p.providerRepository.On("ExistsById", id).Return(false)

	err := p.TestO.Delete(id)
	assert.Equal(p.T(), int_errors.NewErrNotFound("provider with id %s not found", id), err)

	p.providerRepository.AssertNotCalled(p.T(), "Delete", mock.Anything)
}
