package handler

import (
	"country/service"
	"encoding/json"
	helper "helper/service"
	"net/http"
)

type countryHandlerObject struct {
	countryService service.CountryService
}

func NewCountryHandler(countryService service.CountryService) CountryHandler {
	return &countryHandlerObject{countryService}
}

type CountryHandler interface {
	FindAllCountries() http.HandlerFunc
	FindCountryByCode() http.HandlerFunc
}

func (c *countryHandlerObject) FindAllCountries() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		json.NewEncoder(writer).Encode(c.countryService.FindAllCountries())
	}
}

func (c *countryHandlerObject) FindCountryByCode() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if parameter, err := helper.GetRequestParameter(request, "code"); err != nil {
			helper.HandleBadRequestWithError(writer, err)
		} else {
			if err, country := c.countryService.FindCountryByCode(parameter); err != nil {
				helper.HandleErrorResponseWithError(writer, http.StatusNotFound, err)
			} else {
				json.NewEncoder(writer).Encode(country)
			}
		}
	}
}
