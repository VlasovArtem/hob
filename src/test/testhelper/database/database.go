package database

import (
	"fmt"
	database "github.com/VlasovArtem/hob/src/db"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"reflect"
)

type DBTestSuite struct {
	suite.Suite
	Database   database.DatabaseService
	beforeTest []func(service database.DatabaseService)
	afterTest  []func(service database.DatabaseService)
	afterSuite []func(service database.DatabaseService)
}

func (db *DBTestSuite) InitDBTestSuite() {
	config := database.NewDefaultDatabaseConfiguration()
	config.DBName = "hob_test"
	db.Database = database.NewDatabaseService(config)
}

func (db *DBTestSuite) CreateRepository(provider func(service database.DatabaseService)) *DBTestSuite {
	provider(db.Database)

	return db
}

func (db *DBTestSuite) ExecuteMigration(migrators ...any) {
	for _, migrator := range migrators {
		if err := db.Database.D().AutoMigrate(migrator); err != nil {
			log.Fatal().Err(err).Msg("Cannot create table")
		}
	}
}

func (db *DBTestSuite) CreateEntity(entity any) {
	if err := db.Database.Create(entity); err != nil {
		log.Fatal().Err(err).Msg("Cannot create entity")
	}
}

func (db *DBTestSuite) AddBeforeTest(beforeTest func(service database.DatabaseService)) *DBTestSuite {
	db.beforeTest = append(db.beforeTest, beforeTest)

	return db
}

func (db *DBTestSuite) AddAfterTest(afterTest func(service database.DatabaseService)) *DBTestSuite {
	db.afterTest = append(db.afterTest, afterTest)

	return db
}

func (db *DBTestSuite) AddAfterSuite(afterSuite func(service database.DatabaseService)) *DBTestSuite {
	db.afterSuite = append(db.afterSuite, afterSuite)

	return db
}

func (db *DBTestSuite) BeforeTest(suiteName, testName string) {
	for _, function := range db.beforeTest {
		function(db.Database)
	}
}

func (db *DBTestSuite) TearDownSuite() {
	for _, afterSuiteFunc := range db.afterSuite {
		afterSuiteFunc(db.Database)
	}
}

func (db *DBTestSuite) TearDownTest() {
	for _, afterTestFunc := range db.afterTest {
		afterTestFunc(db.Database)
	}
}

func (db *DBTestSuite) Delete(entity any) {
	entityValue := reflect.Indirect(reflect.ValueOf(entity))
	idValue := entityValue.FieldByName("Id")
	valueId := fmt.Sprintf("%v", idValue)
	if parse, err := uuid.Parse(valueId); err != nil {
		log.Fatal().Err(err)
	} else {
		if err = db.Database.DeleteById(entity, parse); err != nil {
			log.Fatal().Err(err)
		}
	}
}

func TruncateTable(service database.DatabaseService, entity any) {
	service.D().Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(entity)
}
