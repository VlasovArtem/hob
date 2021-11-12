package model

import (
	"github.com/google/uuid"
	"time"
)

type Income struct {
	Id          uuid.UUID
	Name        string
	Description string
	Date        time.Time
	Sum         float32
	HouseId     uuid.UUID
}

type CreateIncomeRequest struct {
	Name        string
	Description string
	Date        time.Time
	Sum         float32
	HouseId     uuid.UUID
}

type IncomeResponse struct {
	Income
}

func (i Income) ToResponse() IncomeResponse {
	return IncomeResponse{i}
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
