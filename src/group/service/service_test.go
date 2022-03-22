package service

import (
	"errors"
	interrors "github.com/VlasovArtem/hob/src/common/int-errors"
	"github.com/VlasovArtem/hob/src/group/mocks"
	"github.com/VlasovArtem/hob/src/group/model"
	userMocks "github.com/VlasovArtem/hob/src/user/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"testing"
)

var (
	users           *userMocks.UserService
	groupRepository *mocks.GroupRepository
)

func serviceGenerator() GroupService {
	users = new(userMocks.UserService)
	groupRepository = new(mocks.GroupRepository)

	return NewGroupService(users, groupRepository)
}

func Test_Add(t *testing.T) {
	service := serviceGenerator()

	var entity model.Group
	request := mocks.GenerateCreateGroupRequest()

	users.On("ExistsById", request.OwnerId).Return(true)
	groupRepository.On("Create", mock.Anything).Return(func(group model.Group) model.Group {
		entity = group
		return group
	}, nil)

	income, err := service.Add(request)

	assert.Nil(t, err)
	assert.Equal(t, entity.ToDto(), income)
}

func Test_Add_WithUserNotExists(t *testing.T) {
	service := serviceGenerator()

	request := mocks.GenerateCreateGroupRequest()

	users.On("ExistsById", request.OwnerId).Return(false)

	payment, err := service.Add(request)

	assert.Equal(t, interrors.NewErrNotFound("user with id %s not found", request.OwnerId), err)
	assert.Equal(t, model.GroupDto{}, payment)

	groupRepository.AssertNotCalled(t, "Create", mock.Anything)
}

func Test_Add_WithErrorFromRepository(t *testing.T) {
	service := serviceGenerator()

	expectedError := errors.New("error")
	request := mocks.GenerateCreateGroupRequest()

	users.On("ExistsById", request.OwnerId).Return(true)
	groupRepository.On("Create", mock.Anything).Return(model.Group{}, expectedError)

	income, err := service.Add(request)

	assert.Equal(t, expectedError, err)
	assert.Equal(t, model.GroupDto{}, income)
}

func Test_FindById(t *testing.T) {
	service := serviceGenerator()

	group := mocks.GenerateGroupDto()

	groupRepository.On("FindById", group.Id).Return(group, nil)

	actual, err := service.FindById(group.Id)

	assert.Nil(t, err)
	assert.Equal(t, group, actual)
}

func Test_FindById_WithNotExistingId(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()

	groupRepository.On("FindById", id).Return(model.GroupDto{}, gorm.ErrRecordNotFound)

	actual, err := service.FindById(id)

	assert.Equal(t, interrors.NewErrNotFound("group with id %s not found", id), err)
	assert.Equal(t, model.GroupDto{}, actual)
}

func Test_FindById_WithError(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()
	expectedError := errors.New("test")

	groupRepository.On("FindById", id).Return(model.GroupDto{}, expectedError)

	actual, err := service.FindById(id)

	assert.Equal(t, expectedError, err)
	assert.Equal(t, model.GroupDto{}, actual)
}

func Test_FindByUserId(t *testing.T) {
	service := serviceGenerator()

	groups := []model.GroupDto{mocks.GenerateGroupDto()}

	groupRepository.On("FindByOwnerId", groups[0].OwnerId).Return(groups, nil)

	actual := service.FindByUserId(groups[0].OwnerId)

	assert.Equal(t, groups, actual)
}

func Test_FindByUserId_WithNotExistingRecords(t *testing.T) {
	service := serviceGenerator()

	var groups []model.GroupDto

	userId := uuid.New()

	groupRepository.On("FindByOwnerId", userId).Return(groups, nil)

	actual := service.FindByUserId(userId)

	assert.Equal(t, groups, actual)
}

func Test_ExistsById(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()

	groupRepository.On("ExistsById", id).Return(true)

	assert.True(t, service.ExistsById(id))
}

func Test_ExistsById_WithNotExists(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()

	groupRepository.On("ExistsById", id).Return(false)

	assert.False(t, service.ExistsById(id))
}

func Test_DeleteById(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()

	groupRepository.On("ExistsById", id).Return(true)
	groupRepository.On("DeleteById", id).Return(nil)

	assert.Nil(t, service.DeleteById(id))
}

func Test_DeleteById_WithNotExists(t *testing.T) {
	service := serviceGenerator()

	id := uuid.New()

	groupRepository.On("ExistsById", id).Return(false)

	assert.Equal(t, interrors.NewErrNotFound("group with id %s not found", id), service.DeleteById(id))

	groupRepository.AssertNotCalled(t, "DeleteById", id)
}

func Test_Update(t *testing.T) {
	houseService := serviceGenerator()

	id, request := mocks.GenerateUpdateGroupRequest()

	groupRepository.On("ExistsById", id).Return(true)
	groupRepository.On("Update", id, request).Return(nil)

	assert.Nil(t, houseService.Update(id, request))

	groupRepository.AssertCalled(t, "Update", id, request)
}

func Test_Update_WithErrorFromDatabase(t *testing.T) {
	houseService := serviceGenerator()

	id, request := mocks.GenerateUpdateGroupRequest()

	groupRepository.On("ExistsById", id).Return(true)
	groupRepository.On("Update", id, request).Return(errors.New("test"))

	err := houseService.Update(id, request)
	assert.Equal(t, errors.New("test"), err)
}

func Test_Update_WithNotExists(t *testing.T) {
	houseService := serviceGenerator()

	id, request := mocks.GenerateUpdateGroupRequest()

	groupRepository.On("ExistsById", id).Return(false)

	err := houseService.Update(id, request)
	assert.Equal(t, interrors.NewErrNotFound("group with id %s not found", id), err)

	groupRepository.AssertNotCalled(t, "Update", id, request)
}
