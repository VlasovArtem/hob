package app

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/environment"
	"github.com/VlasovArtem/hob/src/common/handler"
	"github.com/VlasovArtem/hob/src/db"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"reflect"
)

const (
	hostEnvironmentName     = "DB_HOST"
	portEnvironmentName     = "DB_PORT"
	userEnvironmentName     = "DB_USER"
	passwordEnvironmentName = "DB_PASSWORD"
	dbnameEnvironmentName   = "DB_NAME"
)

var (
	migratorType    = reflect.TypeOf((*dependency.ObjectDatabaseMigrator)(nil)).Elem()
	initializerType = reflect.TypeOf((*dependency.ObjectDependencyInitializer)(nil)).Elem()
	handlerType     = reflect.TypeOf((*handler.ApplicationHandler)(nil)).Elem()
)

type ApplicationObject struct {
	DependenciesFactory dependency.DependenciesFactory
	databaseService     db.DatabaseService
	router              *mux.Router
}

type ApplicationService interface {
	AddDependency(interface{}) ApplicationService
	AddAutoDependency(dependency.ObjectDependencyInitializer) ApplicationService
	AddHandler(handler.ApplicationHandler) ApplicationService
}

// NewApplicationService initializers is a list of auto initializing objects. Order is required.
// The application creates db.DatabaseConfiguration automatically
func NewApplicationService(router *mux.Router) ApplicationService {
	applicationService := &ApplicationObject{
		DependenciesFactory: dependency.NewDependenciesFactory(),
		router:              router,
	}

	applicationService.createDatabaseConfiguration()

	applicationService.AddAutoDependency(new(db.DatabaseObject))

	applicationService.databaseService = applicationService.DependenciesFactory.FindRequiredByObject(db.DatabaseObject{}).(db.DatabaseService)

	return applicationService
}

func (a *ApplicationObject) AddDependency(object interface{}) ApplicationService {
	if reflect.TypeOf(object).Implements(handlerType) {
		return a.AddHandler(object.(handler.ApplicationHandler))
	}
	if reflect.TypeOf(object).Implements(initializerType) {
		return a.AddAutoDependency(object.(dependency.ObjectDependencyInitializer))
	}

	a.DependenciesFactory.Add(object)

	if reflect.TypeOf(object).Implements(migratorType) {
		a.migrate(object.(dependency.ObjectDatabaseMigrator))
	}

	return a
}

func (a *ApplicationObject) AddAutoDependency(initializer dependency.ObjectDependencyInitializer) ApplicationService {
	initializer.Initialize(a.DependenciesFactory)

	if reflect.TypeOf(initializer).Implements(migratorType) {
		a.migrate(initializer.(dependency.ObjectDatabaseMigrator))
	}

	return a
}

func (a *ApplicationObject) AddHandler(applicationHandler handler.ApplicationHandler) ApplicationService {
	a.AddAutoDependency(applicationHandler)

	a.DependenciesFactory.FindRequiredByObject(applicationHandler).(handler.ApplicationHandler).Init(a.router)

	return a
}

func (a *ApplicationObject) createDatabaseConfiguration() {
	configuration := db.DatabaseConfiguration{
		Host:     environment.GetEnvironmentVariable(hostEnvironmentName, "localhost"),
		Port:     environment.GetEnvironmentIntVariable(portEnvironmentName, 5432),
		User:     environment.GetEnvironmentVariable(userEnvironmentName, "postgres"),
		Password: environment.GetEnvironmentVariable(passwordEnvironmentName, "postgres"),
		DBName:   environment.GetEnvironmentVariable(dbnameEnvironmentName, "hob"),
	}

	a.AddDependency(configuration)
}

func (a *ApplicationObject) migrate(object dependency.ObjectDatabaseMigrator) {
	if a.databaseService == nil {
		log.Fatal().Msg("DatabaseService is not initialized")
	}
	if err := a.databaseService.D().AutoMigrate(object.GetEntity()); err != nil {
		log.Fatal().Err(err)
	}
}
