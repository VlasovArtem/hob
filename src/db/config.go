package db

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type DatabaseConfiguration struct {
	Host       string
	Port       int
	User       string
	Password   string
	DBName     string
	GormConfig *gorm.Config
}

type myGormLogger zerolog.Logger

func (a *myGormLogger) LogMode(logLevel logger.LogLevel) logger.Interface {
	return a
}

func (a *myGormLogger) Info(context context.Context, format string, v ...interface{}) {
	log.Info().Msgf(format, v...)
}

func (a *myGormLogger) Warn(context context.Context, format string, v ...interface{}) {
	log.Warn().Msgf(format, v...)
}
func (a *myGormLogger) Error(context context.Context, format string, v ...interface{}) {
	log.Error().Msgf(format, v...)
}
func (a *myGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, affected := fc()
	log.Trace().Err(err).Msgf("Begin - %s. SQL - %s. Affected rows - %d", begin.Format(time.RFC3339), sql, affected)
}

func NewDefaultDatabaseConfiguration() DatabaseConfiguration {
	return DatabaseConfiguration{
		Host:     "localhost",
		Port:     5432,
		User:     "hob",
		Password: "magical_password",
		DBName:   "hob",
		GormConfig: &gorm.Config{
			Logger: &myGormLogger{},
		},
	}
}
