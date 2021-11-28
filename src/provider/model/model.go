package model

import "github.com/google/uuid"

type Provider struct {
	Id      uuid.UUID `gorm:"primarykey;type:uuid"`
	Name    string    `gorm:"unique"`
	Details string
}

type CreateProviderRequest struct {
	Name    string
	Details string
}

type ProviderDto struct {
	Id      uuid.UUID
	Name    string
	Details string
}

func (p Provider) ToDto() ProviderDto {
	return ProviderDto{
		Id:      p.Id,
		Name:    p.Name,
		Details: p.Details,
	}
}

func (c CreateProviderRequest) ToEntity() Provider {
	return Provider{
		Id:      uuid.New(),
		Name:    c.Name,
		Details: c.Details,
	}
}
