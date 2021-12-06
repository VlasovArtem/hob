package repository

import (
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/provider/custom/mocks"
	"github.com/VlasovArtem/hob/src/provider/custom/model"
	"github.com/VlasovArtem/hob/src/test/testhelper/database"
	userMocks "github.com/VlasovArtem/hob/src/user/mocks"
	userModel "github.com/VlasovArtem/hob/src/user/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type CustomProviderRepositoryTestSuite struct {
	database.DBTestSuite
	repository  CustomProviderRepository
	createdUser userModel.User
}

func (p *CustomProviderRepositoryTestSuite) SetupSuite() {
	p.InitDBTestSuite()

	p.CreateRepository(
		func(service db.DatabaseService) {
			p.repository = NewCustomProviderRepository(service)
		},
	).
		AddMigrators(userModel.User{}, model.CustomProvider{})

	p.createdUser = userMocks.GenerateUser()
	p.CreateEntity(&p.createdUser)
}

func (p *CustomProviderRepositoryTestSuite) TearDownSuite() {
	p.TearDown()
}

func TestCustomProviderRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CustomProviderRepositoryTestSuite))
}

func (p *CustomProviderRepositoryTestSuite) Test_Create() {
	entity := mocks.GenerateCustomProvider(p.createdUser.Id)

	actual, err := p.repository.Create(entity)

	assert.Nil(p.T(), err)
	assert.Equal(p.T(), entity, actual)

	p.Delete(entity)
}

func (p *CustomProviderRepositoryTestSuite) Test_Create_WithSameNameButDifferentUsers() {
	first := mocks.GenerateCustomProvider(p.createdUser.Id)

	actual, err := p.repository.Create(first)

	assert.Nil(p.T(), err)
	assert.Equal(p.T(), first, actual)

	newUser := userMocks.GenerateUser()
	err = p.Database.Create(&newUser)

	assert.Nil(p.T(), err)

	second := mocks.GenerateCustomProvider(newUser.Id)
	second.Name = first.Name

	actual, err = p.repository.Create(second)

	assert.Nil(p.T(), err)
	assert.Equal(p.T(), second, actual)

	p.Delete(first)
	p.Delete(second)
}

func (p *CustomProviderRepositoryTestSuite) Test_Create_WithSameNameButSameUser() {
	first := mocks.GenerateCustomProvider(p.createdUser.Id)

	actual, err := p.repository.Create(first)

	assert.Nil(p.T(), err)
	assert.Equal(p.T(), first, actual)

	second := mocks.GenerateCustomProvider(p.createdUser.Id)
	second.Name = first.Name

	actual, err = p.repository.Create(second)

	assert.NotNil(p.T(), err)
	assert.Equal(p.T(), second, actual)

	p.Delete(first)
}

func (p *CustomProviderRepositoryTestSuite) Test_Creat_WithMissingUser() {
	entity := mocks.GenerateCustomProvider(uuid.New())

	actual, err := p.repository.Create(entity)

	assert.NotNil(p.T(), err)
	assert.Equal(p.T(), entity, actual)
}

func (p *CustomProviderRepositoryTestSuite) Test_FindById() {
	provider := p.CreateCustomProvider()

	actual, err := p.repository.FindById(provider.Id)

	assert.Nil(p.T(), err)
	assert.Equal(p.T(), provider, actual)
}

func (p *CustomProviderRepositoryTestSuite) Test_FindById_WithNotExistsRecord() {
	id := uuid.New()

	actual, err := p.repository.FindById(id)

	assert.Equal(p.T(), gorm.ErrRecordNotFound, err)
	assert.Equal(p.T(), model.CustomProvider{}, actual)
}

func (p *CustomProviderRepositoryTestSuite) Test_FindByUserId() {
	provider := p.CreateCustomProviderWithNewUser()

	actual := p.repository.FindByUserId(provider.UserId)

	assert.Equal(p.T(), []model.CustomProvider{provider}, actual)
}

func (p *CustomProviderRepositoryTestSuite) Test_FindByUserId_WithNotExistsRecord() {
	actual := p.repository.FindByUserId(uuid.New())

	assert.Equal(p.T(), []model.CustomProvider{}, actual)
}

func (p *CustomProviderRepositoryTestSuite) Test_ExistsById() {
	provider := p.CreateCustomProvider()

	assert.True(p.T(), p.repository.ExistsById(provider.Id))
}

func (p *CustomProviderRepositoryTestSuite) Test_ExistsById_WithNotExistsRecord() {
	id := uuid.New()

	assert.False(p.T(), p.repository.ExistsById(id))
}

func (p *CustomProviderRepositoryTestSuite) Test_ExistsByNameAndUserId() {
	provider := p.CreateCustomProvider()

	assert.True(p.T(), p.repository.ExistsByNameAndUserId(provider.Name, provider.UserId))
}

func (p *CustomProviderRepositoryTestSuite) Test_ExistsByNameAndUserId_WithNotMatchingName() {
	provider := p.CreateCustomProvider()

	assert.False(p.T(), p.repository.ExistsByNameAndUserId("not match", provider.UserId))
}

func (p *CustomProviderRepositoryTestSuite) Test_ExistsByNameAndUserId_WithNotMatchingUserId() {
	provider := p.CreateCustomProvider()

	assert.False(p.T(), p.repository.ExistsByNameAndUserId(provider.Name, uuid.New()))
}

func (p *CustomProviderRepositoryTestSuite) CreateCustomProvider() model.CustomProvider {
	provider := mocks.GenerateCustomProvider(p.createdUser.Id)

	p.CreateEntity(provider)

	return provider
}

func (p *CustomProviderRepositoryTestSuite) CreateCustomProviderWithNewUser() model.CustomProvider {
	createdUser := userMocks.GenerateUser()
	p.CreateEntity(&createdUser)

	provider := mocks.GenerateCustomProvider(createdUser.Id)

	p.CreateEntity(provider)

	return provider
}
