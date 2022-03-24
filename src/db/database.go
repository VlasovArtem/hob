package db

import (
	"fmt"
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseConfiguration struct {
	Host       string
	Port       int
	User       string
	Password   string
	DBName     string
	GormConfig *gorm.Config
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

func (d *DatabaseObject) Initialize(factory dependency.DependenciesProvider) any {
	return NewDatabaseService(factory.FindRequiredByObject(DatabaseConfiguration{}).(DatabaseConfiguration))
}

func NewDatabaseService(config DatabaseConfiguration) DatabaseService {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.User, config.Password, config.DBName)

	var options []gorm.Option

	if config.GormConfig != nil {
		options = append(options, config.GormConfig)
	}

	if db, err := gorm.Open(postgres.Open(psqlconn), options...); err != nil {
		log.Fatal().Err(err)
		return nil
	} else {
		return &DatabaseObject{db}
	}
}

type DatabaseService interface {
	Create(value any, omit ...string) error
	FindById(receiver any, id uuid.UUID) error
	FindByIdModeled(model any, receiver any, id uuid.UUID) error
	FindByQuery(receiver any, query any, conditions ...any) error
	FindByModeled(model any, receiver any, query any, conditions ...any) error
	FirstByModeled(model any, receiver any, query any, conditions ...any) error
	ExistsById(model any, id uuid.UUID) (exists bool)
	ExistsByQuery(model any, query any, args ...any) (exists bool)
	DeleteById(model any, id uuid.UUID) error
	UpdateById(model any, id uuid.UUID, entity any, omit ...string) error
	D() *gorm.DB
	DM(model any) *gorm.DB
}

func (d *DatabaseObject) Create(value any, omit ...string) error {
	tx := d.db.Omit(omit...).Create(value)

	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (d *DatabaseObject) FindById(receiver any, id uuid.UUID) error {
	return d.FindByIdModeled(receiver, receiver, id)
}

func (d *DatabaseObject) FindByIdModeled(model any, receiver any, id uuid.UUID) error {
	tx := d.db.Model(model).First(receiver, id)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (d *DatabaseObject) FindByQuery(receiver any, query any, conditions ...any) error {
	return d.FindByModeled(receiver, receiver, query, conditions)
}

func (d *DatabaseObject) FindByModeled(model any, receiver any, query any, conditions ...any) error {
	tx := d.db.Model(model).Where(query, conditions...).Find(receiver)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (d *DatabaseObject) FirstByModeled(model any, receiver any, query any, conditions ...any) error {
	tx := d.db.Model(model).Where(query, conditions...).First(receiver)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (d *DatabaseObject) ExistsById(model any, id uuid.UUID) (exists bool) {
	tx := d.db.Model(model).Select("count(*) > 0").Where("id = ?", id)
	if err := tx.Find(&exists).Error; err != nil {
		log.Error().Err(err).Msg("")
	}
	return exists
}

func (d *DatabaseObject) ExistsByQuery(model any, query any, args ...any) (exists bool) {
	tx := d.db.Model(model).Select("count(*) > 0").Where(query, args...)
	if err := tx.Find(&exists).Error; err != nil {
		log.Error().Err(err).Msg("")
	}
	return exists
}

func (d *DatabaseObject) DeleteById(model any, id uuid.UUID) error {
	return d.db.Delete(model, id).Error
}

func (d *DatabaseObject) UpdateById(model any, id uuid.UUID, entity any, omit ...string) error {
	omitColumns := []string{"Id"}
	omitColumns = append(omitColumns, omit...)

	return d.db.Model(model).Where("id = ?", id).Omit(omitColumns...).Updates(entity).Error
}

func (d *DatabaseObject) D() *gorm.DB {
	return d.db
}

func (d *DatabaseObject) DM(model any) *gorm.DB {
	return d.db.Model(model)
}
