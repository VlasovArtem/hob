package model

import (
	houseModel "github.com/VlasovArtem/hob/src/house/model"
	"github.com/google/uuid"
	"time"
)

type Income struct {
	Id          uuid.UUID `gorm:"primarykey"`
	Name        string
	Description string
	Date        time.Time
	Sum         float32
	HouseId     uuid.UUID
	House       houseModel.House `gorm:"foreignKey:HouseId"`
}

type CreateIncomeRequest struct {
	Name        string
	Description string
	Date        time.Time
	Sum         float32
	HouseId     uuid.UUID
}

type UpdateIncomeRequest struct {
	Id          uuid.UUID
	Name        string
	Description string
	Date        time.Time
	Sum         float32
}

type IncomeDto struct {
	Id          uuid.UUID
	Name        string
	Description string
	Date        time.Time
	Sum         float32
	HouseId     uuid.UUID
}

func (i Income) ToDto() IncomeDto {
	return IncomeDto{
		Id:          i.Id,
		Name:        i.Name,
		Description: i.Description,
		Date:        i.Date,
		Sum:         i.Sum,
		HouseId:     i.HouseId,
	}
}

func (c CreateIncomeRequest) CreateToEntity() Income {
	return Income{
		Id:          uuid.New(),
		Name:        c.Name,
		Description: c.Description,
		Date:        c.Date,
		Sum:         c.Sum,
		HouseId:     c.HouseId,
	}
}

func (c UpdateIncomeRequest) UpdateToEntity() Income {
	return Income{
		Id:          c.Id,
		Name:        c.Name,
		Description: c.Description,
		Date:        c.Date,
		Sum:         c.Sum,
	}
}
