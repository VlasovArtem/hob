package service

import (
	"country/model"
	"errors"
	"fmt"
	"log"
)

type countryServiceObject struct {
	countriesMap map[string]model.Country
	countries    []model.Country
}

type CountryService interface {
	FindCountryByCode(code string) (error, model.Country)
	FindAllCountries() []model.Country
}

func NewCountryService(countries []model.Country) CountryService {
	object := &countryServiceObject{
		countriesMap: make(map[string]model.Country),
		countries:    countries,
	}

	for _, country := range object.countries {
		object.countriesMap[country.Code] = country
	}

	log.Println("Countries init completed")

	return object
}

func (c *countryServiceObject) FindCountryByCode(code string) (error, model.Country) {
	if country, ok := c.countriesMap[code]; ok {
		return nil, country
	}
	return errors.New(fmt.Sprintf("country with code %s is not found", code)), model.Country{}
}

func (c *countryServiceObject) FindAllCountries() []model.Country {
	return c.countries
}
