package repository

import (
	"fmt"
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/group/mocks"
	"github.com/VlasovArtem/hob/src/group/model"
	houseMocks "github.com/VlasovArtem/hob/src/house/mocks"
	houseModel "github.com/VlasovArtem/hob/src/house/model"
	"github.com/VlasovArtem/hob/src/test/testhelper/database"
	userMocks "github.com/VlasovArtem/hob/src/user/mocks"
	userModel "github.com/VlasovArtem/hob/src/user/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type GroupRepositoryTestSuite struct {
	database.DBTestSuite
	repository  GroupRepository
	createdUser userModel.User
}

func (i *GroupRepositoryTestSuite) SetupSuite() {
	i.InitDBTestSuite()

	i.CreateRepository(
		func(service db.DatabaseService) {
			i.repository = NewGroupRepository(service)
		},
	).
		AddAfterTest(func(service db.DatabaseService) {
			database.TruncateTable(service, model.Group{})
		}).
		AddAfterSuite(func(service db.DatabaseService) {
			database.TruncateTable(service, userModel.User{})
		}).
		ExecuteMigration(userModel.User{}, houseModel.House{}, model.Group{})

	i.createdUser = userMocks.GenerateUser()
	i.CreateEntity(&i.createdUser)
}

func TestGroupRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(GroupRepositoryTestSuite))
}

func (i *GroupRepositoryTestSuite) Test_Create() {
	firstGroup := mocks.GenerateGroup(i.createdUser.Id)

	actual, err := i.repository.Create(firstGroup)

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), firstGroup, actual)
}

func (i *GroupRepositoryTestSuite) Test_FindById() {
	expected := i.createGroup()

	actual, err := i.repository.FindById(expected.Id)

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), expected.ToDto(), actual)
}

func (i *GroupRepositoryTestSuite) Test_FindById_WithMissingId() {
	actual, err := i.repository.FindById(uuid.New())

	assert.ErrorIs(i.T(), err, gorm.ErrRecordNotFound)
	assert.Equal(i.T(), model.GroupDto{}, actual)
}

func (i *GroupRepositoryTestSuite) Test_FindByOwnerId() {
	expected := i.createGroup()

	actual := i.repository.FindByOwnerId(expected.OwnerId)

	assert.Equal(i.T(), []model.GroupDto{expected.ToDto()}, actual)
}

func (i *GroupRepositoryTestSuite) Test_ExistsById() {
	entity := i.createGroup()

	assert.True(i.T(), i.repository.ExistsById(entity.Id))
}

func (i *GroupRepositoryTestSuite) Test_ExistsById_WithMissingId() {
	assert.False(i.T(), i.repository.ExistsById(uuid.New()))
}

func (i *GroupRepositoryTestSuite) Test_DeleteById() {
	group := i.createGroup()

	assert.Nil(i.T(), i.repository.DeleteById(group.Id))
}

func (i *GroupRepositoryTestSuite) Test_DeleteById_WithMissingId() {
	assert.Nil(i.T(), i.repository.DeleteById(uuid.New()))
}

func (i *GroupRepositoryTestSuite) Test_Update() {
	entity := i.createGroup()
	newHouse := houseMocks.GenerateHouse(i.createdUser.Id)
	i.CreateEntity(newHouse)

	updatedIncome := model.UpdateGroupRequest{
		Name: fmt.Sprintf("%s-new", entity.Name),
	}

	err := i.repository.Update(entity.Id, updatedIncome)

	assert.Nil(i.T(), err)

	response, err := i.repository.FindById(entity.Id)
	assert.Nil(i.T(), err)
	assert.Equal(i.T(), model.GroupDto{
		Id:      entity.Id,
		Name:    "Name-new",
		OwnerId: i.createdUser.Id,
	}, response)
}

func (i *GroupRepositoryTestSuite) createGroup() model.Group {
	entity := mocks.GenerateGroup(i.createdUser.Id)

	i.CreateEntity(entity)

	return entity
}
