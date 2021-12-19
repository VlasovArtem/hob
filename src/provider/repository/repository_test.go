package repository

import (
	"fmt"
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/provider/mocks"
	"github.com/VlasovArtem/hob/src/provider/model"
	"github.com/VlasovArtem/hob/src/test/testhelper/database"
	userMocks "github.com/VlasovArtem/hob/src/user/mocks"
	userModel "github.com/VlasovArtem/hob/src/user/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type ProviderRepositoryTestSuite struct {
	database.DBTestSuite
	repository  ProviderRepository
	createdUser userModel.User
}

func (p *ProviderRepositoryTestSuite) SetupSuite() {
	p.InitDBTestSuite()

	p.CreateRepository(
		func(service db.DatabaseService) {
			p.repository = NewProviderRepository(service)
		},
	).
		AddMigrators(userModel.User{}, model.Provider{})

	p.createdUser = userMocks.GenerateUser()
	p.CreateConstantEntity(&p.createdUser)
}

func TestProviderRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ProviderRepositoryTestSuite))
}

func (p *ProviderRepositoryTestSuite) Test_Create() {
	entity := mocks.GenerateProvider(p.createdUser.Id)

	actual, err := p.repository.Create(entity)

	assert.Nil(p.T(), err)
	assert.Equal(p.T(), entity, actual)

	p.Delete(entity)
}

func (p *ProviderRepositoryTestSuite) Test_Create_WithDefaultUser() {
	defaultProvider := mocks.GenerateProvider(uuid.UUID{})
	provider := mocks.GenerateProvider(p.createdUser.Id)

	actual, err := p.repository.Create(defaultProvider)
	actual1, err := p.repository.Create(provider)

	assert.Nil(p.T(), err)
	assert.Equal(p.T(), defaultProvider, actual)

	assert.Nil(p.T(), err)
	assert.Equal(p.T(), provider, actual1)

	providers := p.repository.FindByUserId(p.createdUser.Id)

	assert.EqualValues(p.T(), providers, []model.ProviderDto{defaultProvider.ToDto(), provider.ToDto()})

	p.Delete(defaultProvider)
	p.Delete(provider)
}

func (p *ProviderRepositoryTestSuite) Test_Create_WithSameNameButDifferentUsers() {
	first := mocks.GenerateProvider(p.createdUser.Id)

	actual, err := p.repository.Create(first)

	assert.Nil(p.T(), err)
	assert.Equal(p.T(), first, actual)

	newUser := userMocks.GenerateUser()
	err = p.Database.Create(&newUser)

	assert.Nil(p.T(), err)

	second := mocks.GenerateProvider(newUser.Id)
	second.Name = first.Name

	actual, err = p.repository.Create(second)

	assert.Nil(p.T(), err)
	assert.Equal(p.T(), second, actual)

	p.Delete(first)
	p.Delete(second)
}

func (p *ProviderRepositoryTestSuite) Test_Create_WithSameNameButSameUser() {
	first := mocks.GenerateProvider(p.createdUser.Id)

	actual, err := p.repository.Create(first)

	assert.Nil(p.T(), err)
	assert.Equal(p.T(), first, actual)

	second := mocks.GenerateProvider(p.createdUser.Id)
	second.Name = first.Name

	actual, err = p.repository.Create(second)

	assert.NotNil(p.T(), err)
	assert.Equal(p.T(), second, actual)

	p.Delete(first)
}

func (p *ProviderRepositoryTestSuite) Test_Creat_WithMissingUser() {
	entity := mocks.GenerateProvider(uuid.New())

	actual, err := p.repository.Create(entity)

	assert.NotNil(p.T(), err)
	assert.Equal(p.T(), entity, actual)
}

func (p *ProviderRepositoryTestSuite) Test_FindById() {
	provider := p.createCustomProvider()

	actual, err := p.repository.FindById(provider.Id)

	assert.Nil(p.T(), err)
	assert.Equal(p.T(), provider, actual)
}

func (p *ProviderRepositoryTestSuite) Test_FindById_WithNotExistsRecord() {
	id := uuid.New()

	actual, err := p.repository.FindById(id)

	assert.Equal(p.T(), gorm.ErrRecordNotFound, err)
	assert.Equal(p.T(), model.Provider{}, actual)
}

func (p *ProviderRepositoryTestSuite) Test_Delete() {
	provider := p.createCustomProvider()

	err := p.repository.Delete(provider.Id)

	assert.Nil(p.T(), err)
	assert.False(p.T(), p.repository.ExistsById(provider.Id))
}

