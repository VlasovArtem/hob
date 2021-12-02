package handler

import (
	"encoding/json"
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/rest"
	"github.com/VlasovArtem/hob/src/country/service"
	"github.com/gorilla/mux"
	"net/http"
)

type CountryHandlerObject struct {
	countryService service.CountryService
}

func NewCountryHandler(countryService service.CountryService) CountryHandler {
	return &CountryHandlerObject{countryService}
}

func (c *CountryHandlerObject) Initialize(factory dependency.DependenciesFactory) interface{} {
	return factory.Add(NewCountryHandler(factory.FindRequiredByObject(service.CountryServiceObject{}).(service.CountryService)))
}

func (c *CountryHandlerObject) Init(router *mux.Router) {
	subrouter := router.PathPrefix("/api/v1/country").Subrouter()

	subrouter.Path("/").HandlerFunc(c.FindAll()).Methods("GET")
	subrouter.Path("/{code}").HandlerFunc(c.FindByCode()).Methods("GET")
}

type CountryHandler interface {
	FindAll() http.HandlerFunc
	FindByCode() http.HandlerFunc
}

func (c *CountryHandlerObject) FindAll() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		json.NewEncoder(writer).Encode(c.countryService.FindAllCountries())
	}
}

func (c *CountryHandlerObject) FindByCode() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if parameter, err := rest.GetRequestParameter(request, "code"); err != nil {
			rest.HandleBadRequestWithError(writer, err)
		} else {
			if err, country := c.countryService.FindCountryByCode(parameter); err != nil {
				rest.HandleErrorResponseWithError(writer, http.StatusNotFound, err)
			} else {
				json.NewEncoder(writer).Encode(country)
			}
		}
	}
}
