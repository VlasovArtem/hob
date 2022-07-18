package handler

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/rest"
	"github.com/VlasovArtem/hob/src/country/service"
	"github.com/gorilla/mux"
	"net/http"
)

type CountryHandlerStr struct {
	countryService service.CountryService
}

func NewCountryHandler(countryService service.CountryService) CountryHandler {
	return &CountryHandlerStr{countryService}
}

func (c *CountryHandlerStr) GetRequiredDependencies() []dependency.Requirements {
	return []dependency.Requirements{
		dependency.FindNameAndType(service.CountryServiceStr{}),
	}
}

func (c *CountryHandlerStr) Initialize(factory dependency.DependenciesProvider) any {
	return NewCountryHandler(dependency.FindRequiredDependency[service.CountryServiceStr, service.CountryService](factory))
}

func (c *CountryHandlerStr) Init(router *mux.Router) {
	subrouter := router.PathPrefix("/api/v1/countries").Subrouter()

	subrouter.Path("/").HandlerFunc(c.FindAll()).Methods("GET")
	subrouter.Path("/{code}").HandlerFunc(c.FindByCode()).Methods("GET")
}

type CountryHandler interface {
	FindAll() http.HandlerFunc
	FindByCode() http.HandlerFunc
}

func (c *CountryHandlerStr) FindAll() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		rest.PerformResponse(writer, c.countryService.FindAllCountries(), nil)
	}
}

func (c *CountryHandlerStr) FindByCode() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if parameter, err := rest.GetRequestParameter(request, "code"); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			country, err := c.countryService.FindCountryByCode(parameter)
			rest.PerformResponse(writer, country, err)
		}
	}
}
