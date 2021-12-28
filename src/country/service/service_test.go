package service

import (
	"encoding/json"
	int_errors "github.com/VlasovArtem/hob/src/common/int-errors"
	"github.com/VlasovArtem/hob/src/country/model"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"reflect"
	"testing"
)

var countriesService = func() CountryService {
	file, err := ioutil.ReadFile("../../../content/countries.json")

	if err != nil {
		log.Fatal().Err(err).Msg("countries file not fount")
	}

	var countriesContent []model.Country

	json.Unmarshal(file, &countriesContent)

	return NewCountryService(countriesContent)
}()

func TestFindCountryByCode(t *testing.T) {
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		args    args
		err     error
		country model.Country
	}{
		{
			name: "with existing country",
			args: args{
				"UA",
			},
			err: nil,
			country: model.Country{
				Name:    "Ukraine",
				Code:    "UA",
				Capital: "Kiev",
				Region:  "EU",
				Currency: model.Currency{
					Code:   "UAH",
					Name:   "Ukrainian hryvnia",
					Symbol: "₴",
				},
				Language: model.Language{
					Code: "uk",
					Name: "Ukrainian",
				},
				Flag: "https://restcountries.eu/data/ukr.svg",
			},
		}, {
			name: "with not existing country",
			args: args{
				"INVALID",
			},
			err:     int_errors.NewErrNotFound("country with code INVALID is not found"),
			country: model.Country{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := countriesService.FindCountryByCode(tt.args.code)
			if !reflect.DeepEqual(got1, tt.err) {
				t.Errorf("FindByCode() got = %v, want %v", got1, tt.err)
			}
			if !reflect.DeepEqual(got, tt.country) {
				t.Errorf("FindByCode() got1 = %v, want %v", got, tt.country)
			}
		})
	}
}

func TestFindAllCountries(t *testing.T) {
	assert.Equal(t, 249, len(countriesService.FindAllCountries()), "FindAll() should have 249 countries")
}
