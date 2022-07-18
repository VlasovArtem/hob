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
	"github.com/VlasovArtem/hob/src/pivotal/calculator"
	pivotalRepository "github.com/VlasovArtem/hob/src/pivotal/repository"
	pivotalService "github.com/VlasovArtem/hob/src/pivotal/service"
	providerRepository "github.com/VlasovArtem/hob/src/provider/repository"
	providerService "github.com/VlasovArtem/hob/src/provider/service"
	"github.com/VlasovArtem/hob/src/scheduler"
	userRepository "github.com/VlasovArtem/hob/src/user/repository"
	userService "github.com/VlasovArtem/hob/src/user/service"
	userRequestValidator "github.com/VlasovArtem/hob/src/user/validator"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
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

var migratorType = reflect.TypeOf((*dependency.ObjectDatabaseMigrator[any])(nil)).Elem()

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
	beans := []dependency.ObjectDependencyInitializer{
		new(userRequestValidator.UserRequestValidatorStr),
		new(userRepository.UserRepositoryObject),
		new(userService.UserServiceStr),
		new(repository.GroupRepositoryStr),
		new(groupService.GroupServiceStr),
		new(houseRepository.HouseRepositoryStr),
		new(houseService.HouseServiceStr),
		new(scheduler.SchedulerServiceStr),
		new(providerRepository.ProviderRepositoryStr),
		new(providerService.ProviderServiceStr),
		new(paymentRepository.PaymentRepositoryStr),
		new(paymentService.PaymentServiceStr),
		new(paymentSchedulerRepository.PaymentSchedulerRepositoryStr),
		new(paymentSchedulerService.PaymentSchedulerServiceStr),
		new(meterRepository.MeterRepositoryStr),
		new(meterService.MeterServiceStr),
		new(incomeRepository.IncomeRepositoryStr),
		new(incomeService.IncomeServiceStr),
		new(incomeSchedulerRepository.IncomeRepositorySchedulerStr),
		new(incomeSchedulerService.IncomeSchedulerServiceStr),
		new(pivotalRepository.HousePivotalRepository),
		new(pivotalRepository.GroupPivotalRepository),
		new(pivotalService.PivotalServiceStr),
		new(calculator.PivotalCalculatorServiceStr),
	}

	slices.SortFunc(beans, func(i, j dependency.ObjectDependencyInitializer) bool {
		return len(i.GetRequiredDependencies()) < len(j.GetRequiredDependencies())
	})

	beansRound := beans
	var newBeansRound []dependency.ObjectDependencyInitializer

	for len(beansRound) != 0 {
		for _, initializer := range beansRound {
			if a.isDependencyIsReadyForInit(initializer) {
				autoDependency := a.DependenciesFactory.AddAutoDependency(initializer)

				if reflect.TypeOf(autoDependency).Implements(migratorType) {
					a.migrate(autoDependency.(dependency.ObjectDatabaseMigrator[any]))
				}
			} else {
				newBeansRound = append(newBeansRound, initializer)
			}
		}
		beansRound = newBeansRound
		newBeansRound = []dependency.ObjectDependencyInitializer{}
	}
}

func (a *RootApplication) isDependencyIsReadyForInit(initializer dependency.ObjectDependencyInitializer) bool {
	for _, requirement := range initializer.GetRequiredDependencies() {
		if a.DependenciesFactory.FindByName(requirement.Name, false) == nil &&
			a.DependenciesFactory.FindByType(requirement.Type, false) == nil {
			return false
		}
	}
	return true
}

func (a *RootApplication) migrate(object dependency.ObjectDatabaseMigrator[any]) {
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
