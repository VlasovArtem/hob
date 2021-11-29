package database

import (
	database "github.com/VlasovArtem/hob/src/db"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

func DropTable(db *gorm.DB, model interface{}) {
	err := db.Migrator().DropTable(model)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot drop table")
	}
}

func CreateTable(db *gorm.DB, model interface{}) {
	err := db.AutoMigrate(model)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create table")
	}
}

func RecreateTable(db *gorm.DB, model interface{}) {
	DropTable(db, model)
	CreateTable(db, model)
}

type DBTestSuite struct {
	suite.Suite
	Database        database.DatabaseService
	migrators       []interface{}
	beforeTest      []func(service database.DatabaseService)
	createdEntities []interface{}
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

func (db *DBTestSuite) AddMigrators(migrators ...interface{}) {
	for _, migrator := range migrators {
		if err := db.Database.D().AutoMigrate(migrator); err != nil {
			log.Fatal().Err(err).Msg("Cannot create table")
		}

		db.migrators = append(db.migrators, migrator)
	}
}

func (db *DBTestSuite) CreateEntity(entity interface{}) {
	if err := db.Database.Create(entity); err != nil {
		log.Fatal().Err(err).Msg("Cannot create entity")
	}

	db.createdEntities = append(db.createdEntities, entity)
}

func (db *DBTestSuite) TearDown() {
	for _, entity := range db.createdEntities {
		db.Database.D().Delete(entity)
	}
}

func (db *DBTestSuite) AddBeforeTest(beforeTest func(service database.DatabaseService)) *DBTestSuite {
	db.beforeTest = append(db.beforeTest, beforeTest)

	return db
}

func (db *DBTestSuite) BeforeTest(suiteName, testName string) {
	for _, beforeTestFunction := range db.beforeTest {
		beforeTestFunction(db.Database)
	}
}
