package model

import (
	userModel "github.com/VlasovArtem/hob/src/user/model"
	"github.com/google/uuid"
)

type Group struct {
	Id      uuid.UUID `gorm:"primarykey"`
	Name    string    `gorm:"index:idx_name_group,unique"`
	OwnerId uuid.UUID
	Owner   userModel.User `gorm:"foreignKey:OwnerId"`
}

type GroupDto struct {
	Id      uuid.UUID
	Name    string
	OwnerId uuid.UUID
}

type CreateGroupRequest struct {
	Name    string
	OwnerId uuid.UUID
}

type UpdateGroupRequest struct {
	Name string
}

func (g Group) ToDto() GroupDto {
	return GroupDto{
		Id:      g.Id,
		Name:    g.Name,
		OwnerId: g.OwnerId,
	}
}

func (c CreateGroupRequest) ToEntity() Group {
	return Group{
		Id:      uuid.New(),
		Name:    c.Name,
		OwnerId: c.OwnerId,
	}
}
