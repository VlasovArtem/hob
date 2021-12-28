package service

import (
	"github.com/VlasovArtem/hob/src/common/int-errors"
	"github.com/VlasovArtem/hob/src/country/model"
	"github.com/rs/zerolog/log"
)

type CountryServiceObject struct {
	countriesMap map[string]model.Country
	countries    []model.Country
}

type CountryService interface {
	FindCountryByCode(code string) (model.Country, error)
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

func (c *CountryServiceObject) FindCountryByCode(code string) (country model.Country, err error) {
	if country, ok := c.countriesMap[code]; ok {
		return country, err
	}
	return country, int_errors.NewErrNotFound("country with code %s is not found", code)
}

func (c *CountryServiceObject) FindAllCountries() []model.Country {
	return c.countries
}
