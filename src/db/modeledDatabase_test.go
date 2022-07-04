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

type ModeledDatabaseTestSuite struct {
	suite.Suite
	database ModeledDatabase[testEntity]
}

func (i *ModeledDatabaseTestSuite) SetupSuite() {
	config := NewDefaultDatabaseConfiguration()
	config.Port = 5444
	config.DBName = "hob_test"
	i.database = NewModeledDatabase[testEntity](testEntity{}, NewDatabaseService(config))

	err := i.database.DB().AutoMigrate(testEntity{})

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create new entity")
	}
}

func (i *ModeledDatabaseTestSuite) TearDownSuite() {
	err := i.database.DB().Migrator().DropTable(testEntity{})

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot drop table")
	}
}

func (i *ModeledDatabaseTestSuite) TearDownTest() {
	err := i.database.DB().Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(testEntity{}).Error

	if err != nil {
		log.Err(err).Msg("Cannot truncate table")
	}
}

func TestModeledDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(ModeledDatabaseTestSuite))
}

func (i *ModeledDatabaseTestSuite) Test_Find() {
	entity := i.createTestEntity()

	actual, err := i.database.Find(entity.Id)

	assert.Nil(i.T(), err)

	assert.Equal(i.T(), entity, actual)
}

func (i *ModeledDatabaseTestSuite) Test_Find_WithNotExistsId() {
	actual, err := i.database.Find(uuid.New())

	assert.Nil(i.T(), err)

	assert.Equal(i.T(), testEntity{}, actual)
}

func (i *ModeledDatabaseTestSuite) Test_FindReceiver() {
	entity := i.createTestEntity()

	var receiver testEntityDto

	err := i.database.FindReceiver(&receiver, entity.Id)

	assert.Nil(i.T(), err)

	assert.Equal(i.T(), entity.toDto(), receiver)
}

func (i *ModeledDatabaseTestSuite) Test_FindReceiver_WithNotExistsId() {
	var receiver testEntityDto

	err := i.database.FindReceiver(&receiver, uuid.New())

	assert.Nil(i.T(), err)

	assert.Equal(i.T(), testEntityDto{}, receiver)
}

func (i *ModeledDatabaseTestSuite) Test_FindBy() {
	entity := i.createTestEntity()

	actual, err := i.database.FindBy("name = ?", entity.Name)

	assert.Nil(i.T(), err)

	assert.Equal(i.T(), entity, actual)
}

func (i *ModeledDatabaseTestSuite) Test_FindBy_WithNotExistsId() {
	actual, err := i.database.FindBy("name = ?", "not_exists")

	assert.Nil(i.T(), err)

	assert.Equal(i.T(), testEntity{}, actual)
}

func (i *ModeledDatabaseTestSuite) Test_FindReceiverBy() {
	entity := i.createTestEntity()

	var receiver testEntityDto

	err := i.database.FindReceiverBy(&receiver, "name = ?", entity.Name)

	assert.Nil(i.T(), err)

	assert.Equal(i.T(), entity.toDto(), receiver)
}

func (i *ModeledDatabaseTestSuite) Test_FindReceiverBy_WithNotExistsId() {
	var receiver testEntityDto

	err := i.database.FindReceiverBy(&receiver, "name = ?", "not_exists")

	assert.Nil(i.T(), err)

	assert.Equal(i.T(), testEntityDto{}, receiver)
}

func (i *ModeledDatabaseTestSuite) Test_First() {
	entity := i.createTestEntity()

	actual, err := i.database.First(entity.Id)

	assert.Nil(i.T(), err)

	assert.Equal(i.T(), entity, actual)
}

func (i *ModeledDatabaseTestSuite) Test_First_WithNotExistsId() {
	actual, err := i.database.First(uuid.New())

	assert.Equal(i.T(), gorm.ErrRecordNotFound, err)

	assert.Equal(i.T(), testEntity{}, actual)
}

func (i *ModeledDatabaseTestSuite) Test_FirstReceiver() {
	entity := i.createTestEntity()

	var receiver testEntityDto

	err := i.database.FirstReceiver(&receiver, entity.Id)

	assert.Nil(i.T(), err)

	assert.Equal(i.T(), entity.toDto(), receiver)
}

func (i *ModeledDatabaseTestSuite) Test_FirstReceiver_WithNotExistsId() {
	var receiver testEntityDto

	err := i.database.FirstReceiver(&receiver, uuid.New())

	assert.Equal(i.T(), gorm.ErrRecordNotFound, err)

	assert.Equal(i.T(), testEntityDto{}, receiver)
}

