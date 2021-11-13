package test

import (
	country "country/model"
	"fmt"
	"github.com/google/uuid"
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

func GetUserResponse(id uuid.UUID, email string) usermodel.UserResponse {
	return usermodel.UserResponse{
		Id:        id,
		FirstName: "First name",
		LastName:  "Last name",
		Email:     email,
	}
}

func GetCreateUserRequest() usermodel.CreateUserRequest {
	return usermodel.CreateUserRequest{
		FirstName: "First name",
		LastName:  "Last name",
		Email:     fmt.Sprintf("mail%s@mail.com", uuid.New()),
		Password:  []byte("password"),
	}
}
