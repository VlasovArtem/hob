package mocks

import (
	"fmt"
	"github.com/VlasovArtem/hob/src/provider/custom/model"
	"github.com/google/uuid"
)

func GenerateCustomProvider(userId uuid.UUID) model.CustomProvider {
	id := uuid.New()
	return model.CustomProvider{
		Id:      id,
		Name:    fmt.Sprintf("Name%s", id),
		Details: "Details",
		UserId:  userId,
	}
}

func GenerateCustomProviderRequest() model.CreateCustomProviderRequest {
	return model.CreateCustomProviderRequest{
		Name:    "Name",
		Details: "Details",
		UserId:  uuid.New(),
	}
}

func GenerateCustomProviderDto() model.CustomProviderDto {
	return model.CustomProviderDto{
		Id:      uuid.New(),
		Name:    "Name",
		Details: "Details",
		UserId:  uuid.New(),
	}
}