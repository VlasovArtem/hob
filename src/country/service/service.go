package service

import (
	"country/model"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
)

type CountryServiceObject struct {
	countriesMap map[string]model.Country
	countries    []model.Country
}

type CountryService interface {
	FindCountryByCode(code string) (error, model.Country)
	FindAllCountries() []model.Country
}

func NewCountryService(countries []model.Country) CountryService {
	object := &CountryServiceObject{
		countriesMap: make(map[string]model.Country),
		countries:    countries,
	}

	for _, country := range object.countries {
		object.countriesMap[country.Code] = country
	}

	log.Info().Msg("Countries init completed")

	return object
}

func (c *CountryServiceObject) FindCountryByCode(code string) (error, model.Country) {
	if country, ok := c.countriesMap[code]; ok {
		return nil, country
	}
	return errors.New(fmt.Sprintf("country with code %s is not found", code)), model.Country{}
}

func (c *CountryServiceObject) FindAllCountries() []model.Country {
	return c.countries
}
