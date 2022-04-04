package service

import (
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/common"
	interrors "github.com/VlasovArtem/hob/src/common/int-errors"
	"github.com/VlasovArtem/hob/src/group/mocks"
	"github.com/VlasovArtem/hob/src/group/model"
	"github.com/VlasovArtem/hob/src/test/testhelper"
	userMocks "github.com/VlasovArtem/hob/src/user/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type GroupServiceTestSuite struct {
	testhelper.MockTestSuite[GroupService]
	users           *userMocks.UserService
	groupRepository *mocks.GroupRepository
}

func TestGroupServiceTestSuite(t *testing.T) {
	testingSuite := &GroupServiceTestSuite{}
	testingSuite.TestObjectGenerator = func() GroupService {
		testingSuite.users = new(userMocks.UserService)
		testingSuite.groupRepository = new(mocks.GroupRepository)
		return NewGroupService(testingSuite.users, testingSuite.groupRepository)
	}

	suite.Run(t, testingSuite)
}

func (g *GroupServiceTestSuite) Test_Add() {
	var entity model.Group
	request := mocks.GenerateCreateGroupRequest()

	g.users.On("ExistsById", request.OwnerId).Return(true)
	g.groupRepository.On("Create", mock.Anything).Return(func(group model.Group) model.Group {
		entity = group
		return group
	}, nil)

	income, err := g.TestO.Add(request)

	assert.Nil(g.T(), err)
	assert.Equal(g.T(), entity.ToDto(), income)
}

func (g *GroupServiceTestSuite) Test_Add_WithUserNotExists() {
	request := mocks.GenerateCreateGroupRequest()

	g.users.On("ExistsById", request.OwnerId).Return(false)

	payment, err := g.TestO.Add(request)

	assert.Equal(g.T(), interrors.NewErrNotFound("user with id %s not found", request.OwnerId), err)
	assert.Equal(g.T(), model.GroupDto{}, payment)

	g.groupRepository.AssertNotCalled(g.T(), "Create", mock.Anything)
}

func (g *GroupServiceTestSuite) Test_Add_WithErrorFromRepository() {
	expectedError := errors.New("error")
	request := mocks.GenerateCreateGroupRequest()

	g.users.On("ExistsById", request.OwnerId).Return(true)
	g.groupRepository.On("Create", mock.Anything).Return(model.Group{}, expectedError)

	income, err := g.TestO.Add(request)

	assert.Equal(g.T(), expectedError, err)
	assert.Equal(g.T(), model.GroupDto{}, income)
}

func (g *GroupServiceTestSuite) Test_AddBatch() {
	var entities []model.Group
	request := mocks.GenerateCreateGroupBatchRequest(2)

	g.users.On("ExistsById", mock.Anything).Return(true)
	g.groupRepository.On("CreateBatch", mock.Anything).Return(func(groups []model.Group) []model.Group {
		entities = groups
		return groups
	}, nil)

	result, err := g.TestO.AddBatch(request)

	assert.Nil(g.T(), err)
	assert.Equal(g.T(), common.MapSlice(entities, model.GroupToGroupDto), result)

	g.users.AssertCalled(g.T(), "ExistsById", entities[0].OwnerId)
	g.users.AssertCalled(g.T(), "ExistsById", entities[1].OwnerId)
	g.groupRepository.AssertCalled(g.T(), "CreateBatch", entities)
}

func (g *GroupServiceTestSuite) Test_AddBatch_WithUserNotExists() {
	request := mocks.GenerateCreateGroupBatchRequest(2)

	g.users.On("ExistsById", request.Groups[0].OwnerId).Return(true)
	g.users.On("ExistsById", request.Groups[1].OwnerId).Return(false)

	result, err := g.TestO.AddBatch(request)

	var expectedResult []model.GroupDto

	assert.Equal(g.T(), expectedResult, result)

	builder := interrors.NewBuilder()
	builder.WithMessage("Create Group Batch Request Issue")
	builder.WithDetail(fmt.Sprintf("user with id %s not found", request.Groups[1].OwnerId.String()))

	expectedError := interrors.NewErrResponse(builder).(*interrors.ErrResponse)
	actualError := err.(*interrors.ErrResponse)

	assert.Equal(g.T(), *expectedError, *actualError)

	g.groupRepository.AssertNotCalled(g.T(), "CreateBatch", mock.Anything)
}

func (g *GroupServiceTestSuite) Test_AddBatch_WithEmptyGroups() {
	request := mocks.GenerateCreateGroupBatchRequest(0)

	result, err := g.TestO.AddBatch(request)

	assert.Nil(g.T(), err)
	assert.Equal(g.T(), []model.GroupDto{}, result)

	g.users.AssertNotCalled(g.T(), "ExistsById", mock.Anything)
	g.groupRepository.AssertNotCalled(g.T(), "CreateBatch", mock.Anything)
}