func (p *ProviderRepositoryTestSuite) Test_FindByUserId() {
	provider := p.createCustomProviderWithNewUser()

	actual := p.repository.FindByUserId(provider.UserId)

	assert.Equal(p.T(), []model.ProviderDto{provider.ToDto()}, actual)
}

func (p *ProviderRepositoryTestSuite) Test_FindByUserId_WithNotExistsRecord() {
	actual := p.repository.FindByUserId(uuid.New())

	assert.Equal(p.T(), []model.ProviderDto{}, actual)
}

func (p *ProviderRepositoryTestSuite) Test_FindByNameLikeAndUserIds() {
	provider := p.createCustomProviderWithNewUser()

	actual := p.repository.FindByNameLikeAndUserId("Provider", 0, 10, provider.UserId)

	assert.Equal(p.T(), []model.ProviderDto{provider.ToDto()}, actual)
}

func (p *ProviderRepositoryTestSuite) Test_FindByNameLikeAndUserIds_WithNotMatchingName() {
	provider := p.createCustomProviderWithNewUser()

	actual := p.repository.FindByNameLikeAndUserId("invalid", 0, 10, provider.UserId)

	assert.Equal(p.T(), []model.ProviderDto{}, actual)
}

func (p *ProviderRepositoryTestSuite) Test_FindByNameLikeAndUserIds_WithNotMatchingUserId() {
	actual := p.repository.FindByNameLikeAndUserId("Provider", 0, 10, uuid.New())

	assert.Equal(p.T(), []model.ProviderDto{}, actual)
}

func (p *ProviderRepositoryTestSuite) Test_ExistsById() {
	provider := p.createCustomProvider()

	assert.True(p.T(), p.repository.ExistsById(provider.Id))
}

func (p *ProviderRepositoryTestSuite) Test_ExistsById_WithNotExistsRecord() {
	id := uuid.New()

	assert.False(p.T(), p.repository.ExistsById(id))
}

func (p *ProviderRepositoryTestSuite) Test_ExistsByNameAndUserId() {
	provider := p.createCustomProvider()

	assert.True(p.T(), p.repository.ExistsByNameAndUserId(provider.Name, provider.UserId))
}

func (p *ProviderRepositoryTestSuite) Test_ExistsByNameAndUserId_WithNotMatchingName() {
	provider := p.createCustomProvider()

	assert.False(p.T(), p.repository.ExistsByNameAndUserId("not match", provider.UserId))
}

func (p *ProviderRepositoryTestSuite) Test_ExistsByNameAndUserId_WithNotMatchingUserId() {
	provider := p.createCustomProvider()

	assert.False(p.T(), p.repository.ExistsByNameAndUserId(provider.Name, uuid.New()))
}

func (p *ProviderRepositoryTestSuite) Test_Update() {
	provider := p.createCustomProvider()

	updated := model.Provider{
		Id:      provider.Id,
		Name:    fmt.Sprintf("%s-new", provider.Name),
		Details: fmt.Sprintf("%s-new", provider.Details),
	}

	err := p.repository.Update(updated)

	assert.Nil(p.T(), err)

	response, err := p.repository.FindById(provider.Id)
	assert.Nil(p.T(), err)
	assert.Equal(p.T(), model.Provider{
		Id:      provider.Id,
		Name:    fmt.Sprintf("%s-new", provider.Name),
		Details: "Details-new",
		UserId:  provider.UserId,
		User:    provider.User,
	}, response)
}

func (p *ProviderRepositoryTestSuite) Test_Update_WithMatchingName() {
	first := p.createCustomProvider()
	provider := p.createCustomProvider()

	updated := model.Provider{
		Id:      provider.Id,
		Name:    first.Name,
		Details: fmt.Sprintf("%s-new", provider.Details),
	}

	err := p.repository.Update(updated)

	assert.NotNil(p.T(), err)
}

func (p *ProviderRepositoryTestSuite) Test_Update_WithMissingId() {
	assert.Nil(p.T(), p.repository.Update(model.Provider{Id: uuid.New()}))
}

func (p *ProviderRepositoryTestSuite) createCustomProvider() model.Provider {
	provider := mocks.GenerateProvider(p.createdUser.Id)

	p.CreateEntity(provider)

	return provider
}

func (p *ProviderRepositoryTestSuite) createCustomProviderWithNewUser() model.Provider {
	createdUser := userMocks.GenerateUser()
	p.CreateEntity(&createdUser)

	provider := mocks.GenerateProvider(createdUser.Id)

	p.CreateEntity(provider)

	return provider
}
