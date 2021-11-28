package test

import (
	country "country/model"
)

var CountryObject = &country.Country{
	Name:    "Ukraine",
	Code:    "UA",
	Capital: "Kiev",
	Region:  "EU",
	Currency: country.Currency{
		Code:   "UAH",
		Name:   "Ukrainian hryvnia",
		Symbol: "₴",
	},
	Language: country.Language{
		Code: "uk",
		Name: "Ukrainian",
	},
	Flag: "https://restcountries.eu/data/ukr.svg",
}