func (g *GroupServiceTestSuite) Test_AddBatch_WithErrorFromRepository() {
	expectedError := errors.New("error")
	request := mocks.GenerateCreateGroupBatchRequest(1)

	g.users.On("ExistsById", mock.Anything).Return(true)
	g.groupRepository.On("CreateBatch", mock.Anything).Return(nil, expectedError)

	result, err := g.TestO.AddBatch(request)

	var expected []model.GroupDto

	assert.Equal(g.T(), expectedError, err)
	assert.Equal(g.T(), expected, result)
}

func (g *GroupServiceTestSuite) Test_FindById() {
	group := mocks.GenerateGroupDto()

	g.groupRepository.On("FindById", group.Id).Return(group, nil)

	actual, err := g.TestO.FindById(group.Id)

	assert.Nil(g.T(), err)
	assert.Equal(g.T(), group, actual)
}

func (g *GroupServiceTestSuite) Test_FindById_WithNotExistingId() {
	id := uuid.New()

	g.groupRepository.On("FindById", id).Return(model.GroupDto{}, gorm.ErrRecordNotFound)

	actual, err := g.TestO.FindById(id)

	assert.Equal(g.T(), interrors.NewErrNotFound("group with id %s not found", id), err)
	assert.Equal(g.T(), model.GroupDto{}, actual)
}

func (g *GroupServiceTestSuite) Test_FindById_WithError() {
	id := uuid.New()
	expectedError := errors.New("test")

	g.groupRepository.On("FindById", id).Return(model.GroupDto{}, expectedError)

	actual, err := g.TestO.FindById(id)

	assert.Equal(g.T(), expectedError, err)
	assert.Equal(g.T(), model.GroupDto{}, actual)
}

func (g *GroupServiceTestSuite) Test_FindByUserId() {
	groups := []model.GroupDto{mocks.GenerateGroupDto()}

	g.groupRepository.On("FindByOwnerId", groups[0].OwnerId).Return(groups, nil)

	actual := g.TestO.FindByUserId(groups[0].OwnerId)

	assert.Equal(g.T(), groups, actual)
}

func (g *GroupServiceTestSuite) Test_FindByUserId_WithNotExistingRecords() {
	var groups []model.GroupDto

	userId := uuid.New()

	g.groupRepository.On("FindByOwnerId", userId).Return(groups, nil)

	actual := g.TestO.FindByUserId(userId)

	assert.Equal(g.T(), groups, actual)
}

func (g *GroupServiceTestSuite) Test_ExistsById() {
	id := uuid.New()

	g.groupRepository.On("ExistsById", id).Return(true)

	assert.True(g.T(), g.TestO.ExistsById(id))
}

func (g *GroupServiceTestSuite) Test_ExistsById_WithNotExists() {
	id := uuid.New()

	g.groupRepository.On("ExistsById", id).Return(false)

	assert.False(g.T(), g.TestO.ExistsById(id))
}

func (g *GroupServiceTestSuite) Test_ExistsByIds() {
	ids := []uuid.UUID{uuid.New()}

	g.groupRepository.On("ExistsByIds", ids).Return(true)

	assert.True(g.T(), g.TestO.ExistsByIds(ids))
}

func (g *GroupServiceTestSuite) Test_ExistsByIds_WithNotExists() {
	ids := []uuid.UUID{uuid.New()}

	g.groupRepository.On("ExistsByIds", ids).Return(false)

	assert.False(g.T(), g.TestO.ExistsByIds(ids))
}

func (g *GroupServiceTestSuite) Test_DeleteById() {
	id := uuid.New()

	g.groupRepository.On("ExistsById", id).Return(true)
	g.groupRepository.On("DeleteById", id).Return(nil)

	assert.Nil(g.T(), g.TestO.DeleteById(id))
}

func (g *GroupServiceTestSuite) Test_DeleteById_WithNotExists() {
	id := uuid.New()

	g.groupRepository.On("ExistsById", id).Return(false)

	assert.Equal(g.T(), interrors.NewErrNotFound("group with id %s not found", id), g.TestO.DeleteById(id))

	g.groupRepository.AssertNotCalled(g.T(), "DeleteById", id)
}

func (g *GroupServiceTestSuite) Test_Update() {
	id, request := mocks.GenerateUpdateGroupRequest()

	g.groupRepository.On("ExistsById", id).Return(true)
	g.groupRepository.On("Update", id, request).Return(nil)

	assert.Nil(g.T(), g.TestO.Update(id, request))

	g.groupRepository.AssertCalled(g.T(), "Update", id, request)
}

func (g *GroupServiceTestSuite) Test_Update_WithErrorFromDatabase() {
	id, request := mocks.GenerateUpdateGroupRequest()

	g.groupRepository.On("ExistsById", id).Return(true)
	g.groupRepository.On("Update", id, request).Return(errors.New("test"))

	err := g.TestO.Update(id, request)
	assert.Equal(g.T(), errors.New("test"), err)
}

func (g *GroupServiceTestSuite) Test_Update_WithNotExists() {
	id, request := mocks.GenerateUpdateGroupRequest()

	g.groupRepository.On("ExistsById", id).Return(false)

	err := g.TestO.Update(id, request)
	assert.Equal(g.T(), interrors.NewErrNotFound("group with id %s not found", id), err)

	g.groupRepository.AssertNotCalled(g.T(), "Update", id, request)
}
