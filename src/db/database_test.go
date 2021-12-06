package db

import (
	"fmt"
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
	config.DBName = "hob_test"
	i.database = NewDatabaseService(config)

	err := i.database.D().AutoMigrate(testEntity{})

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create new entity")
	}
}

func (i *DatabaseTestSuite) TearDownSuite() {
	err := i.database.D().Migrator().DropTable(testEntity{})

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot drop table")
	}
}

func TestDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}

func (i *DatabaseTestSuite) Test_Create() {
	entity := generateTestEntity()

	err := i.database.Create(&entity)

	assert.Nil(i.T(), err)

	i.database.ExistsById(testEntity{}, entity.Id)
}

func (i *DatabaseTestSuite) Test_Create_WithExistingId() {
	entity := generateTestEntity()

	err := i.database.Create(&entity)

	assert.Nil(i.T(), err)

	i.database.ExistsById(testEntity{}, entity.Id)

	err = i.database.Create(&entity)

	assert.NotNil(i.T(), err)
}

func (i *DatabaseTestSuite) Test_FindById() {
	entity := i.createTestEntity()

	receiver := testEntity{}
	err := i.database.FindById(&receiver, entity.Id)

	assert.Nil(i.T(), err)

	assert.Equal(i.T(), entity, receiver)
}

func (i *DatabaseTestSuite) Test_FindById_WithNotExistsId() {
	receiver := testEntity{}
	err := i.database.FindById(&receiver, uuid.New())

	assert.Equal(i.T(), gorm.ErrRecordNotFound, err)

	assert.Equal(i.T(), testEntity{}, receiver)
}

func (i *DatabaseTestSuite) Test_FindByIdModeled() {
	entity := i.createTestEntity()

	receiver := testEntityDto{}
	err := i.database.FindByIdModeled(testEntity{}, &receiver, entity.Id)

	assert.Nil(i.T(), err)

	assert.Equal(i.T(), testEntityDto{
		Id:          entity.Id,
		Name:        entity.Name,
		Description: entity.Description,
	}, receiver)
}

func (i *DatabaseTestSuite) Test_FindByIdModeled_WithNotExistsId() {
	receiver := testEntityDto{}
	err := i.database.FindByIdModeled(testEntity{}, &receiver, uuid.New())

	assert.Equal(i.T(), gorm.ErrRecordNotFound, err)

	assert.Equal(i.T(), testEntityDto{}, receiver)
}

func (i *DatabaseTestSuite) Test_ExistsById() {
	entity := i.createTestEntity()

	exists := i.database.ExistsById(testEntity{}, entity.Id)

	assert.True(i.T(), exists)
}

func (i *DatabaseTestSuite) Test_ExistsById_WithNotExistsId() {
	exists := i.database.ExistsById(testEntity{}, uuid.New())

	assert.False(i.T(), exists)
}

func (i *DatabaseTestSuite) Test_ExistsByQuery() {
	entity := i.createTestEntity()

	exists := i.database.ExistsByQuery(testEntity{}, "name = ?", entity.Name)

	assert.True(i.T(), exists)
}

func (i *DatabaseTestSuite) Test_ExistsByQuery_WithNotExists() {
	exists := i.database.ExistsByQuery(testEntity{}, "name = ?", "name match")

	assert.False(i.T(), exists)
}

func (i *DatabaseTestSuite) Test_DeleteById() {
	entity := i.createTestEntity()

	err := i.database.DeleteById(entity, entity.Id)

	assert.Nil(i.T(), err)
	assert.False(i.T(), i.database.ExistsById(testEntity{}, entity.Id))
}

func (i *DatabaseTestSuite) Test_DeleteById_WithNotExists() {
	err := i.database.DeleteById(testEntity{}, uuid.New())

	assert.Nil(i.T(), err)
}

func (i *DatabaseTestSuite) Test_UpdateById() {
	entity := i.createTestEntity()

	err := i.database.UpdateById(testEntity{}, entity.Id, testEntityUpdateDto{
		Id:          entity.Id,
		Name:        fmt.Sprintf("%s-new", entity.Name),
		Description: fmt.Sprintf("%s-new", entity.Description),
		Value:       entity.Value + 100,
	})

	assert.Nil(i.T(), err)

	newEntity := testEntity{}

	err = i.database.FindById(&newEntity, entity.Id)

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), testEntity{
		Id:          entity.Id,
		Name:        "Name-new",
		Description: "Description-new",
		Value:       200,
	}, newEntity)
}

func (i *DatabaseTestSuite) Test_UpdateById_WithOmitColumns() {
	entity := i.createTestEntity()

	err := i.database.UpdateById(testEntity{}, entity.Id, testEntityUpdateDto{
		Id:          entity.Id,
		Name:        fmt.Sprintf("%s-new", entity.Name),
		Description: fmt.Sprintf("%s-new", entity.Description),
		Value:       entity.Value + 100,
	}, "Value")

	assert.Nil(i.T(), err)

	newEntity := testEntity{}

	err = i.database.FindById(&newEntity, entity.Id)

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), testEntity{
		Id:          entity.Id,
		Name:        "Name-new",
		Description: "Description-new",
		Value:       100,
	}, newEntity)
}

func (i *DatabaseTestSuite) Test_UpdateById_WithNotExists() {
	err := i.database.UpdateById(testEntity{}, uuid.New(), testEntityUpdateDto{
		Id:          uuid.New(),
		Name:        "new name",
		Description: "new description",
		Value:       100,
	})

	assert.Nil(i.T(), err)
}

func (i *DatabaseTestSuite) Test_DM() {
	entity := generateTestEntity()
	entity.Name = uuid.New().String()

	err := i.database.Create(&entity)

	assert.Nil(i.T(), err)

	receiver := testEntityDto{}
	i.database.DM(testEntity{}).
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
	i.database.DM(testEntity{}).
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
