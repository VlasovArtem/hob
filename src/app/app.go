package app

import (
	"encoding/json"
	"fmt"
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/environment"
	helper "github.com/VlasovArtem/hob/src/common/service"
	"github.com/VlasovArtem/hob/src/config"
	"github.com/VlasovArtem/hob/src/country/model"
	countries "github.com/VlasovArtem/hob/src/country/service"
	"github.com/VlasovArtem/hob/src/db"
	houseRepository "github.com/VlasovArtem/hob/src/house/respository"
	houseService "github.com/VlasovArtem/hob/src/house/service"
	incomeRepository "github.com/VlasovArtem/hob/src/income/repository"
	incomeSchedulerRepository "github.com/VlasovArtem/hob/src/income/scheduler/repository"
	incomeSchedulerService "github.com/VlasovArtem/hob/src/income/scheduler/service"
	incomeService "github.com/VlasovArtem/hob/src/income/service"
	meterRepository "github.com/VlasovArtem/hob/src/meter/repository"
	meterService "github.com/VlasovArtem/hob/src/meter/service"
	paymentRepository "github.com/VlasovArtem/hob/src/payment/repository"
	paymentSchedulerRepository "github.com/VlasovArtem/hob/src/payment/scheduler/repository"
	paymentSchedulerService "github.com/VlasovArtem/hob/src/payment/scheduler/service"
	paymentService "github.com/VlasovArtem/hob/src/payment/service"
	customProviderRepository "github.com/VlasovArtem/hob/src/provider/custom/repository"
	customProviderService "github.com/VlasovArtem/hob/src/provider/custom/service"
	providerRepository "github.com/VlasovArtem/hob/src/provider/repository"
	providerService "github.com/VlasovArtem/hob/src/provider/service"
	"github.com/VlasovArtem/hob/src/scheduler"
	userRepository "github.com/VlasovArtem/hob/src/user/repository"
	userService "github.com/VlasovArtem/hob/src/user/service"
	userRequestValidator "github.com/VlasovArtem/hob/src/user/validator"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
	"reflect"
)

const (
	hostEnvironmentName     = "DB_HOST"
	portEnvironmentName     = "DB_PORT"
	userEnvironmentName     = "DB_USER"
	passwordEnvironmentName = "DB_PASSWORD"
	dbnameEnvironmentName   = "DB_NAME"
	countriesDirVariable    = "COUNTRIES_DIR"
)

var migratorType = reflect.TypeOf((*dependency.ObjectDatabaseMigrator)(nil)).Elem()

type RootApplication struct {
	DependenciesFactory dependency.DependenciesFactory
	databaseService     db.DatabaseService
	Config              *config.Config
}

// NewRootApplication initializers is a list of auto initializing objects. Order is required.
// The application creates db.DatabaseConfiguration automatically
func NewRootApplication(config *config.Config) *RootApplication {
	applicationService := &RootApplication{
		DependenciesFactory: dependency.NewDependenciesFactory(),
		Config:              config,
	}

	applicationService.createDatabaseConfiguration()

	applicationService.DependenciesFactory.AddAutoDependency(new(db.DatabaseObject))

	applicationService.databaseService = applicationService.DependenciesFactory.FindRequiredByObject(db.DatabaseObject{}).(db.DatabaseService)

	applicationService.createCountriesService()

	applicationService.addAutoInitializingDependencies()

	return applicationService
}

func (a *RootApplication) createDatabaseConfiguration() {
	configuration := db.DatabaseConfiguration{
		Host:     environment.GetEnvironmentVariable(hostEnvironmentName, "localhost"),
		Port:     environment.GetEnvironmentIntVariable(portEnvironmentName, 5432),
		User:     environment.GetEnvironmentVariable(userEnvironmentName, "postgres"),
		Password: environment.GetEnvironmentVariable(passwordEnvironmentName, "postgres"),
		DBName:   environment.GetEnvironmentVariable(dbnameEnvironmentName, "hob"),
	}

	a.DependenciesFactory.Add(configuration)
}

func (a *RootApplication) addAutoInitializingDependencies() {
	initializers := []dependency.ObjectDependencyInitializer{
		new(userRequestValidator.UserRequestValidatorObject),
		new(userRepository.UserRepositoryObject),
		new(userService.UserServiceObject),
		new(houseRepository.HouseRepositoryObject),
		new(houseService.HouseServiceObject),
		new(scheduler.SchedulerServiceObject),
		new(providerRepository.ProviderRepositoryObject),
		new(providerService.ProviderServiceObject),
		new(customProviderRepository.CustomProviderRepositoryObject),
		new(customProviderService.CustomProviderServiceObject),
		new(paymentRepository.PaymentRepositoryObject),
		new(paymentService.PaymentServiceObject),
		new(paymentSchedulerRepository.PaymentSchedulerRepositoryObject),
		new(paymentSchedulerService.PaymentSchedulerServiceObject),
		new(meterRepository.MeterRepositoryObject),
		new(meterService.MeterServiceObject),
		new(incomeRepository.IncomeRepositoryObject),
		new(incomeService.IncomeServiceObject),
		new(incomeSchedulerRepository.IncomeSchedulerRepositoryObject),
		new(incomeSchedulerService.IncomeSchedulerServiceObject),
	}

	for _, initializer := range initializers {
		autoDependency := a.DependenciesFactory.AddAutoDependency(initializer)

		if reflect.TypeOf(autoDependency).Implements(migratorType) {
			a.migrate(autoDependency.(dependency.ObjectDatabaseMigrator))
		}
	}
}

func (a *RootApplication) migrate(object dependency.ObjectDatabaseMigrator) {
	if a.databaseService == nil {
		log.Fatal().Msg("DatabaseService is not initialized")
	}
	if err := a.databaseService.D().AutoMigrate(object.GetEntity()); err != nil {
		log.Fatal().Err(err)
	}
}

func (a *RootApplication) createCountriesService() {
	file, err := ioutil.ReadFile(fmt.Sprintf("%scontent/countries.json", environment.GetEnvironmentVariable(countriesDirVariable, "./")))

	if helper.LogError(err, "Countries is not found") {
		os.Exit(1)
	}

	var countriesContent []model.Country

	if err = json.Unmarshal(file, &countriesContent); err != nil {
		log.Fatal().Err(err)
	}

	if len(countriesContent) == 0 {
		log.Fatal().Msg("countries content is empty")
	}

	a.DependenciesFactory.Add(countries.NewCountryService(countriesContent))
}
