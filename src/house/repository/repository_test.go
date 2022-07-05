package repository

import (
	"fmt"
	"github.com/VlasovArtem/hob/src/db"
	groupModel "github.com/VlasovArtem/hob/src/group/model"
	"github.com/VlasovArtem/hob/src/house/mocks"
	"github.com/VlasovArtem/hob/src/house/model"
	"github.com/VlasovArtem/hob/src/test/testhelper"
	"github.com/VlasovArtem/hob/src/test/testhelper/database"
	userMocks "github.com/VlasovArtem/hob/src/user/mocks"
	userModel "github.com/VlasovArtem/hob/src/user/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type HouseRepositoryTestSuite struct {
	database.DBTestSuite[model.House]
	createdUser userModel.User
	repository  HouseRepository
}

func (h *HouseRepositoryTestSuite) SetupSuite() {
	h.InitDBTestSuite()

	h.CreateRepository(
		func(service db.ModeledDatabase[model.House]) {
			h.repository = NewHouseRepository(service)
		},
	).
		AddAfterSuite(
			func(service db.ModeledDatabase[model.House]) {
				testhelper.TruncateTable(service, userModel.User{})
			},
		).
		AddAfterTest(
			func(service db.ModeledDatabase[model.House]) {
				service.DB().Exec("DELETE FROM house_groups")
				testhelper.TruncateTable(service, groupModel.Group{})
				testhelper.TruncateTable(service, model.House{})
			},
		).
		ExecuteMigration(userModel.User{}, model.House{}, groupModel.Group{})

	h.createdUser = userMocks.GenerateUser()
	h.CreateEntity(&h.createdUser)
}

func TestHouseRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(HouseRepositoryTestSuite))
}

func (h *HouseRepositoryTestSuite) Test_CreateBatch() {
	first := mocks.GenerateHouse(h.createdUser.Id)
	first.Name = "Name First"
	second := mocks.GenerateHouse(h.createdUser.Id)
	second.Name = "Name Second"

	actual, err := h.repository.CreateBatch([]model.House{first, second})

	assert.Nil(h.T(), err)
	assert.Equal(h.T(), []model.House{first, second}, actual)
}

func (h *HouseRepositoryTestSuite) Test_CreateBatch_WithGroupId() {
	group := groupModel.Group{
		Id:      uuid.New(),
		Name:    "Test Group",
		OwnerId: h.createdUser.Id,
	}
	h.CreateEntity(group)

	first := mocks.GenerateHouse(h.createdUser.Id)
	first.Name = "Name First"
	first.Groups = []groupModel.Group{group}
	second := mocks.GenerateHouse(h.createdUser.Id)
	second.Name = "Name Second"
	second.Groups = []groupModel.Group{group}

	actual, err := h.repository.CreateBatch([]model.House{first, second})

	assert.Nil(h.T(), err)
	assert.Equal(h.T(), []model.House{first, second}, actual)
}

func (h *HouseRepositoryTestSuite) Test_FindById() {
	house := h.createHouse()

	actual, err := h.repository.FindById(house.Id)

	house.Groups = []groupModel.Group{}

	assert.Nil(h.T(), err)
	assert.Equal(h.T(), house, actual)
}

func (h *HouseRepositoryTestSuite) Test_FindById_WithMissingId() {
	actual, err := h.repository.FindById(uuid.New())

	assert.ErrorIs(h.T(), err, gorm.ErrRecordNotFound)
	assert.Equal(h.T(), model.House{}, actual)
}

func (h *HouseRepositoryTestSuite) Test_FindByUserId() {
	house := h.createHouse()

	actual := h.repository.FindByUserId(house.UserId)

	house.Groups = []groupModel.Group{}

	assert.Equal(h.T(), []model.House{house}, actual)
}

func (h *HouseRepositoryTestSuite) Test_FindByUserId_WithMissingId() {
	actual := h.repository.FindByUserId(uuid.New())

	assert.Equal(h.T(), []model.House{}, actual)
}

func (h *HouseRepositoryTestSuite) Test_Update() {
	house := h.createHouse()

	updatedHouse := model.UpdateHouseRequest{
		Name:        fmt.Sprintf("%s-new", house.Name),
		CountryCode: "US",
		City:        fmt.Sprintf("%s-new", house.City),
		StreetLine1: fmt.Sprintf("%s-new", house.StreetLine1),
		StreetLine2: fmt.Sprintf("%s-new", house.StreetLine2),
	}

	err := h.repository.UpdateByRequest(house.Id, updatedHouse)

	assert.Nil(h.T(), err)

	response, err := h.repository.FindById(house.Id)
	assert.Nil(h.T(), err)
	assert.Equal(h.T(), model.House{
		Id:          house.Id,
		Name:        "Name-new",
		CountryCode: "US",
		City:        "City-new",
		StreetLine1: "Street Line 1-new",
		StreetLine2: "Street Line 2-new",
		UserId:      house.UserId,
		Groups:      []groupModel.Group{},
	}, response)
}

