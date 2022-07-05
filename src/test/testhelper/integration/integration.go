package integration

import (
	"fmt"
	database "github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/test/testhelper"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"reflect"
)

type Suite[T any] struct {
	suite.Suite
	O               T
	DatabaseService database.DatabaseService
	beforeTest      []func(object T)
	afterTest       []func(object T)
	afterSuite      []func(object T)
	migrators       []any
}

func (db *Suite[T]) InitSuite(provider func(service database.DatabaseService) T) {
	config := database.NewDefaultDatabaseConfiguration()
	config.Port = 5444
	config.DBName = "hob_test"
	config.GormConfig = &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}
	db.DatabaseService = database.NewDatabaseService(config)
	db.O = provider(db.DatabaseService)
}

func (db *Suite[T]) ExecuteMigration(migrators ...any) {
	for _, migrator := range migrators {
		if err := db.DatabaseService.DB().AutoMigrate(migrator); err != nil {
			log.Fatal().Err(err).Msg("Cannot create table")
		}
	}
	db.migrators = migrators
}

func (db *Suite[T]) CreateEntity(entity any, omit ...string) {
	if err := db.DatabaseService.Create(entity, omit...); err != nil {
		log.Fatal().Err(err).Msg("Cannot create entity")
	}
}

func (db *Suite[T]) AddBeforeTest(beforeTest func(object T)) *Suite[T] {
	db.beforeTest = append(db.beforeTest, beforeTest)

	return db
}

func (db *Suite[T]) AddAfterTest(afterTest func(object T)) *Suite[T] {
	db.afterTest = append(db.afterTest, afterTest)

	return db
}

func (db *Suite[T]) AddAfterSuite(afterSuite func(object T)) *Suite[T] {
	db.afterSuite = append(db.afterSuite, afterSuite)

	return db
}

func (db *Suite[T]) BeforeTest(suiteName, testName string) {
	for _, function := range db.beforeTest {
		function(db.O)
	}
}

func (db *Suite[T]) TearDownSuite() {
	for _, afterSuiteFunc := range db.afterSuite {
		afterSuiteFunc(db.O)
	}
	for i := len(db.migrators) - 1; i >= 0; i-- {
		testhelper.TruncateTable(db.DatabaseService, db.migrators[i])
	}
}

func (db *Suite[T]) TearDownTest() {
	for _, afterTestFunc := range db.afterTest {
		afterTestFunc(db.O)
	}
}

func (db *Suite[T]) Delete(entity any) {
	entityValue := reflect.Indirect(reflect.ValueOf(entity))
	idValue := entityValue.FieldByName("Id")
	valueId := fmt.Sprintf("%v", idValue)
	if parse, err := uuid.Parse(valueId); err != nil {
		log.Fatal().Err(err)
	} else {
		if err = db.DatabaseService.DB().Delete(entity, parse).Error; err != nil {
			log.Fatal().Err(err)
		}
	}
}
