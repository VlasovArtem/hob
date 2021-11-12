package handler

import (
	"common/rest"
	"country/service"
	"encoding/json"
	"net/http"
)

type countryHandlerObject struct {
	countryService service.CountryService
}

func NewCountryHandler(countryService service.CountryService) CountryHandler {
	return &countryHandlerObject{countryService}
}

type CountryHandler interface {
	FindAll() http.HandlerFunc
	FindByCode() http.HandlerFunc
}

func (c *countryHandlerObject) FindAll() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		json.NewEncoder(writer).Encode(c.countryService.FindAllCountries())
	}
}

func (c *countryHandlerObject) FindByCode() http.HandlerFunc {
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
