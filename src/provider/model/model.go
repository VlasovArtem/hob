package model

import (
	userModel "github.com/VlasovArtem/hob/src/user/model"
	"github.com/google/uuid"
)

type Provider struct {
	Id      uuid.UUID `gorm:"primarykey;type:uuid"`
	Name    string    `gorm:"index:idx_name_userid,unique"`
	Details string
	UserId  uuid.UUID      `gorm:"index:idx_name_userid,unique"`
	User    userModel.User `gorm:"foreignKey:UserId"`
}

type CreateProviderRequest struct {
	Name    string
	Details string
	UserId  uuid.UUID
}

type UpdateProviderRequest struct {
	Name    string
	Details string
}

type ProviderDto struct {
	Id      uuid.UUID
	Name    string
	Details string
	UserId  uuid.UUID
}

func (p Provider) ToDto() ProviderDto {
	return ProviderDto{
		Id:      p.Id,
		Name:    p.Name,
		Details: p.Details,
		UserId:  p.UserId,
	}
}

func (c CreateProviderRequest) ToEntity() Provider {
	return Provider{
		Id:      uuid.New(),
		Name:    c.Name,
		Details: c.Details,
		UserId:  c.UserId,
	}
}

func (u UpdateProviderRequest) ToEntity(id uuid.UUID) Provider {
	return Provider{
		Id:      id,
		Name:    u.Name,
		Details: u.Details,
	}
}
