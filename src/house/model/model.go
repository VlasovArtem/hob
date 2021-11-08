package model

import "github.com/google/uuid"

var DEFAULT = House{}

type House struct {
	Id          uuid.UUID
	Name        string
	Country     string
	City        string
	StreetLine1 string
	StreetLine2 string
	Deleted     bool
}

type CreateHouseRequest struct {
	Name        string
	Country     string
	City        string
	StreetLine1 string
	StreetLine2 string
}

type CreateHouseResponse struct {
	Id uuid.UUID
}
