package mocks

import (
	"github.com/VlasovArtem/hob/src/group/model"
	"github.com/google/uuid"
)

func GenerateGroup(ownerId uuid.UUID) model.Group {
	return model.Group{
		Id:      uuid.New(),
		Name:    "Name",
		OwnerId: ownerId,
	}
}

func GenerateGroupDto() model.GroupDto {
	return model.GroupDto{
		Id:      uuid.New(),
		Name:    "Name",
		OwnerId: uuid.New(),
	}
}

func GenerateCreateGroupRequest() model.CreateGroupRequest {
	return model.CreateGroupRequest{
		Name:    "name",
		OwnerId: uuid.New(),
	}
}

func GenerateUpdateGroupRequest() (uuid.UUID, model.UpdateGroupRequest) {
	return uuid.New(), model.UpdateGroupRequest{
		Name: "new-name",
	}
}
