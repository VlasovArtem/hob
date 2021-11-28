package repository

import (
	"db"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"provider/mocks"
	"provider/model"
	"test/testhelper/database"
	"testing"
)

type ProviderRepositoryTestSuite struct {
	suite.Suite
	database   db.DatabaseService
	repository ProviderRepository
}

func (p *ProviderRepositoryTestSuite) SetupSuite() {
	config := db.NewDefaultDatabaseConfiguration()
	config.DBName = "hob_test"
	p.database = db.NewDatabaseService(config)
	p.repository = NewProviderRepository(p.database)
	database.CreateTable(p.database.D(), model.Provider{})
}

func (p *ProviderRepositoryTestSuite) TearDownSuite() {
	database.DropTable(p.database.D(), model.Provider{})
}

func TestProviderRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ProviderRepositoryTestSuite))
}

func (p *ProviderRepositoryTestSuite) Test_Create() {
	provider := mocks.GenerateProvider()

	create, err := p.repository.Create(provider)

	assert.Nil(p.T(), err)
	assert.Equal(p.T(), provider, create)
}

func (p *ProviderRepositoryTestSuite) Test_CreateWithMatchingName() {
	expected := mocks.GenerateProvider()

	actual, err := p.repository.Create(expected)

	assert.Nil(p.T(), err)
	assert.Equal(p.T(), expected, actual)

	actual, err = p.repository.Create(expected)

	assert.Equal(p.T(), expected, actual)
	assert.NotNil(p.T(), err)
}

func (p *ProviderRepositoryTestSuite) Test_FindById() {
	provider := p.CreateProvider()

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

func (p *ProviderRepositoryTestSuite) Test_FindByNameLike() {
	var expectedProviders []model.Provider

	for i := 0; i < 10; i++ {
		provider := mocks.GenerateProvider()
		provider.Name = fmt.Sprintf("findByName-%d", i)

		created, _ := p.repository.Create(provider)

		expectedProviders = append(expectedProviders, created)
	}

	firstPage := p.repository.FindByNameLike("findByName", 0, 5)

	assert.Equal(p.T(), expectedProviders[:5], firstPage)

	secondPage := p.repository.FindByNameLike("findByName", 1, 5)

	assert.Equal(p.T(), expectedProviders[5:], secondPage)

	database.RecreateTable(p.database.D(), model.Provider{})
}

func (p *ProviderRepositoryTestSuite) Test_FindByNameLike_WithOutMatch() {
	providers := p.repository.FindByNameLike("invalid", 0, 100)

	assert.Equal(p.T(), []model.Provider{}, providers)
}

func (p *ProviderRepositoryTestSuite) Test_FindByNameLike_WithEmptyString() {
	database.RecreateTable(p.database.D(), model.Provider{})

	var expectedProviders []model.Provider

	for i := 0; i < 10; i++ {
		provider := mocks.GenerateProvider()
		provider.Name = fmt.Sprintf("%d", i)

		created, _ := p.repository.Create(provider)

		expectedProviders = append(expectedProviders, created)
	}

	page := p.repository.FindByNameLike("", 0, len(expectedProviders))

	assert.ElementsMatch(p.T(), expectedProviders, page)

	database.RecreateTable(p.database.D(), model.Provider{})
}

func (p *ProviderRepositoryTestSuite) Test_ExistsById() {
	provider := p.CreateProvider()

	assert.True(p.T(), p.repository.ExistsById(provider.Id))
}

func (p *ProviderRepositoryTestSuite) Test_ExistsById_WithNotExistsRecord() {
	id := uuid.New()

	assert.False(p.T(), p.repository.ExistsById(id))
}

func (p *ProviderRepositoryTestSuite) Test_ExistsByName() {
	provider := p.CreateProvider()

	assert.True(p.T(), p.repository.ExistsByName(provider.Name))
}

func (p *ProviderRepositoryTestSuite) Test_ExistsByName_WithNotExistsRecord() {
	assert.False(p.T(), p.repository.ExistsByName("not exists"))
}

func (p *ProviderRepositoryTestSuite) CreateProvider() model.Provider {
	provider := mocks.GenerateProvider()

	saved, err := p.repository.Create(provider)

	assert.Nil(p.T(), err)

	return saved
}
