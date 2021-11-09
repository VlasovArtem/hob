package helper

import (
	country "country/model"
	"fmt"
	"github.com/google/uuid"
	"house/model"
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

func GenerateCreateHouseRequest() model.CreateHouseRequest {
	return model.CreateHouseRequest{
		Name:        fmt.Sprintf("Test House %s", uuid.New()),
		Country:     "UA",
		City:        "City",
		StreetLine1: "StreetLine1",
		StreetLine2: "StreetLine2",
	}
}

func GenerateHouseResponse(id uuid.UUID, name string) model.HouseResponse {
	return model.HouseResponse{
		Id:          id,
		Name:        name,
		Country:     CountryObject.Name,
		City:        "City",
		StreetLine1: "StreetLine1",
		StreetLine2: "StreetLine2",
	}
}
