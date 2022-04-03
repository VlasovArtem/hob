package service

import (
	"github.com/VlasovArtem/hob/src/common/database"
	"github.com/VlasovArtem/hob/src/common/dependency"
	interrors "github.com/VlasovArtem/hob/src/common/int-errors"
	"github.com/VlasovArtem/hob/src/group/model"
	"github.com/VlasovArtem/hob/src/group/repository"
	userService "github.com/VlasovArtem/hob/src/user/service"
	"github.com/google/uuid"
)

type GroupServiceObject struct {
	userService userService.UserService
	repository  repository.GroupRepository
}

func NewGroupService(
	userService userService.UserService,
	repository repository.GroupRepository,
) GroupService {
	return &GroupServiceObject{
		userService: userService,
		repository:  repository,
	}
}

func (h *GroupServiceObject) Initialize(factory dependency.DependenciesProvider) any {
	return NewGroupService(
		dependency.FindRequiredDependency[userService.UserServiceObject, userService.UserService](factory),
		dependency.FindRequiredDependency[repository.GroupRepositoryObject, repository.GroupRepository](factory),
	)
}

type GroupService interface {
	Add(house model.CreateGroupRequest) (model.GroupDto, error)
	FindById(id uuid.UUID) (model.GroupDto, error)
	FindByUserId(userId uuid.UUID) []model.GroupDto
	ExistsById(id uuid.UUID) bool
	ExistsByIds(ids []uuid.UUID) bool
	DeleteById(id uuid.UUID) error
	Update(id uuid.UUID, request model.UpdateGroupRequest) error
}

func (h *GroupServiceObject) Add(request model.CreateGroupRequest) (response model.GroupDto, err error) {
	if !h.userService.ExistsById(request.OwnerId) {
		return response, interrors.NewErrNotFound("user with id %s not found", request.OwnerId)
	} else {
		entity := request.ToEntity()

		if entity, err := h.repository.Create(entity); err != nil {
			return response, err
		} else {
			return entity.ToDto(), nil
		}
	}
}

func (h *GroupServiceObject) FindById(id uuid.UUID) (response model.GroupDto, err error) {
	if response, err = h.repository.FindById(id); err != nil {
		return response, database.HandlerFindError(err, "group with id %s not found", id)
	} else {
		return response, nil
	}
}

func (h *GroupServiceObject) FindByUserId(userId uuid.UUID) []model.GroupDto {
	return h.repository.FindByOwnerId(userId)
}

func (h *GroupServiceObject) ExistsById(id uuid.UUID) bool {
	return h.repository.ExistsById(id)
}

func (h *GroupServiceObject) ExistsByIds(ids []uuid.UUID) bool {
	return h.repository.ExistsByIds(ids)
}

func (h *GroupServiceObject) DeleteById(id uuid.UUID) error {
	if !h.ExistsById(id) {
		return interrors.NewErrNotFound("group with id %s not found", id)
	}
	return h.repository.DeleteById(id)
}

func (h *GroupServiceObject) Update(id uuid.UUID, request model.UpdateGroupRequest) error {
	if !h.ExistsById(id) {
		return interrors.NewErrNotFound("group with id %s not found", id)
	}
	return h.repository.Update(id, request)
}
