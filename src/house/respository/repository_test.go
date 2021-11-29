package respository

import (
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
	h.CreateEntity(&h.createdUser)
}

func TestHouseRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(HouseRepositoryTestSuite))
}

func (h *HouseRepositoryTestSuite) TearDownSuite() {
	h.TearDown()
}

func (h *HouseRepositoryTestSuite) Test_Create() {
	house := mocks.GenerateHouse(h.createdUser.Id)

	actual, err := h.repository.Create(house)

	assert.Nil(h.T(), err)
	assert.Equal(h.T(), house, actual)
}

func (h *HouseRepositoryTestSuite) Test_Creat_WithMissingUser() {
	house := mocks.GenerateHouse(uuid.New())

	actual, err := h.repository.Create(house)

	assert.NotNil(h.T(), err)
	assert.Equal(h.T(), house, actual)
}

func (h *HouseRepositoryTestSuite) Test_FindById() {
	house := h.createHouse()

	actual, err := h.repository.FindResponseById(house.Id)

	assert.Nil(h.T(), err)
	assert.Equal(h.T(), house.ToDto(), actual)
}

func (h *HouseRepositoryTestSuite) Test_FindById_WithMissingId() {
	actual, err := h.repository.FindResponseById(uuid.New())

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

func (h *HouseRepositoryTestSuite) createHouse() model.House {
	house := mocks.GenerateHouse(h.createdUser.Id)

	create, err := h.repository.Create(house)

	assert.Nil(h.T(), err)

	return create
}
