package repository

import (
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/group/mocks"
	"github.com/VlasovArtem/hob/src/group/model"
	"github.com/VlasovArtem/hob/src/test/testhelper"
	"github.com/VlasovArtem/hob/src/test/testhelper/database"
	userMocks "github.com/VlasovArtem/hob/src/user/mocks"
	userModel "github.com/VlasovArtem/hob/src/user/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type GroupRepositoryTestSuite struct {
	database.DBTestSuite[model.Group]
	repository  GroupRepository
	createdUser userModel.User
}

func (g *GroupRepositoryTestSuite) SetupSuite() {
	g.InitDBTestSuite()

	g.CreateRepository(
		func(service db.ModeledDatabase[model.Group]) {
			g.repository = NewGroupRepository(service)
		},
	).
		AddAfterTest(func(service db.ModeledDatabase[model.Group]) {
			testhelper.TruncateTable(service, model.Group{})
		}).
		AddAfterSuite(func(service db.ModeledDatabase[model.Group]) {
			testhelper.TruncateTable(service, userModel.User{})
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

	err := g.repository.Create(&firstGroup)

	assert.Nil(g.T(), err)

	actual, err := g.repository.Find(firstGroup.Id)

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

func (g *GroupRepositoryTestSuite) Test_FindByOwnerId() {
	expected := g.createGroup()

	actual := g.repository.FindByOwnerId(expected.OwnerId)

	assert.Equal(g.T(), []model.GroupDto{expected.ToDto()}, actual)
}

func (g *GroupRepositoryTestSuite) Test_ExistsByIds() {
	entity := g.createGroup()

	assert.True(g.T(), g.repository.ExistsByIds([]uuid.UUID{entity.Id}))
}

func (g *GroupRepositoryTestSuite) Test_ExistsByIds_WithMissingId() {
	entity := g.createGroup()

	assert.False(g.T(), g.repository.ExistsByIds([]uuid.UUID{entity.Id, uuid.New()}))
}

func (g *GroupRepositoryTestSuite) createGroup() model.Group {
	entity := mocks.GenerateGroup(g.createdUser.Id)

	g.CreateEntity(entity)

	return entity
}
