package app

import (
	"encoding/json"
	"fmt"
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/environment"
	"github.com/VlasovArtem/hob/src/config"
	"github.com/VlasovArtem/hob/src/country/model"
	countries "github.com/VlasovArtem/hob/src/country/service"
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/group/repository"
	groupService "github.com/VlasovArtem/hob/src/group/service"
	houseRepository "github.com/VlasovArtem/hob/src/house/repository"
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
	"github.com/VlasovArtem/hob/src/pivotal/cache"
	pivotalModel "github.com/VlasovArtem/hob/src/pivotal/model"
	pivotalRepository "github.com/VlasovArtem/hob/src/pivotal/repository"
	pivotalService "github.com/VlasovArtem/hob/src/pivotal/service"
	providerRepository "github.com/VlasovArtem/hob/src/provider/repository"
	providerService "github.com/VlasovArtem/hob/src/provider/service"
	"github.com/VlasovArtem/hob/src/scheduler"
	userRepository "github.com/VlasovArtem/hob/src/user/repository"
	userService "github.com/VlasovArtem/hob/src/user/service"
	userRequestValidator "github.com/VlasovArtem/hob/src/user/validator"
	"github.com/rs/zerolog/log"
	"io/ioutil"
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
	DependenciesFactory dependency.DependenciesProvider
	databaseService     db.DatabaseService
	Config              *config.Config
}

// NewRootApplication initializers is a list of auto initializing objects. Order is required.
// The application creates db.DatabaseConfiguration automatically
func NewRootApplication(config *config.Config) *RootApplication {
	applicationService := &RootApplication{
		DependenciesFactory: dependency.NewDependenciesProvider(),
		Config:              config,
	}

	applicationService.createDatabaseConfiguration()

	applicationService.DependenciesFactory.AddAutoDependency(new(db.Database))

	applicationService.databaseService = applicationService.DependenciesFactory.FindRequiredByObject(db.Database{}).(db.DatabaseService)

	applicationService.createCountriesService()

	applicationService.addAutoInitializingDependencies()

	return applicationService
}

func (a *RootApplication) createDatabaseConfiguration() {
	configuration := db.NewDefaultDatabaseConfiguration()
	configuration.Host = environment.GetEnvironmentVariable(hostEnvironmentName, "localhost")
	configuration.Port = environment.GetEnvironmentIntVariable(portEnvironmentName, 5432)
	configuration.User = environment.GetEnvironmentVariable(userEnvironmentName, "hob")
	configuration.Password = environment.GetEnvironmentVariable(passwordEnvironmentName, "magical_password")
	configuration.DBName = environment.GetEnvironmentVariable(dbnameEnvironmentName, "hob")

	a.DependenciesFactory.Add(configuration)
}

func (a *RootApplication) addAutoInitializingDependencies() {
	initializers := []dependency.ObjectDependencyInitializer{
		new(cache.PivotalCacheObject),
		new(userRequestValidator.UserRequestValidatorObject),
		new(userRepository.UserRepositoryObject),
		new(userService.UserServiceObject),
		new(repository.GroupRepositoryObject),
		new(groupService.GroupServiceStr),
		new(houseRepository.houseRepositoryStruct),
		new(houseService.HouseServiceStr),
		new(scheduler.SchedulerServiceObject),
		new(providerRepository.ProviderRepositoryStr),
		new(providerService.ProviderServiceStr),
		new(paymentRepository.PaymentRepositoryStr),
		new(paymentService.PaymentServiceStr),
		new(paymentSchedulerRepository.PaymentSchedulerRepositoryStr),
		new(paymentSchedulerService.PaymentSchedulerServiceStr),
		new(meterRepository.MeterRepositoryStr),
		new(meterService.MeterServiceObject),
		new(incomeRepository.IncomeRepositoryStr),
		new(incomeService.IncomeServiceStr),
		new(incomeSchedulerRepository.IncomeRepositorySchedulerStr),
		new(incomeSchedulerService.IncomeSchedulerServiceStr),
		new(pivotalRepository.HousePivotalRepository[pivotalModel.HousePivotal]),
		new(pivotalRepository.GroupPivotalRepository[pivotalModel.GroupPivotal]),
		new(pivotalService.PivotalServiceStr),
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
	if err := a.databaseService.DB().AutoMigrate(object.GetEntity()); err != nil {
		log.Fatal().Err(err)
	}
}

func (a *RootApplication) createCountriesService() {
	file, err := ioutil.ReadFile(fmt.Sprintf("%scontent/countries.json", environment.GetEnvironmentVariable(countriesDirVariable, "./")))

	if err != nil {
		log.Fatal().Msg("Countries is not found")
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
