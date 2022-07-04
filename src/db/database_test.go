package db

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type testEntity struct {
	Id          uuid.UUID `gorm:"primarykey"`
	Name        string
	Description string
	Value       int
}

type testEntityDto struct {
	Id          uuid.UUID
	Name        string
	Description string
}

type testEntityUpdateDto struct {
	Id          uuid.UUID
	Name        string
	Description string
	Value       int
}

type DatabaseTestSuite struct {
	suite.Suite
	database DatabaseService
}

func (i *DatabaseTestSuite) SetupSuite() {
	config := NewDefaultDatabaseConfiguration()
	config.Port = 5444
	config.DBName = "hob_test"
	i.database = NewDatabaseService(config)

	err := i.database.DB().AutoMigrate(testEntity{})

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create new entity")
	}
}

func (i *DatabaseTestSuite) TearDownSuite() {
	err := i.database.DB().Migrator().DropTable(testEntity{})

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot drop table")
	}
}

func (i *DatabaseTestSuite) TearDownTest() {
	err := i.database.DB().Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(testEntity{}).Error

	if err != nil {
		log.Err(err).Msg("Cannot truncate table")
	}
}

func TestDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}

func (i *DatabaseTestSuite) Test_Create() {
	entity := generateTestEntity()

	err := i.database.Create(&entity)

	assert.Nil(i.T(), err)
}

func (i *DatabaseTestSuite) Test_Create_WithExistingId() {
	entity := generateTestEntity()

	err := i.database.Create(&entity)

	assert.Nil(i.T(), err)

	err = i.database.Create(&entity)

	assert.NotNil(i.T(), err)
}

func (i *DatabaseTestSuite) Test_DM() {
	entity := generateTestEntity()
	entity.Name = uuid.New().String()

	err := i.database.Create(&entity)

	assert.Nil(i.T(), err)

	receiver := testEntityDto{}
	i.database.DBModeled(testEntity{}).
		Where("name = ?", entity.Name).
		Find(&receiver)

	assert.Equal(i.T(), testEntityDto{
		Id:          entity.Id,
		Name:        entity.Name,
		Description: entity.Description,
	}, receiver)
}

func (i *DatabaseTestSuite) Test_DM_WithNotMatchedRecord() {
	receiver := testEntityDto{}
	i.database.DBModeled(testEntity{}).
		Where("name = ?", "not match").
		Find(&receiver)

	assert.Equal(i.T(), testEntityDto{}, receiver)
}

func generateTestEntity() testEntity {
	return testEntity{
		Id:          uuid.New(),
		Name:        "Name",
		Description: "Description",
		Value:       100,
	}
}

func (i *DatabaseTestSuite) createTestEntity() testEntity {
	entity := generateTestEntity()

	err := i.database.Create(&entity)

	assert.Nil(i.T(), err)

	return entity
}
