package repository

import (
	"fmt"
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/group/mocks"
	"github.com/VlasovArtem/hob/src/group/model"
	houseMocks "github.com/VlasovArtem/hob/src/house/mocks"
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

func (g *GroupRepositoryTestSuite) SetupSuite() {
	g.InitDBTestSuite()

	g.CreateRepository(
		func(service db.DatabaseService) {
			g.repository = NewGroupRepository(service)
		},
	).
		AddAfterTest(func(service db.DatabaseService) {
			database.TruncateTable(service, model.Group{})
		}).
		AddAfterSuite(func(service db.DatabaseService) {
			database.TruncateTable(service, userModel.User{})
		}).
		ExecuteMigration(userModel.User{}, model.Group{})

	g.createdUser = userMocks.GenerateUser()
	g.CreateEntity(&g.createdUser)
}

func TestGroupRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(GroupRepositoryTestSuite))
}

func (g *GroupRepositoryTestSuite) Test_Create() {
	firstGroup := mocks.GenerateGroup(g.createdUser.Id)

	actual, err := g.repository.Create(firstGroup)

	assert.Nil(g.T(), err)
	assert.Equal(g.T(), firstGroup, actual)
}

func (g *GroupRepositoryTestSuite) Test_CreateBatch() {
	firstGroup := mocks.GenerateGroup(g.createdUser.Id)
	firstGroup.Name = "Name First"
	secondGroup := mocks.GenerateGroup(g.createdUser.Id)
	secondGroup.Name = "Name Second"

	actual, err := g.repository.CreateBatch([]model.Group{firstGroup, secondGroup})

	assert.Nil(g.T(), err)
	assert.Equal(g.T(), []model.Group{firstGroup, secondGroup}, actual)
}

func (g *GroupRepositoryTestSuite) Test_FindById() {
	expected := g.createGroup()

	actual, err := g.repository.FindById(expected.Id)

	assert.Nil(g.T(), err)
	assert.Equal(g.T(), expected.ToDto(), actual)
}

func (g *GroupRepositoryTestSuite) Test_FindById_WithMissingId() {
	actual, err := g.repository.FindById(uuid.New())

	assert.ErrorIs(g.T(), err, gorm.ErrRecordNotFound)
	assert.Equal(g.T(), model.GroupDto{}, actual)
}

func (g *GroupRepositoryTestSuite) Test_FindByOwnerId() {
	expected := g.createGroup()

	actual := g.repository.FindByOwnerId(expected.OwnerId)

	assert.Equal(g.T(), []model.GroupDto{expected.ToDto()}, actual)
}

func (g *GroupRepositoryTestSuite) Test_ExistsById() {
	entity := g.createGroup()

	assert.True(g.T(), g.repository.ExistsById(entity.Id))
}

func (g *GroupRepositoryTestSuite) Test_ExistsById_WithMissingId() {
	assert.False(g.T(), g.repository.ExistsById(uuid.New()))
}

func (g *GroupRepositoryTestSuite) Test_ExistsByIds() {
	entity := g.createGroup()

	assert.True(g.T(), g.repository.ExistsByIds([]uuid.UUID{entity.Id}))
}

func (g *GroupRepositoryTestSuite) Test_ExistsByIds_WithMissingId() {
	entity := g.createGroup()

	assert.False(g.T(), g.repository.ExistsByIds([]uuid.UUID{entity.Id, uuid.New()}))
}

func (g *GroupRepositoryTestSuite) Test_DeleteById() {
	group := g.createGroup()

	assert.Nil(g.T(), g.repository.DeleteById(group.Id))
}

func (g *GroupRepositoryTestSuite) Test_DeleteById_WithMissingId() {
	assert.Nil(g.T(), g.repository.DeleteById(uuid.New()))
}

func (g *GroupRepositoryTestSuite) Test_Update() {
	entity := g.createGroup()
	newHouse := houseMocks.GenerateHouse(g.createdUser.Id)
	g.CreateEntity(newHouse)

	updatedIncome := model.UpdateGroupRequest{
		Name: fmt.Sprintf("%s-new", entity.Name),
	}

	err := g.repository.Update(entity.Id, updatedIncome)

	assert.Nil(g.T(), err)

	response, err := g.repository.FindById(entity.Id)
	assert.Nil(g.T(), err)
	assert.Equal(g.T(), model.GroupDto{
		Id:      entity.Id,
		Name:    "Name-new",
		OwnerId: g.createdUser.Id,
	}, response)
}

func (g *GroupRepositoryTestSuite) createGroup() model.Group {
	entity := mocks.GenerateGroup(g.createdUser.Id)

	g.CreateEntity(entity)

	return entity
}
