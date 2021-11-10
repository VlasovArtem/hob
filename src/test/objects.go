package test

import (
	country "country/model"
	"fmt"
	"github.com/google/uuid"
	housemodel "house/model"
	usermodel "user/model"
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

func GenerateCreateHouseRequest() housemodel.CreateHouseRequest {
	return housemodel.CreateHouseRequest{
		Name:        fmt.Sprintf("Test House %s", uuid.New()),
		Country:     "UA",
		City:        "City",
		StreetLine1: "StreetLine1",
		StreetLine2: "StreetLine2",
	}
}

func GenerateHouseResponse(id uuid.UUID, name string) housemodel.HouseResponse {
	return housemodel.HouseResponse{
		Id:          id,
		Name:        name,
		Country:     CountryObject.Name,
		City:        "City",
		StreetLine1: "StreetLine1",
		StreetLine2: "StreetLine2",
	}
}

func GetUserResponse(id uuid.UUID, email string) usermodel.UserResponse {
	return usermodel.UserResponse{
		Id:        id,
		FirstName: "First name",
		LastName:  "Last name",
		Email:     email,
	}
}

func GetUser(id uuid.UUID) usermodel.User {
	user := GetCreateUserRequest().ToEntity()

	user.Id = id

	return user
}

func GetCreateUserRequest() usermodel.CreateUserRequest {
	return usermodel.CreateUserRequest{
		FirstName: "First name",
		LastName:  "Last name",
		Email:     fmt.Sprintf("mail%s@mail.com", uuid.New()),
		Password:  []byte("password"),
	}
}
