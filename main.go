package main

import "C"
import (
	"app"
	helper "common/service"
	"country/model"
	countries "country/service"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	houseHandler "house/handler"
	houseRepository "house/respository"
	houseService "house/service"
	incomeHandler "income/handler"
	incomeRepository "income/repository"
	incomeSchedulerHandler "income/scheduler/handler"
	incomeSchedulerRepository "income/scheduler/repository"
	incomeSchedulerService "income/scheduler/service"
	incomeService "income/service"
	"io/ioutil"
	meterHandler "meter/handler"
	meterRepository "meter/repository"
	meterService "meter/service"
	"net/http"
	"os"
	paymentHandler "payment/handler"
	paymentRepository "payment/repository"
	paymentSchedulerHandler "payment/scheduler/handler"
	paymentSchedulerRepository "payment/scheduler/repository"
	paymentSchedulerService "payment/scheduler/service"
	paymentService "payment/service"
	customProviderHandler "provider/custom/handler"
	customProviderRepository "provider/custom/repository"
	customProviderService "provider/custom/service"
	providerHandler "provider/handler"
	providerRepository "provider/repository"
	providerService "provider/service"
	"scheduler"
	userHandler "user/handler"
	userRepository "user/repository"
	userService "user/service"
	userRequestValidator "user/validator"
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
	file, err := ioutil.ReadFile(fmt.Sprintf("%s/content/countries.json", os.Getenv("GOPATH")))

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
