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
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

	return router
}

func initCountryHandler(router *mux.Router, handler *application) {
	countryRouter := router.PathPrefix("/api/v1/country").Subrouter()

	countryRouter.Path("/").HandlerFunc(handler.country.countryHandler.FindAllCountries()).Methods("GET")
	countryRouter.Path("/{code}").HandlerFunc(handler.country.countryHandler.FindAllCountries()).Methods("GET")
}

func initHouseHandler(router *mux.Router, handler *application) {
	houseRouter := router.PathPrefix("/api/v1/house").Subrouter()

	houseRouter.Path("/").HandlerFunc(handler.house.houseHandler.AddHouseHandler()).Methods("POST")
	houseRouter.Path("/").HandlerFunc(handler.house.houseHandler.FindAllHousesHandler()).Methods("GET")
	houseRouter.Path("/{id}").HandlerFunc(handler.house.houseHandler.FindHouseByIdHandler()).Methods("GET")
}

func initUserHandler(router *mux.Router, handler *application) {
	houseRouter := router.PathPrefix("/api/v1/user").Subrouter()

	houseRouter.Path("/").HandlerFunc(handler.user.userHandler.AddUser()).Methods("POST")
	houseRouter.Path("/{id}").HandlerFunc(handler.user.userHandler.FindUserById()).Methods("GET")
}

func initApplication() *application {
	countriesService := initCountriesService()
	houseService := houses.NewHouseService(countriesService)
	userService := users.NewUserService()

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
