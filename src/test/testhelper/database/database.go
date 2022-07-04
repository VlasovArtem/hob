package database

import (
	"fmt"
	database "github.com/VlasovArtem/hob/src/db"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"reflect"
)

type DBTestSuite[T any] struct {
	suite.Suite
	Database   database.ModeledDatabase[T]
	beforeTest []func(service database.ModeledDatabase[T])
	afterTest  []func(service database.ModeledDatabase[T])
	afterSuite []func(service database.ModeledDatabase[T])
}

func (db *DBTestSuite[T]) InitDBTestSuite() {
	config := database.NewDefaultDatabaseConfiguration()
	config.Port = 5444
	config.DBName = "hob_test"
	config.GormConfig = &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}
	var model T
	db.Database = database.NewModeledDatabase(model, database.NewDatabaseService(config))
}

func (db *DBTestSuite[T]) CreateRepository(provider func(service database.ModeledDatabase[T])) *DBTestSuite[T] {
	provider(db.Database)

	return db
}

func (db *DBTestSuite[T]) ExecuteMigration(migrators ...any) {
	for _, migrator := range migrators {
		if err := db.Database.DB().AutoMigrate(migrator); err != nil {
			log.Fatal().Err(err).Msg("Cannot create table")
		}
	}
}

func (db *DBTestSuite[T]) CreateEntity(entity any) {
	if err := db.Database.Create(entity); err != nil {
		log.Fatal().Err(err).Msg("Cannot create entity")
	}
}

func (db *DBTestSuite[T]) AddBeforeTest(beforeTest func(service database.ModeledDatabase[T])) *DBTestSuite[T] {
	db.beforeTest = append(db.beforeTest, beforeTest)

	return db
}

func (db *DBTestSuite[T]) AddAfterTest(afterTest func(service database.ModeledDatabase[T])) *DBTestSuite[T] {
	db.afterTest = append(db.afterTest, afterTest)

	return db
}

func (db *DBTestSuite[T]) AddAfterSuite(afterSuite func(service database.ModeledDatabase[T])) *DBTestSuite[T] {
	db.afterSuite = append(db.afterSuite, afterSuite)

	return db
}

func (db *DBTestSuite[T]) BeforeTest(suiteName, testName string) {
	for _, function := range db.beforeTest {
		function(db.Database)
	}
}

func (db *DBTestSuite[T]) TearDownSuite() {
	for _, afterSuiteFunc := range db.afterSuite {
		afterSuiteFunc(db.Database)
	}
}

func (db *DBTestSuite[T]) TearDownTest() {
	for _, afterTestFunc := range db.afterTest {
		afterTestFunc(db.Database)
	}
}

func (db *DBTestSuite[T]) Delete(entity any) {
	entityValue := reflect.Indirect(reflect.ValueOf(entity))
	idValue := entityValue.FieldByName("Id")
	valueId := fmt.Sprintf("%v", idValue)
	if parse, err := uuid.Parse(valueId); err != nil {
		log.Fatal().Err(err)
	} else {
		if err = db.Database.Delete(parse); err != nil {
			log.Fatal().Err(err)
		}
	}
}

func TruncateTable[T any](service database.ModeledDatabase[T], model any) {
	err := service.DB().Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(model).Error

	if err != nil {
		log.Err(err).Msg("Cannot truncate table")
	}
}

func TruncateTableCascade[T any](service database.ModeledDatabase[T], tableName string) {
	err := service.DB().Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", tableName)).Error

	if err != nil {
		log.Err(err).Msg("Cannot truncate table")
	}
}
