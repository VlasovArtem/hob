package model

import (
	"github.com/VlasovArtem/hob/src/country/model"
	groupModel "github.com/VlasovArtem/hob/src/group/model"
	userModel "github.com/VlasovArtem/hob/src/user/model"
	"github.com/google/uuid"
)

type House struct {
	Id          uuid.UUID `gorm:"primarykey"`
	Name        string
	CountryCode string
	City        string
	StreetLine1 string
	StreetLine2 string
	UserId      uuid.UUID
	User        userModel.User     `gorm:"foreignKey:UserId"`
	Groups      []groupModel.Group `gorm:"many2many:house_groups"`
}

type HouseDto struct {
	Id          uuid.UUID
	Name        string
	CountryCode string
	City        string
	StreetLine1 string
	StreetLine2 string
	UserId      uuid.UUID
	Groups      []groupModel.GroupDto
}

type CreateHouseRequest struct {
	Name        string
	CountryCode string
	City        string
	StreetLine1 string
	StreetLine2 string
	UserId      uuid.UUID
	GroupIds    []uuid.UUID
}

type UpdateHouseRequest struct {
	Name        string
	CountryCode string
	City        string
	StreetLine1 string
	StreetLine2 string
	GroupIds    []uuid.UUID
}

func (h House) ToDto() HouseDto {
	var groups []groupModel.GroupDto
	for _, group := range h.Groups {
		groups = append(groups, group.ToDto())
	}

	return HouseDto{
		Id:          h.Id,
		Name:        h.Name,
		CountryCode: h.CountryCode,
		City:        h.City,
		StreetLine1: h.StreetLine1,
		StreetLine2: h.StreetLine2,
		UserId:      h.UserId,
		Groups:      groups,
	}
}

func (c CreateHouseRequest) ToEntity(country *model.Country) House {
	var groups []groupModel.Group

	for _, id := range c.GroupIds {
		groups = append(groups, groupModel.Group{Id: id})
	}

	return House{
		Id:          uuid.New(),
		Name:        c.Name,
		CountryCode: country.Code,
		City:        c.City,
		StreetLine1: c.StreetLine1,
		StreetLine2: c.StreetLine2,
		UserId:      c.UserId,
		Groups:      groups,
	}
}