func (h *HouseRepositoryTestSuite) Test_UpdateByRequest_WithAddingGroups() {
	group := groupModel.Group{
		Id:      uuid.New(),
		Name:    "Test Group",
		OwnerId: h.createdUser.Id,
	}
	h.CreateEntity(group)

	house := h.createHouse()

	updatedHouse := model.UpdateHouseRequest{
		GroupIds: []uuid.UUID{group.Id},
	}

	err := h.repository.UpdateByRequest(house.Id, updatedHouse)

	assert.Nil(h.T(), err)

	response, err := h.repository.FindById(house.Id)
	assert.Nil(h.T(), err)
	assert.Equal(h.T(), []groupModel.Group{group}, response.Groups)
}

func (h *HouseRepositoryTestSuite) Test_UpdateByRequest_WithUpdatingGroups() {
	group := groupModel.Group{
		Id:      uuid.New(),
		Name:    "Test Group",
		OwnerId: h.createdUser.Id,
	}

	house := h.createHouseWithGroups([]groupModel.Group{group})

	newGroup := groupModel.Group{
		Id:      uuid.New(),
		Name:    "New Test Group",
		OwnerId: h.createdUser.Id,
	}

	h.CreateEntity(newGroup)

	updatedHouse := model.UpdateHouseRequest{
		GroupIds: []uuid.UUID{newGroup.Id},
	}

	err := h.repository.UpdateByRequest(house.Id, updatedHouse)

	assert.Nil(h.T(), err)

	response, err := h.repository.FindById(house.Id)
	assert.Nil(h.T(), err)
	assert.Equal(h.T(), []groupModel.Group{newGroup}, response.Groups)
}

func (h *HouseRepositoryTestSuite) Test_UpdateByRequest_WithExtendingGroups() {
	group := groupModel.Group{
		Id:      uuid.New(),
		Name:    "Test Group",
		OwnerId: h.createdUser.Id,
	}

	house := h.createHouseWithGroups([]groupModel.Group{group})

	newGroup := groupModel.Group{
		Id:      uuid.New(),
		Name:    "New Test Group",
		OwnerId: h.createdUser.Id,
	}

	h.CreateEntity(newGroup)

	updatedHouse := model.UpdateHouseRequest{
		GroupIds: []uuid.UUID{group.Id, newGroup.Id},
	}

	err := h.repository.UpdateByRequest(house.Id, updatedHouse)

	assert.Nil(h.T(), err)

	response, err := h.repository.FindById(house.Id)
	assert.Nil(h.T(), err)
	assert.Equal(h.T(), []groupModel.Group{group, newGroup}, response.Groups)
}

func (h *HouseRepositoryTestSuite) Test_UpdateByRequest_WithDeletingGroups() {
	group := groupModel.Group{
		Id:      uuid.New(),
		Name:    "Test Group",
		OwnerId: h.createdUser.Id,
	}

	house := h.createHouseWithGroups([]groupModel.Group{group})

	updatedHouse := model.UpdateHouseRequest{
		GroupIds: []uuid.UUID{},
	}

	err := h.repository.UpdateByRequest(house.Id, updatedHouse)

	assert.Nil(h.T(), err)

	response, err := h.repository.FindById(house.Id)
	assert.Nil(h.T(), err)
	assert.Equal(h.T(), []groupModel.Group{}, response.Groups)
}

func (h *HouseRepositoryTestSuite) Test_UpdateByRequest_WithMissingId() {
	assert.NotNil(h.T(), h.repository.UpdateByRequest(uuid.New(), model.UpdateHouseRequest{}))
}

func (h *HouseRepositoryTestSuite) createHouse() (house model.House) {
	house = mocks.GenerateHouse(h.createdUser.Id)

	h.CreateEntity(&house)

	return house
}

func (h *HouseRepositoryTestSuite) createHouseWithGroups(groups []groupModel.Group) (house model.House) {
	for _, group := range groups {
		h.CreateEntity(group)
	}

	house = mocks.GenerateHouse(h.createdUser.Id)
	house.Groups = groups

	h.CreateEntity(&house)

	return house

}
