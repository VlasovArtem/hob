package model

import (
	groupModel "github.com/VlasovArtem/hob/src/group/model"
	houseModel "github.com/VlasovArtem/hob/src/house/model"
	"github.com/google/uuid"
	"time"
)

type Pivotal struct {
	Id                      uuid.UUID `gorm:"primarykey"`
	Income                  float64
	Payments                float64
	Total                   float64
	LatestIncomeUpdateDate  time.Time
	LatestPaymentUpdateDate time.Time
}

type HousePivotal struct {
	Pivotal Pivotal `gorm:"embedded"`
	HouseId uuid.UUID
	House   houseModel.House `gorm:"foreignKey:HouseId"`
}

type GroupPivotal struct {
	Pivotal Pivotal `gorm:"embedded"`
	GroupId uuid.UUID
	Group   groupModel.Group `gorm:"foreignKey:GroupId"`
}

type HousePivotalDto struct {
	Pivotal
	HouseId uuid.UUID
}

type GroupPivotalDto struct {
	Pivotal
	Group groupModel.Group
}

func (h HousePivotal) ToDto() HousePivotalDto {
	return HousePivotalDto{
		Pivotal: h.Pivotal,
		HouseId: h.HouseId,
	}
}

func (g GroupPivotal) ToDto() GroupPivotalDto {
	return GroupPivotalDto{
		Pivotal: g.Pivotal,
		Group:   g.Group,
	}
}

type TotalPivotalDto struct {
	Income   float64
	Payments float64
	Total    float64
}

type PivotalResponseDto struct {
	House  HousePivotalDto
	Groups []GroupPivotalDto
	Total  TotalPivotalDto
}
