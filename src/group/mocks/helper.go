package mocks

import (
	"fmt"
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

func GenerateCreateGroupBatchRequest(number int) model.CreateGroupBatchRequest {
	requests := make([]model.CreateGroupRequest, 0)
	if number == 0 {
		return model.CreateGroupBatchRequest{
			Groups: requests,
		}
	}
	for i := 0; i < number; i++ {
		requests = append(requests, model.CreateGroupRequest{
			Name:    fmt.Sprintf("name-%d", i),
			OwnerId: uuid.New(),
		})
	}

	return model.CreateGroupBatchRequest{Groups: requests}
}

func GenerateUpdateGroupRequest() (uuid.UUID, model.UpdateGroupRequest) {
	return uuid.New(), model.UpdateGroupRequest{
		Name: "new-name",
	}
}
