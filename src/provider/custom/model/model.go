package model

import (
	"github.com/google/uuid"
	userModel "user/model"
)

type CustomProvider struct {
	Id      uuid.UUID `gorm:"primarykey"`
	Name    string    `gorm:"index:idx_name_userid,unique"`
	Details string
	UserId  uuid.UUID      `gorm:"index:idx_name_userid,unique"`
	User    userModel.User `gorm:"foreignKey:UserId"`
}

type CreateCustomProviderRequest struct {
	Name    string
	Details string
	UserId  uuid.UUID
}

type CustomProviderDto struct {
	Id      uuid.UUID
	Name    string
	Details string
	UserId  uuid.UUID
}

func (c CustomProvider) ToDto() CustomProviderDto {
	return CustomProviderDto{
		Id:      c.Id,
		Name:    c.Name,
		Details: c.Details,
		UserId:  c.UserId,
	}
}

func (c CreateCustomProviderRequest) ToEntity() CustomProvider {
	return CustomProvider{
		Id:      uuid.New(),
		Name:    c.Name,
		Details: c.Details,
		UserId:  c.UserId,
	}
}
