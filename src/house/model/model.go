package model

import (
	"country/model"
	"github.com/google/uuid"
)

type House struct {
	Id          uuid.UUID
	Name        string
	Country     *model.Country
	City        string
	StreetLine1 string
	StreetLine2 string
	Deleted     bool
}

type HouseResponse struct {
	Id          uuid.UUID
	Name        string
	Country     string
	City        string
	StreetLine1 string
	StreetLine2 string
}

type CreateHouseRequest struct {
	Name        string
	Country     string
	City        string
	StreetLine1 string
	StreetLine2 string
}

func (h House) ToResponse() HouseResponse {
	return HouseResponse{
		Id:          h.Id,
		Name:        h.Name,
		Country:     h.Country.Name,
		City:        h.City,
		StreetLine1: h.StreetLine1,
		StreetLine2: h.StreetLine2,
	}
}

func (c CreateHouseRequest) ToEntity(country *model.Country) House {
	return House{
		Id:          uuid.New(),
		Name:        c.Name,
		Country:     country,
		City:        c.City,
		StreetLine1: c.StreetLine1,
		StreetLine2: c.StreetLine2,
		Deleted:     true,
	}
}
