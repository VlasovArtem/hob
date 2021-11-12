package main

import (
	helper "common/service"
	countryHandler "country/handler"
	"country/model"
	countries "country/service"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	houseHandler "house/handler"
	houses "house/service"
	incomeHandler "income/handler"
	incomes "income/service"
	"io/ioutil"
	"log"
	meterHandler "meter/handler"
	meters "meter/service"
	"net/http"
	"os"
	paymentHandler "payment/handler"
	payments "payment/service"
	userHandler "user/handler"
	users "user/service"
)

type application struct {
	country struct {
		countryService countries.CountryService
		countryHandler countryHandler.CountryHandler
	}
	house struct {
		houseService houses.HouseService
		houseHandler houseHandler.HouseHandler
	}
	user struct {
		userService users.UserService
		userHandler userHandler.UserHandler
	}
	payment struct {
		paymentService payments.PaymentService
		paymentHandler paymentHandler.PaymentHandler
	}
	meter struct {
		meterService meters.MeterService
		meterHandler meterHandler.MeterHandler
	}
	income struct {
		incomeService incomes.IncomeService
		incomeHandler incomeHandler.IncomeHandler
	}
}

func main() {
	app := initApplication()

	router := initRouter(app)

	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":3000", router))
}

func initRouter(app *application) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	initCountryHandler(router, app)
	initHouseHandler(router, app)
	initUserHandler(router, app)
	initPaymentHandler(router, app)
	initMeterHandler(router, app)
	initIncomeHandler(router, app)

	return router
}

func initCountryHandler(router *mux.Router, handler *application) {
	countryRouter := router.PathPrefix("/api/v1/country").Subrouter()

	countryRouter.Path("/").HandlerFunc(handler.country.countryHandler.FindAll()).Methods("GET")
	countryRouter.Path("/{code}").HandlerFunc(handler.country.countryHandler.FindByCode()).Methods("GET")
}

func initHouseHandler(router *mux.Router, handler *application) {
	houseRouter := router.PathPrefix("/api/v1/house").Subrouter()

	houseRouter.Path("/").HandlerFunc(handler.house.houseHandler.Add()).Methods("POST")
	houseRouter.Path("/").HandlerFunc(handler.house.houseHandler.FindAll()).Methods("GET")
	houseRouter.Path("/{id}").HandlerFunc(handler.house.houseHandler.FindById()).Methods("GET")
}

func initUserHandler(router *mux.Router, handler *application) {
	houseRouter := router.PathPrefix("/api/v1/user").Subrouter()

	houseRouter.Path("/").HandlerFunc(handler.user.userHandler.Add()).Methods("POST")
	houseRouter.Path("/{id}").HandlerFunc(handler.user.userHandler.FindById()).Methods("GET")
}

func initPaymentHandler(router *mux.Router, handler *application) {
	houseRouter := router.PathPrefix("/api/v1/payment").Subrouter()

	houseRouter.Path("/").HandlerFunc(handler.payment.paymentHandler.Add()).Methods("POST")
	houseRouter.Path("/{id}").HandlerFunc(handler.payment.paymentHandler.FindById()).Methods("GET")
	houseRouter.Path("/house/{id}").HandlerFunc(handler.payment.paymentHandler.FindByHouseId()).Methods("GET")
	houseRouter.Path("/user/{id}").HandlerFunc(handler.payment.paymentHandler.FindByUserId()).Methods("GET")
}

func initMeterHandler(router *mux.Router, handler *application) {
	houseRouter := router.PathPrefix("/api/v1/meter").Subrouter()

	houseRouter.Path("/").HandlerFunc(handler.meter.meterHandler.Add()).Methods("POST")
	houseRouter.Path("/{id}").HandlerFunc(handler.meter.meterHandler.FindById()).Methods("GET")
	houseRouter.Path("/payment/{id}").HandlerFunc(handler.meter.meterHandler.FindByPaymentId()).Methods("GET")
}

func initIncomeHandler(router *mux.Router, handler *application) {
	houseRouter := router.PathPrefix("/api/v1/income").Subrouter()

	houseRouter.Path("/").HandlerFunc(handler.income.incomeHandler.Add()).Methods("POST")
	houseRouter.Path("/{id}").HandlerFunc(handler.income.incomeHandler.FindById()).Methods("GET")
	houseRouter.Path("/house/{id}").HandlerFunc(handler.income.incomeHandler.FindByHouseId()).Methods("GET")
}

func initApplication() *application {
	countriesService := initCountriesService()
	houseService := houses.NewHouseService(countriesService)
	userService := users.NewUserService()
	paymentService := payments.NewPaymentService(userService, houseService)
	meterService := meters.NewMeterService(paymentService)
	incomeService := incomes.NewIncomeService(houseService)

	return &application{
		country: struct {
			countryService countries.CountryService
			countryHandler countryHandler.CountryHandler
		}{
			countryService: countriesService,
			countryHandler: countryHandler.NewCountryHandler(countriesService),
		},
		house: struct {
			houseService houses.HouseService
			houseHandler houseHandler.HouseHandler
		}{
			houseService: houseService,
			houseHandler: houseHandler.NewHouseHandler(houseService),
		},
		user: struct {
			userService users.UserService
			userHandler userHandler.UserHandler
		}{
			userService: userService,
			userHandler: userHandler.NewUserHandler(userService),
		},
		payment: struct {
			paymentService payments.PaymentService
			paymentHandler paymentHandler.PaymentHandler
		}{
			paymentService: paymentService,
			paymentHandler: paymentHandler.NewPaymentHandler(paymentService),
		},
		meter: struct {
			meterService meters.MeterService
			meterHandler meterHandler.MeterHandler
		}{
			meterService: meterService,
			meterHandler: meterHandler.NewMeterHandler(meterService),
		},
		income: struct {
			incomeService incomes.IncomeService
			incomeHandler incomeHandler.IncomeHandler
		}{
			incomeService: incomeService,
			incomeHandler: incomeHandler.NewIncomeHandler(incomeService),
		},
	}
}

func initCountriesService() countries.CountryService {
	file, err := ioutil.ReadFile(fmt.Sprintf("%s/content/countries.json", os.Getenv("GOPATH")))

	if helper.HandleError(err, "Countries is not found") {
		os.Exit(1)
	}

	var countriesContent []model.Country

	if err = json.Unmarshal(file, &countriesContent); err != nil {
		log.Fatal(err)
	}

	if len(countriesContent) == 0 {
		log.Fatal("countries content is empty")
	}

	return countries.NewCountryService(countriesContent)
}
