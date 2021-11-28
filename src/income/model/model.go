package model

import (
	"github.com/google/uuid"
	houseModel "house/model"
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

type IncomeResponse struct {
	Id          uuid.UUID
	Name        string
	Description string
	Date        time.Time
	Sum         float32
	HouseId     uuid.UUID
}

func (i Income) ToResponse() IncomeResponse {
	return IncomeResponse{
		Id:          i.Id,
		Name:        i.Name,
		Description: i.Description,
		Date:        i.Date,
		Sum:         i.Sum,
		HouseId:     i.HouseId,
	}
}

func (c CreateIncomeRequest) ToEntity() Income {
	return Income{
		Id:          uuid.New(),
		Name:        c.Name,
		Description: c.Description,
		Date:        c.Date,
		Sum:         c.Sum,
		HouseId:     c.HouseId,
	}
}