func (i *ModeledDatabaseTestSuite) Test_FirstBy() {
	entity := i.createTestEntity()

	actual, err := i.database.FirstBy("name = ?", entity.Name)

	assert.Nil(i.T(), err)

	assert.Equal(i.T(), entity, actual)
}

func (i *ModeledDatabaseTestSuite) Test_FirstBy_WithNotExistsId() {
	actual, err := i.database.FirstBy("name = ?", "not_exists")

	assert.Equal(i.T(), gorm.ErrRecordNotFound, err)

	assert.Equal(i.T(), testEntity{}, actual)
}

func (i *ModeledDatabaseTestSuite) Test_FirstReceiverBy() {
	entity := i.createTestEntity()

	var receiver testEntityDto

	err := i.database.FirstReceiverBy(&receiver, "name = ?", entity.Name)

	assert.Nil(i.T(), err)

	assert.Equal(i.T(), entity.toDto(), receiver)
}

func (i *ModeledDatabaseTestSuite) Test_FirstReceiverBy_WithNotExistsId() {
	var receiver testEntityDto

	err := i.database.FirstReceiverBy(&receiver, "name = ?", "not_exists")

	assert.Equal(i.T(), gorm.ErrRecordNotFound, err)

	assert.Equal(i.T(), testEntityDto{}, receiver)
}

func (i *ModeledDatabaseTestSuite) Test_Exists() {
	entity := i.createTestEntity()

	exists := i.database.Exists(entity.Id)

	assert.True(i.T(), exists)
}

func (i *ModeledDatabaseTestSuite) Test_Exists_WithNotExistsId() {
	exists := i.database.Exists(uuid.New())

	assert.False(i.T(), exists)
}

func (i *ModeledDatabaseTestSuite) Test_ExistsBy() {
	entity := i.createTestEntity()

	exists := i.database.ExistsBy("name = ?", entity.Name)

	assert.True(i.T(), exists)
}

func (i *ModeledDatabaseTestSuite) Test_ExistsBy_WithNotExists() {
	exists := i.database.ExistsBy("name = ?", "name match")

	assert.False(i.T(), exists)
}

func (i *ModeledDatabaseTestSuite) Test_Delete() {
	entity := i.createTestEntity()

	err := i.database.Delete(entity.Id)

	assert.Nil(i.T(), err)
	assert.False(i.T(), i.database.Exists(entity.Id))
}

func (i *ModeledDatabaseTestSuite) Test_Delete_WithNotExists() {
	err := i.database.Delete(uuid.New())

	assert.Nil(i.T(), err)
}

func (i *ModeledDatabaseTestSuite) Test_Update() {
	entity := i.createTestEntity()

	err := i.database.Update(entity.Id, testEntityUpdateDto{
		Id:          entity.Id,
		Name:        fmt.Sprintf("%s-new", entity.Name),
		Description: fmt.Sprintf("%s-new", entity.Description),
		Value:       entity.Value + 100,
	})

	assert.Nil(i.T(), err)

	newEntity, err := i.database.Find(entity.Id)

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), testEntity{
		Id:          entity.Id,
		Name:        "Name-new",
		Description: "Description-new",
		Value:       200,
	}, newEntity)
}

func (i *ModeledDatabaseTestSuite) Test_Update_WithOmitColumns() {
	entity := i.createTestEntity()

	err := i.database.Update(entity.Id, testEntityUpdateDto{
		Id:          entity.Id,
		Name:        fmt.Sprintf("%s-new", entity.Name),
		Description: fmt.Sprintf("%s-new", entity.Description),
		Value:       entity.Value + 100,
	}, "Value")

	assert.Nil(i.T(), err)

	newEntity, err := i.database.Find(entity.Id)

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), testEntity{
		Id:          entity.Id,
		Name:        "Name-new",
		Description: "Description-new",
		Value:       100,
	}, newEntity)
}

func (i *ModeledDatabaseTestSuite) Test_UpdateById_WithNotExists() {
	err := i.database.Update(uuid.New(), testEntityUpdateDto{
		Id:          uuid.New(),
		Name:        "new name",
		Description: "new description",
		Value:       100,
	})

	assert.Nil(i.T(), err)
}

func (i *ModeledDatabaseTestSuite) createTestEntity() testEntity {
	entity := generateTestEntity()

	err := i.database.Create(&entity)

	assert.Nil(i.T(), err)

	return entity
}

func (t testEntity) toDto() testEntityDto {
	return testEntityDto{
		Id:          t.Id,
		Name:        t.Name,
		Description: t.Description,
	}
}
