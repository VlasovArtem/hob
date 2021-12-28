package respository

import (
	"fmt"
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/house/mocks"
	"github.com/VlasovArtem/hob/src/house/model"
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
	database.DBTestSuite
	createdUser userModel.User
	repository  HouseRepository
}

func (h *HouseRepositoryTestSuite) SetupSuite() {
	h.InitDBTestSuite()

	h.CreateRepository(
		func(service db.DatabaseService) {
			h.repository = NewHouseRepository(service)
		},
	).
		AddMigrators(userModel.User{}, model.House{})

	h.createdUser = userMocks.GenerateUser()
	h.CreateConstantEntity(&h.createdUser)
}

func TestHouseRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(HouseRepositoryTestSuite))
}

func (h *HouseRepositoryTestSuite) Test_Create() {
	house := mocks.GenerateHouse(h.createdUser.Id)

	actual, err := h.repository.Create(house)

	assert.Nil(h.T(), err)
	assert.Equal(h.T(), house, actual)

	h.Delete(house)
}

func (h *HouseRepositoryTestSuite) Test_Creat_WithMissingUser() {
	house := mocks.GenerateHouse(uuid.New())

	actual, err := h.repository.Create(house)

	assert.NotNil(h.T(), err)
	assert.Equal(h.T(), house, actual)
}

func (h *HouseRepositoryTestSuite) Test_FindById() {
	house := h.createHouse()

	actual, err := h.repository.FindDtoById(house.Id)

	assert.Nil(h.T(), err)
	assert.Equal(h.T(), house.ToDto(), actual)
}

func (h *HouseRepositoryTestSuite) Test_FindById_WithMissingId() {
	actual, err := h.repository.FindDtoById(uuid.New())

	assert.ErrorIs(h.T(), err, gorm.ErrRecordNotFound)
	assert.Equal(h.T(), model.HouseDto{}, actual)
}

func (h *HouseRepositoryTestSuite) Test_FindByUserId() {
	house := h.createHouse()

	actual := h.repository.FindResponseByUserId(house.UserId)

	var actualResponse model.HouseDto

	for _, response := range actual {
		if response.Id == house.Id {
			actualResponse = response
			break
		}
	}
	assert.Equal(h.T(), house.ToDto(), actualResponse)
}

func (h *HouseRepositoryTestSuite) Test_FindByUserId_WithMissingId() {
	actual := h.repository.FindResponseByUserId(uuid.New())

	assert.Equal(h.T(), []model.HouseDto{}, actual)
}

func (h *HouseRepositoryTestSuite) Test_ExistsById() {
	house := h.createHouse()

	actual := h.repository.ExistsById(house.Id)

	assert.True(h.T(), actual)
}

func (h *HouseRepositoryTestSuite) Test_ExistsById_WithMissingId() {
	actual := h.repository.ExistsById(uuid.New())

	assert.False(h.T(), actual)
}

func (h *HouseRepositoryTestSuite) Test_DeleteById() {
	house := h.createHouse()

	assert.Nil(h.T(), h.repository.DeleteById(house.Id))
}

func (h *HouseRepositoryTestSuite) Test_DeleteById_WithMissingId() {
	assert.Nil(h.T(), h.repository.DeleteById(uuid.New()))
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

	err := h.repository.Update(house.Id, updatedHouse)

	assert.Nil(h.T(), err)

	response, err := h.repository.FindDtoById(house.Id)
	assert.Nil(h.T(), err)
	assert.Equal(h.T(), model.HouseDto{
		Id:          house.Id,
		Name:        "Name-new",
		CountryCode: "US",
		City:        "City-new",
		StreetLine1: "Street Line 1-new",
		StreetLine2: "Street Line 2-new",
		UserId:      house.UserId,
	}, response)
}

func (h *HouseRepositoryTestSuite) Test_Update_WithMissingId() {
	assert.Nil(h.T(), h.repository.DeleteById(uuid.New()))
}

func (h *HouseRepositoryTestSuite) createHouse() (house model.House) {
	house = mocks.GenerateHouse(h.createdUser.Id)

	h.CreateEntity(&house)

	return house
}
