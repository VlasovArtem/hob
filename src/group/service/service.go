package service

import (
	"fmt"
	"github.com/VlasovArtem/hob/src/common"
	"github.com/VlasovArtem/hob/src/common/database"
	"github.com/VlasovArtem/hob/src/common/dependency"
	interrors "github.com/VlasovArtem/hob/src/common/int-errors"
	"github.com/VlasovArtem/hob/src/common/transactional"
	"github.com/VlasovArtem/hob/src/group/model"
	"github.com/VlasovArtem/hob/src/group/repository"
	userService "github.com/VlasovArtem/hob/src/user/service"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GroupServiceStr struct {
	userService userService.UserService
	repository  repository.GroupRepository
}

func NewGroupService(
	userService userService.UserService,
	repository repository.GroupRepository,
) GroupService {
	return &GroupServiceStr{
		userService: userService,
		repository:  repository,
	}
}

func (g *GroupServiceStr) GetRequiredDependencies() []dependency.Requirements {
	return []dependency.Requirements{
		dependency.FindNameAndType(userService.UserServiceStr{}),
		dependency.FindNameAndType(repository.GroupRepositoryStr{}),
	}
}

func (g *GroupServiceStr) Initialize(factory dependency.DependenciesProvider) any {
	return NewGroupService(
		dependency.FindRequiredDependency[userService.UserServiceStr, userService.UserService](factory),
		dependency.FindRequiredDependency[repository.GroupRepositoryStr, repository.GroupRepository](factory),
	)
}

type GroupService interface {
	transactional.Transactional[GroupService]
	Add(request model.CreateGroupRequest) (model.GroupDto, error)
	AddBatch(request model.CreateGroupBatchRequest) ([]model.GroupDto, error)
	FindById(id uuid.UUID) (model.GroupDto, error)
	FindByUserId(userId uuid.UUID) []model.GroupDto
	ExistsById(id uuid.UUID) bool
	ExistsByIds(ids []uuid.UUID) bool
	DeleteById(id uuid.UUID) error
	Update(id uuid.UUID, request model.UpdateGroupRequest) error
}

func (g *GroupServiceStr) Add(request model.CreateGroupRequest) (response model.GroupDto, err error) {
	if !g.userService.ExistsById(request.OwnerId) {
		return response, interrors.NewErrNotFound("user with id %s not found", request.OwnerId)
	} else {
		entity := request.ToEntity()

		if err = g.repository.Create(&entity); err != nil {
			return response, err
		} else {
			return entity.ToDto(), nil
		}
	}
}

func (g *GroupServiceStr) AddBatch(request model.CreateGroupBatchRequest) (response []model.GroupDto, err error) {
	if len(request.Groups) == 0 {
		return make([]model.GroupDto, 0), nil
	}

	ownerIds := make(map[string]uuid.UUID)

	entities := common.MapSlice(request.Groups, func(r model.CreateGroupRequest) model.Group {
		ownerIds[r.OwnerId.String()] = r.OwnerId
		return r.ToEntity()
	})

	builder := interrors.NewBuilder()

	for stringOwnerId, ownerId := range ownerIds {
		if !g.userService.ExistsById(ownerId) {
			builder.WithDetail(fmt.Sprintf("user with id %s not found", stringOwnerId))
		}
	}

	if builder.HasErrors() {
		builder.WithMessage("Create Group Batch Request Issue")
		return nil, interrors.NewErrResponse(builder)
	}

	if batch, err := g.repository.CreateBatch(entities); err != nil {
		return nil, err
	} else {
		return common.MapSlice(batch, model.GroupToGroupDto), nil
	}
}

func (g *GroupServiceStr) FindById(id uuid.UUID) (response model.GroupDto, err error) {
	if err = g.repository.FindReceiver(&response, id); err != nil {
		return response, database.HandlerFindError(err, "group with id %s not found", id)
	} else {
		return response, nil
	}
}

func (g *GroupServiceStr) FindByUserId(userId uuid.UUID) []model.GroupDto {
	return g.repository.FindByOwnerId(userId)
}

func (g *GroupServiceStr) ExistsById(id uuid.UUID) bool {
	return g.repository.Exists(id)
}

func (g *GroupServiceStr) ExistsByIds(ids []uuid.UUID) bool {
	return g.repository.ExistsByIds(ids)
}

func (g *GroupServiceStr) DeleteById(id uuid.UUID) error {
	if !g.ExistsById(id) {
		return interrors.NewErrNotFound("group with id %s not found", id)
	}
	return g.repository.Delete(id)
}

func (g *GroupServiceStr) Update(id uuid.UUID, request model.UpdateGroupRequest) error {
	if !g.ExistsById(id) {
		return interrors.NewErrNotFound("group with id %s not found", id)
	}
	return g.repository.Update(id, request)
}

func (g *GroupServiceStr) Transactional(tx *gorm.DB) GroupService {
	return &GroupServiceStr{
		userService: g.userService.Transactional(tx),
		repository:  g.repository.Transactional(tx),
	}
}
