package model

import (
	"github.com/VlasovArtem/hob/src/common"
	groupModel "github.com/VlasovArtem/hob/src/group/model"
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
	House       houseModel.House   `gorm:"foreignKey:HouseId"`
	Groups      []groupModel.Group `gorm:"many2many:income_groups"`
}

type CreateIncomeRequest struct {
	Name        string
	Description string
	Date        time.Time
	Sum         float32
	HouseId     uuid.UUID
	GroupIds    []uuid.UUID
}

type CreateIncomeBatchRequest struct {
	Incomes []CreateIncomeRequest
}

type UpdateIncomeRequest struct {
	Name        string
	Description string
	Date        time.Time
	Sum         float32
	GroupIds    []uuid.UUID
}

type IncomeDto struct {
	Id          uuid.UUID
	Name        string
	Description string
	Date        time.Time
	Sum         float32
	HouseId     uuid.UUID
	Groups      []groupModel.GroupDto
}

func (i Income) ToDto() IncomeDto {
	return IncomeDto{
		Id:          i.Id,
		Name:        i.Name,
		Description: i.Description,
		Date:        i.Date,
		Sum:         i.Sum,
		HouseId:     i.HouseId,
		Groups:      common.MapSlice(i.Groups, groupModel.GroupToGroupDto),
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
		Groups: common.MapSlice(c.GroupIds, func(groupId uuid.UUID) groupModel.Group {
			return groupModel.Group{Id: groupId}
		}),
	}
}

func IncomeToDto(income Income) IncomeDto {
	return income.ToDto()
}
