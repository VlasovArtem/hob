package db

import (
	"fmt"
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Provider struct {
	db *gorm.DB
}

func (p *Provider) DB() *gorm.DB {
	return p.db
}

func (p *Provider) DBModeled(model any) *gorm.DB {
	return p.db.Model(model)
}

type ProviderInterface interface {
	DB() *gorm.DB
	DBModeled(model any) *gorm.DB
}

type Database struct {
	ProviderInterface
}

func (d *Database) GetRequiredDependencies() []dependency.Requirements {
	return []dependency.Requirements{
		dependency.FindNameAndType(DatabaseConfiguration{}),
	}
}

func (d *Database) Initialize(factory dependency.DependenciesProvider) any {
	return NewDatabaseService(factory.FindRequiredByObject(DatabaseConfiguration{}).(DatabaseConfiguration))
}

func NewDatabaseService(config DatabaseConfiguration) DatabaseService {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.User, config.Password, config.DBName)

	var options []gorm.Option

	if config.GormConfig != nil {
		options = append(options, config.GormConfig)
	}

	if db, err := gorm.Open(postgres.Open(psqlconn), options...); err != nil {
		log.Fatal().Err(err).Msg("failed to connect database")
		return nil
	} else {
		return &Database{&Provider{db: db}}
	}
}

type DatabaseService interface {
	ProviderInterface
	Create(value any, omit ...string) error
	GetProvider() ProviderInterface
}

func (d *Database) Create(value any, omit ...string) error {
	return d.DB().Omit(omit...).Create(value).Error
}

func (d *Database) GetProvider() ProviderInterface {
	return d.ProviderInterface
}
