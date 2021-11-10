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
	"house/service"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type application struct {
	country struct {
		countryService countries.CountryService
		countryHandler countryHandler.CountryHandler
	}
	house struct {
		houseService service.HouseService
		houseHandler houseHandler.HouseHandler
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

func initApplication() *application {
	countriesService := initCountriesService()

	houseService := service.NewHouseService(countriesService)

	return &application{
		country: struct {
			countryService countries.CountryService
			countryHandler countryHandler.CountryHandler
		}{
			countryService: countriesService,
			countryHandler: countryHandler.NewCountryHandler(countriesService),
		},
		house: struct {
			houseService service.HouseService
			houseHandler houseHandler.HouseHandler
		}{
			houseService: houseService,
			houseHandler: houseHandler.NewHouseHandler(houseService),
		},
	}
}

func initCountriesService() countries.CountryService {
	file, err := ioutil.ReadFile(fmt.Sprintf("%s/content/countries.json", os.Getenv("GOPATH")))

	if helper.HandleError(err, "Countries is not found") {
		os.Exit(1)
	}

	var countriesContent []model.Country

	json.Unmarshal(file, &countriesContent)

	if len(countriesContent) == 0 {
		log.Fatal("countries content is empty")
	}

	return countries.NewCountryService(countriesContent)
}
