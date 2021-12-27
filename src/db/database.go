package db

import (
	"fmt"
	"github.com/VlasovArtem/hob/src/common/dependency"
	helper "github.com/VlasovArtem/hob/src/common/service"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

type DatabaseConfiguration struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func NewDefaultDatabaseConfiguration() DatabaseConfiguration {
	return DatabaseConfiguration{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		DBName:   "hob",
	}
}

type DatabaseObject struct {
	db *gorm.DB
}

func (d *DatabaseObject) Initialize(factory dependency.DependenciesProvider) interface{} {
	return NewDatabaseService(factory.FindRequiredByObject(DatabaseConfiguration{}).(DatabaseConfiguration))
}

func NewDatabaseService(config DatabaseConfiguration) DatabaseService {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.User, config.Password, config.DBName)

	if db, err := gorm.Open(postgres.Open(psqlconn)); err != nil {
		helper.LogError(err, "")
		os.Exit(1)
		return nil
	} else {
		return &DatabaseObject{db}
	}
}

type DatabaseService interface {
	Create(value interface{}) error
	FindById(receiver interface{}, id uuid.UUID) error
	FindByIdModeled(model interface{}, receiver interface{}, id uuid.UUID) error
	FindByQuery(receiver interface{}, query interface{}, conditions ...interface{}) error
	FindByModeled(model interface{}, receiver interface{}, query interface{}, conditions ...interface{}) error
	FirstByModeled(model interface{}, receiver interface{}, query interface{}, conditions ...interface{}) error
	ExistsById(model interface{}, id uuid.UUID) (exists bool)
	ExistsByQuery(model interface{}, query interface{}, args ...interface{}) (exists bool)
	DeleteById(model interface{}, id uuid.UUID) error
	UpdateById(model interface{}, id uuid.UUID, entity interface{}, omit ...string) error
	D() *gorm.DB
	DM(model interface{}) *gorm.DB
}

func (d *DatabaseObject) Create(value interface{}) error {
	tx := d.db.Create(value)

	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (d *DatabaseObject) FindById(receiver interface{}, id uuid.UUID) error {
	return d.FindByIdModeled(receiver, receiver, id)
}

func (d *DatabaseObject) FindByIdModeled(model interface{}, receiver interface{}, id uuid.UUID) error {
	tx := d.db.Model(model).First(receiver, id)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (d *DatabaseObject) FindByQuery(receiver interface{}, query interface{}, conditions ...interface{}) error {
	return d.FindByModeled(receiver, receiver, query, conditions)
}

func (d *DatabaseObject) FindByModeled(model interface{}, receiver interface{}, query interface{}, conditions ...interface{}) error {
	tx := d.db.Model(model).Where(query, conditions...).Find(receiver)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (d *DatabaseObject) FirstByModeled(model interface{}, receiver interface{}, query interface{}, conditions ...interface{}) error {
	tx := d.db.Model(model).Where(query, conditions...).First(receiver)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (d *DatabaseObject) ExistsById(model interface{}, id uuid.UUID) (exists bool) {
	tx := d.db.Model(model).Select("count(*) > 0").Where("id = ?", id)
	if err := tx.Find(&exists).Error; err != nil {
		log.Error().Err(err).Msg("")
	}
	return exists
}

func (d *DatabaseObject) ExistsByQuery(model interface{}, query interface{}, args ...interface{}) (exists bool) {
	tx := d.db.Model(model).Select("count(*) > 0").Where(query, args...)
	if err := tx.Find(&exists).Error; err != nil {
		log.Error().Err(err).Msg("")
	}
	return exists
}

func (d *DatabaseObject) DeleteById(model interface{}, id uuid.UUID) error {
	return d.db.Delete(model, id).Error
}

func (d *DatabaseObject) UpdateById(model interface{}, id uuid.UUID, entity interface{}, omit ...string) error {
	omitColumns := []string{"Id"}
	omitColumns = append(omitColumns, omit...)

	return d.db.Model(model).Where("id = ?", id).Omit(omitColumns...).Updates(entity).Error
}

func (d *DatabaseObject) D() *gorm.DB {
	return d.db
}

func (d *DatabaseObject) DM(model interface{}) *gorm.DB {
	return d.db.Model(model)
}
