package main

import "C"
import (
	"encoding/json"
	"github.com/VlasovArtem/hob/src/app"
	helper "github.com/VlasovArtem/hob/src/common/service"
	"github.com/VlasovArtem/hob/src/country/model"
	countries "github.com/VlasovArtem/hob/src/country/service"
	houseHandler "github.com/VlasovArtem/hob/src/house/handler"
	houseRepository "github.com/VlasovArtem/hob/src/house/respository"
	houseService "github.com/VlasovArtem/hob/src/house/service"
	incomeHandler "github.com/VlasovArtem/hob/src/income/handler"
	incomeRepository "github.com/VlasovArtem/hob/src/income/repository"
	incomeSchedulerHandler "github.com/VlasovArtem/hob/src/income/scheduler/handler"
	incomeSchedulerRepository "github.com/VlasovArtem/hob/src/income/scheduler/repository"
	incomeSchedulerService "github.com/VlasovArtem/hob/src/income/scheduler/service"
	incomeService "github.com/VlasovArtem/hob/src/income/service"
	meterHandler "github.com/VlasovArtem/hob/src/meter/handler"
	meterRepository "github.com/VlasovArtem/hob/src/meter/repository"
	meterService "github.com/VlasovArtem/hob/src/meter/service"
	paymentHandler "github.com/VlasovArtem/hob/src/payment/handler"
	paymentRepository "github.com/VlasovArtem/hob/src/payment/repository"
	paymentSchedulerHandler "github.com/VlasovArtem/hob/src/payment/scheduler/handler"
	paymentSchedulerRepository "github.com/VlasovArtem/hob/src/payment/scheduler/repository"
	paymentSchedulerService "github.com/VlasovArtem/hob/src/payment/scheduler/service"
	paymentService "github.com/VlasovArtem/hob/src/payment/service"
	customProviderHandler "github.com/VlasovArtem/hob/src/provider/custom/handler"
	customProviderRepository "github.com/VlasovArtem/hob/src/provider/custom/repository"
	customProviderService "github.com/VlasovArtem/hob/src/provider/custom/service"
	providerHandler "github.com/VlasovArtem/hob/src/provider/handler"
	providerRepository "github.com/VlasovArtem/hob/src/provider/repository"
	providerService "github.com/VlasovArtem/hob/src/provider/service"
	"github.com/VlasovArtem/hob/src/scheduler"
	userHandler "github.com/VlasovArtem/hob/src/user/handler"
	userRepository "github.com/VlasovArtem/hob/src/user/repository"
	userService "github.com/VlasovArtem/hob/src/user/service"
	userRequestValidator "github.com/VlasovArtem/hob/src/user/validator"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)

	application := app.NewApplicationService(router)

	createCountriesService(application)

	addAutoInitializingDependencies(application)

	http.Handle("/", router)

	log.Fatal().Err(http.ListenAndServe(":3000", router))
}

func addAutoInitializingDependencies(application app.ApplicationService) {
	application.
		AddAutoDependency(new(userRequestValidator.UserRequestValidatorObject)).
		AddAutoDependency(new(userRepository.UserRepositoryObject)).
		AddAutoDependency(new(userService.UserServiceObject)).
		AddHandler(new(userHandler.UserHandlerObject))

	application.
		AddAutoDependency(new(houseRepository.HouseRepositoryObject)).
		AddAutoDependency(new(houseService.HouseServiceObject)).
		AddHandler(new(houseHandler.HouseHandlerObject))

	application.
		AddAutoDependency(new(scheduler.SchedulerServiceObject))

	application.
		AddAutoDependency(new(providerRepository.ProviderRepositoryObject)).
		AddAutoDependency(new(providerService.ProviderServiceObject)).
		AddHandler(new(providerHandler.ProviderHandlerObject))

	application.
		AddAutoDependency(new(customProviderRepository.CustomProviderRepositoryObject)).
		AddAutoDependency(new(customProviderService.CustomProviderServiceObject)).
		AddHandler(new(customProviderHandler.CustomProviderHandlerObject))

	application.
		AddAutoDependency(new(paymentRepository.PaymentRepositoryObject)).
		AddAutoDependency(new(paymentService.PaymentServiceObject)).
		AddHandler(new(paymentHandler.PaymentHandlerObject))

	application.
		AddAutoDependency(new(paymentSchedulerRepository.PaymentSchedulerRepositoryObject)).
		AddAutoDependency(new(paymentSchedulerService.PaymentSchedulerServiceObject)).
		AddHandler(new(paymentSchedulerHandler.PaymentSchedulerHandlerObject))

	application.
		AddAutoDependency(new(meterRepository.MeterRepositoryObject)).
		AddAutoDependency(new(meterService.MeterServiceObject)).
		AddHandler(new(meterHandler.MeterHandlerObject))

	application.
		AddAutoDependency(new(incomeRepository.IncomeRepositoryObject)).
		AddAutoDependency(new(incomeService.IncomeServiceObject)).
		AddHandler(new(incomeHandler.IncomeHandlerObject))

	application.
		AddAutoDependency(new(incomeSchedulerRepository.IncomeSchedulerRepositoryObject)).
		AddAutoDependency(new(incomeSchedulerService.IncomeSchedulerServiceObject)).
		AddHandler(new(incomeSchedulerHandler.IncomeSchedulerHandlerObject))
}

func createCountriesService(app app.ApplicationService) {
	file, err := ioutil.ReadFile("./content/countries.json")

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

	app.AddDependency(countries.NewCountryService(countriesContent))
}
