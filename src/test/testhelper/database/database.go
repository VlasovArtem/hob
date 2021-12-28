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

func DropTable(db *gorm.DB, model any) {
	err := db.Migrator().DropTable(model)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot drop table")
	}
}

func CreateTable(db *gorm.DB, model any) {
	err := db.AutoMigrate(model)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create table")
	}
}

func RecreateTable(db *gorm.DB, model any) {
	DropTable(db, model)
	CreateTable(db, model)
}

type DBTestSuite struct {
	suite.Suite
	Database         database.DatabaseService
	migrators        []any
	beforeTest       []func(service database.DatabaseService)
	createdEntities  []any
	constantEntities []any
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

func (db *DBTestSuite) AddMigrators(migrators ...any) {
	for _, migrator := range migrators {
		if err := db.Database.D().AutoMigrate(migrator); err != nil {
			log.Fatal().Err(err).Msg("Cannot create table")
		}

		db.migrators = append(db.migrators, migrator)
	}
}

func (db *DBTestSuite) CreateEntity(entity any) {
	if err := db.Database.Create(entity); err != nil {
		log.Fatal().Err(err).Msg("Cannot create entity")
	}

	db.createdEntities = append(db.createdEntities, reflect.Indirect(reflect.ValueOf(entity)).Interface())
}

func (db *DBTestSuite) CreateConstantEntity(entity any) {
	if err := db.Database.Create(entity); err != nil {
		log.Fatal().Err(err).Msg("Cannot create entity")
	}

	db.constantEntities = append(db.constantEntities, reflect.Indirect(reflect.ValueOf(entity)).Interface())
}

func (db *DBTestSuite) AddBeforeTest(beforeTest func(service database.DatabaseService)) *DBTestSuite {
	db.beforeTest = append(db.beforeTest, beforeTest)

	return db
}

func (db *DBTestSuite) BeforeTest(suiteName, testName string) {
	for _, function := range db.beforeTest {
		function(db.Database)
	}
}

func (db *DBTestSuite) TearDownSuite() {
	db.deleteConstant()
}

func (db *DBTestSuite) TearDownTest() {
	db.deleteCreated()
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

func (db *DBTestSuite) deleteCreated() {
	for i := len(db.createdEntities) - 1; i >= 0; i-- {
		db.Delete(db.createdEntities[i])
	}
}

func (db *DBTestSuite) deleteConstant() {
	for i := len(db.constantEntities) - 1; i >= 0; i-- {
		db.Delete(db.constantEntities[i])
	}
}
