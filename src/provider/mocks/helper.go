package mocks

import (
	"fmt"
	"github.com/VlasovArtem/hob/src/provider/model"
	"github.com/google/uuid"
)

func GenerateProvider(userId uuid.UUID) model.Provider {
	id := uuid.New()
	return model.Provider{
		Id:      id,
		Name:    fmt.Sprintf("%s-Provider", id),
		Details: "Details",
		UserId:  userId,
	}
}

func GenerateCreateProviderRequest() model.CreateProviderRequest {
	return model.CreateProviderRequest{
		Name:    "Name",
		Details: "Details",
		UserId:  uuid.New(),
	}
}

func GenerateUpdateProviderRequest() model.UpdateProviderRequest {
	return model.UpdateProviderRequest{
		Name:    "Name",
		Details: "Details",
	}
}

func GenerateProviderDto() model.ProviderDto {
	return model.ProviderDto{
		Id:      uuid.New(),
		Name:    "Name",
		Details: "Details",
		UserId:  uuid.New(),
	}
}
