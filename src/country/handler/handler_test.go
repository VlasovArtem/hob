package handler

import (
	"encoding/json"
	"github.com/VlasovArtem/hob/src/country/model"
	testHelperService "github.com/VlasovArtem/hob/src/test"
	"github.com/VlasovArtem/hob/src/test/testhelper"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var countryService = testhelper.InitCountryService()

var handler = NewCountryHandler(countryService)

func TestFindAllCountries(t *testing.T) {
	countries := countryService.FindAllCountries()

	testRequest := testhelper.NewTestRequest().
		WithHandler(handler.FindAll()).
		WithURL("https://test.com/api/v1/country").
		WithMethod("GET")

	body := testRequest.Verify(t, http.StatusOK)

	var responses []model.Country
	json.Unmarshal(body, &responses)

	assert.Equal(t, countries, responses)
}

func TestFindCountryByCode(t *testing.T) {
	expectedCountry := testHelperService.CountryObject

	testRequest := testhelper.NewTestRequest().
		WithHandler(handler.FindByCode()).
		WithVar("code", "UA").
		WithURL("https://test.com/api/v1/country/{code}").
		WithMethod("GET")

	body := testRequest.Verify(t, http.StatusOK)

	var response model.Country
	json.Unmarshal(body, &response)

	assert.Equal(t, *expectedCountry, response)
}
