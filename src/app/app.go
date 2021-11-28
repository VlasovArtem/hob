package app

import (
	"common/dependency"
	"common/handler"
	"db"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"os"
	"reflect"
	"strconv"
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
	host := getEnvironmentVariable(hostEnvironmentName, "localhost")
	port := getEnvironmentIntVariable(portEnvironmentName, 5432)
	user := getEnvironmentVariable(userEnvironmentName, "postgres")
	password := getEnvironmentVariable(passwordEnvironmentName, "postgres")
	dbname := getEnvironmentVariable(dbnameEnvironmentName, "hob")

	configuration := db.DatabaseConfiguration{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbname,
	}

	a.AddDependency(configuration)
}

func getEnvironmentVariable(name string, defaultValue string) string {
	if variable := os.Getenv(name); variable == "" {
		log.Info().Msgf("Environment variable with name '%s' not found, default used '%s'", name, defaultValue)
		return defaultValue
	} else {
		return variable
	}
}

func getEnvironmentIntVariable(name string, defaultValue int) int {
	if variable := os.Getenv(name); variable == "" {
		log.Info().Msgf("Environment variable with name '%s' not found, default used '%d'", name, defaultValue)
		return defaultValue
	} else {
		if intVariable, err := strconv.Atoi(variable); err != nil {
			log.Fatal().Err(err)
			return 0
		} else {
			return intVariable
		}
	}
}

func (a *ApplicationObject) migrate(object dependency.ObjectDatabaseMigrator) {
	if a.databaseService == nil {
		log.Fatal().Msg("DatabaseService is not initialized")
	}
	if err := a.databaseService.D().AutoMigrate(object.GetEntity()); err != nil {
		log.Fatal().Err(err)
	}
}
